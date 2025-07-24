package config

import "time"

// Token_t holds BOTH local (self‑signed) and remote (Keycloak) token settings.
type Token_t struct {
	// ─── Local mode (self-signed) ────────────────────────────────────────
	JWTPrivateKey string        // path to RSA PKCS#8 private key
	UserID        string        // "sub" claim
	Roles         string        // comma-separated roles
	TTL           time.Duration // token lifetime

	// ─── Remote mode (Keycloak / OAuth2) ────────────────────────────────
	JWTRealmURL  string // Keycloak realm URL
	ClientID     string // OAuth2 client id
	ClientSecret string // OAuth2 client secret (if confidential)

	// ROPC (password grant)
	Username string
	Password string

	// PKCE (auth code)
	Code         string
	CodeVerifier string
	RedirectURI  string
}

type Start_t struct {
	Verbose bool
}

type Server_t struct {
	Start                 Start_t
	JWTPublicKey          string
	JWTRealmURL           string
	JWTKeyRefreshInterval time.Duration
}

type Config_t struct {
	Token  Token_t
	Server Server_t
}

var Config = Config_t{
	Server: Server_t{
		Start: Start_t{
			Verbose: false,
		},
		JWTPublicKey:          "",
		JWTRealmURL:           "",
		JWTKeyRefreshInterval: time.Minute, // default to 1m
	},
	Token: Token_t{
		// local defaults
		JWTPrivateKey: "",
		UserID:        "",
		Roles:         "",
		TTL:           time.Hour,

		// remote defaults
		JWTRealmURL:  "",
		ClientID:     "",
		ClientSecret: "",
		Username:     "",
		Password:     "",
		Code:         "",
		CodeVerifier: "",
		RedirectURI:  "",
	},
}
