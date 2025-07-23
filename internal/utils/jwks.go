package utils

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

type jwk struct {
	Kid string   `json:"kid"`
	Use string   `json:"use"` // "sig" or "enc"
	X5c []string `json:"x5c"` // certificate chain, first is leaf
}

type jwksResponse struct {
	Keys []jwk `json:"keys"`
}

type discovery struct {
	JwksURI string `json:"jwks_uri"`
}

type jwkManager struct {
	sync.RWMutex
	discoveryURL string
	jwksURL      string
	keys         map[string]*rsa.PublicKey // kid → RSA public key (signing keys only)
	interval     time.Duration
	quit         chan struct{}
}

var manager *jwkManager

// InitJWKs discovers the JWKS endpoint and keeps it refreshed.
func InitJWKs(realmURL string, interval time.Duration) error {
	discURL := strings.TrimRight(realmURL, "/") + "/.well-known/openid-configuration"

	resp, err := http.Get(discURL)
	if err != nil {
		return fmt.Errorf("discovery GET %s: %w", discURL, err)
	}
	defer resp.Body.Close()

	var d discovery
	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return fmt.Errorf("decode discovery JSON: %w", err)
	}
	if d.JwksURI == "" {
		return errors.New("openid‑configuration missing jwks_uri")
	}

	keys, err := fetchKeys(d.JwksURI)
	if err != nil {
		return err
	}

	manager = &jwkManager{
		discoveryURL: discURL,
		jwksURL:      d.JwksURI,
		keys:         keys,
		interval:     interval,
		quit:         make(chan struct{}),
	}
	logrus.Infof("JWKS initialised from %s (refresh every %s)", d.JwksURI, interval)
	go manager.autoRefresh()
	return nil
}

// fetchKeys downloads jwksURL and keeps only "use":"sig" RSA keys.
func fetchKeys(jwksURL string) (map[string]*rsa.PublicKey, error) {
	resp, err := http.Get(jwksURL)
	if err != nil {
		return nil, fmt.Errorf("GET %s: %w", jwksURL, err)
	}
	defer resp.Body.Close()

	var jr jwksResponse
	if err := json.NewDecoder(resp.Body).Decode(&jr); err != nil {
		return nil, fmt.Errorf("decode JWKS: %w", err)
	}

	out := make(map[string]*rsa.PublicKey)
	for _, k := range jr.Keys {
		if k.Use != "sig" {
			continue
		}
		if len(k.X5c) == 0 {
			logrus.Warnf("skip kid %s: empty x5c", k.Kid)
			continue
		}
		der, err := base64.StdEncoding.DecodeString(k.X5c[0])
		if err != nil {
			logrus.Warnf("kid %s: x5c decode error: %v", k.Kid, err)
			continue
		}
		cert, err := x509.ParseCertificate(der)
		if err != nil {
			logrus.Warnf("kid %s: parse cert error: %v", k.Kid, err)
			continue
		}
		rsaKey, ok := cert.PublicKey.(*rsa.PublicKey)
		if !ok {
			logrus.Warnf("kid %s: public key is not RSA", k.Kid)
			continue
		}
		out[k.Kid] = rsaKey
	}
	if len(out) == 0 {
		return nil, errors.New("JWKS contained no signing keys")
	}
	return out, nil
}

// autoRefresh polls discovery & JWKS at the configured interval.
func (j *jwkManager) autoRefresh() {
	t := time.NewTicker(j.interval)
	defer t.Stop()

	for {
		select {
		case <-t.C:
			keys, err := fetchKeys(j.jwksURL)
			if err != nil {
				logrus.Errorf("JWKS refresh failed: %v", err)
				continue
			}
			j.Lock()
			j.keys = keys
			j.Unlock()
			logrus.Infof("JWKS refreshed: %d signing keys", len(keys))

		case <-j.quit:
			return
		}
	}
}

// Keyfunc returns a jwt.Keyfunc that picks RSA / RSAPSS signing keys only.
func Keyfunc() jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		// allow any RSA or RSA‑PSS signing algorithm
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			if _, ok := token.Method.(*jwt.SigningMethodRSAPSS); !ok {
				return nil, fmt.Errorf("alg %q not RSA/RSAPSS", token.Header["alg"])
			}
		}
		if manager == nil {
			return nil, errors.New("jwks not initialised")
		}

		kid, _ := token.Header["kid"].(string)

		manager.RLock()
		defer manager.RUnlock()

		if pub, ok := manager.keys[kid]; ok {
			return pub, nil
		}
		// fallback to first available signing key
		for _, pub := range manager.keys {
			return pub, nil
		}
		return nil, errors.New("no signing keys in memory")
	}
}

// StopJWKs stops the background refresh goroutine.
func StopJWKs() {
	if manager != nil {
		close(manager.quit)
	}
}
