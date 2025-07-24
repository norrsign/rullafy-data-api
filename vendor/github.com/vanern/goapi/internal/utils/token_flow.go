// internal/utils/token_flow.go
package utils

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"

	"github.com/vanern/goapi/config"
)

// RunTokenFlows picks the right flow based on which fields are set.
// NOTE: TokenConfig has been removed; we now use config.Token_t directly.
func RunTokenFlows(cfg config.Token_t) error {
	// remote Keycloak mode
	if cfg.JWTRealmURL != "" {
		return runKeycloakFlows(cfg)
	}
	// local JWT generation
	return runLocalFlow(cfg)
}

func runKeycloakFlows(cfg config.Token_t) error {
	if cfg.ClientID == "" {
		return errors.New("--client-id is required for Keycloak mode")
	}

	// PKCE start: no username, no password, no code, no verifier
	if cfg.Username == "" && cfg.Password == "" && cfg.Code == "" && cfg.CodeVerifier == "" {
		return pkceStart(cfg)
	}
	// PKCE exchange
	if cfg.Code != "" && cfg.CodeVerifier != "" {
		_, err := pkceExchange(cfg)
		return err
	}
	// ROPC
	if cfg.Username != "" && cfg.Password != "" {
		_, err := ropc(cfg)
		return err
	}
	// Client Credentials
	_, err := clientCredentials(cfg)
	return err
}

func pkceStart(cfg config.Token_t) error {
	verifier, err := newCodeVerifier()
	if err != nil {
		return fmt.Errorf("generate code_verifier: %w", err)
	}
	challenge := codeChallenge(verifier)

	params := url.Values{
		"client_id":             {cfg.ClientID},
		"response_type":         {"code"},
		"redirect_uri":          {cfg.RedirectURI},
		"scope":                 {"openid profile email"},
		"code_challenge":        {challenge},
		"code_challenge_method": {"S256"},
	}
	authURL := fmt.Sprintf("%s/protocol/openid-connect/auth?%s",
		strings.TrimRight(cfg.JWTRealmURL, "/"),
		params.Encode(),
	)

	u, err := url.Parse(cfg.RedirectURI)
	if err != nil {
		return fmt.Errorf("invalid redirect-uri: %w", err)
	}

	// 1) Create a local mux so we don't stomp on the global DefaultServeMux
	var srv *http.Server
	mux := http.NewServeMux()
	mux.HandleFunc(u.Path, func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code == "" {
			http.Error(w, "code not found in query", http.StatusBadRequest)
			return
		}

		// inject code & verifier into a copy of cfg
		cfg2 := cfg
		cfg2.Code = code
		cfg2.CodeVerifier = verifier

		token, err := pkceExchange(cfg2)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, `<html><body><h1>Token received</h1><pre>%s</pre></body></html>`, token)

		// shutdown after responding
		go shutdownServer(srv)
	})

	// 2) Use that mux as the server's handler
	srv = &http.Server{
		Addr:    u.Host,
		Handler: mux,
	}

	logrus.Infof("1) Open this URL in your browser:\n\n  %s\n", authURL)
	logrus.Infof("2) Listening on %s and waiting for callback", u.Host)

	// 3) Listen, but treat ErrServerClosed as a successful exit
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func pkceExchange(cfg config.Token_t) (string, error) {
	form := url.Values{"grant_type": {"authorization_code"}}
	form.Set("client_id", cfg.ClientID)
	if cfg.ClientSecret != "" {
		form.Set("client_secret", cfg.ClientSecret)
	}
	form.Set("code", cfg.Code)
	form.Set("redirect_uri", cfg.RedirectURI)
	form.Set("code_verifier", cfg.CodeVerifier)

	return postForToken(cfg.JWTRealmURL, form)
}

func ropc(cfg config.Token_t) (string, error) {
	form := url.Values{
		"grant_type": {"password"},
		"client_id":  {cfg.ClientID},
		"username":   {cfg.Username},
		"password":   {cfg.Password},
	}
	if cfg.ClientSecret != "" {
		form.Set("client_secret", cfg.ClientSecret)
	}
	return postForToken(cfg.JWTRealmURL, form)
}

func clientCredentials(cfg config.Token_t) (string, error) {
	form := url.Values{
		"grant_type": {"client_credentials"},
		"client_id":  {cfg.ClientID},
	}
	if cfg.ClientSecret != "" {
		form.Set("client_secret", cfg.ClientSecret)
	}
	return postForToken(cfg.JWTRealmURL, form)
}

func postForToken(realmURL string, form url.Values) (string, error) {
	endpoint := strings.TrimRight(realmURL, "/") + "/protocol/openid-connect/token"
	resp, err := http.PostForm(endpoint, form)
	if err != nil {
		return "", fmt.Errorf("token request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		return "", fmt.Errorf("token endpoint returned %d: %s", resp.StatusCode, buf.String())
	}

	var out map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", fmt.Errorf("invalid JSON: %w", err)
	}

	accessToken, ok := out["access_token"].(string)
	if !ok {
		return "", fmt.Errorf("access_token not found or not a string")
	}
	return accessToken, nil
}

func runLocalFlow(cfg config.Token_t) error {
	if cfg.JWTPrivateKey == "" {
		return errors.New("--private-key is required for local mode")
	}
	if cfg.UserID == "" {
		return errors.New("--user is required for local mode")
	}
	if cfg.Roles == "" {
		return errors.New("--roles is required for local mode")
	}

	data, err := os.ReadFile(cfg.JWTPrivateKey)
	if err != nil {
		return fmt.Errorf("read private key: %w", err)
	}
	block, _ := pem.Decode(data)
	if block == nil {
		return errors.New("invalid PEM data")
	}
	keyAny, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return fmt.Errorf("parse private key: %w", err)
	}
	priv, ok := keyAny.(*rsa.PrivateKey)
	if !ok {
		return errors.New("not RSA private key")
	}

	now := time.Now()
	claims := jwt.MapClaims{
		"sub":   cfg.UserID,
		"roles": strings.Split(cfg.Roles, ","),
		"exp":   now.Add(cfg.TTL).Unix(),
		"iat":   now.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signed, err := token.SignedString(priv)
	if err != nil {
		return fmt.Errorf("sign token: %w", err)
	}
	fmt.Println(signed)
	return nil
}

// shutdownServer gracefully stops srv after callback completes.
func shutdownServer(srv *http.Server) {
	logrus.Infof("3) shutting down server on %s", srv.Addr)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
}

func newCodeVerifier() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func codeChallenge(verifier string) string {
	sum := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(sum[:])
}
