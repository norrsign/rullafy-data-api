// internal/cmd/token.go
package cmd

import (
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/vanern/goapi/config"
	"github.com/vanern/goapi/internal/utils"
	"github.com/vanern/goapi/types"
)

var tokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Fetch a token from Keycloak (PKCE, ROPC, Client Credentials) or generate one locally",
	RunE:  runToken,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		afterConfigHook := hooks[types.TokenAfterConfigHook]
		if afterConfigHook != nil {
			return afterConfigHook()
		}
		return nil
	},
}

func init() {
	RootCmd.AddCommand(tokenCmd)

	// Keycloak / remote‑mode flags (write into config.Config.Token.* via viper keys)
	tokenCmd.Flags().String("jwt-realm-url", "", "Keycloak realm URL (e.g. https://auth.example.com/realms/your-realm)")
	viper.BindPFlag("jwt_realm_url", tokenCmd.Flags().Lookup("jwt-realm-url"))

	tokenCmd.Flags().String("client-id", "", "Keycloak client ID")
	viper.BindPFlag("client_id", tokenCmd.Flags().Lookup("client-id"))

	tokenCmd.Flags().String("client-secret", "", "Keycloak client secret (if confidential)")
	viper.BindPFlag("client_secret", tokenCmd.Flags().Lookup("client-secret"))

	tokenCmd.Flags().String("username", "", "Username for ROPC flow")
	viper.BindPFlag("username", tokenCmd.Flags().Lookup("username"))

	tokenCmd.Flags().String("password", "", "Password for ROPC flow")
	viper.BindPFlag("password", tokenCmd.Flags().Lookup("password"))

	tokenCmd.Flags().String("code", "", "Authorization code from browser (PKCE flow)")
	viper.BindPFlag("code", tokenCmd.Flags().Lookup("code"))

	tokenCmd.Flags().String("code-verifier", "", "PKCE code verifier (from initial step)")
	viper.BindPFlag("code_verifier", tokenCmd.Flags().Lookup("code-verifier"))

	tokenCmd.Flags().String("redirect-uri", "http://localhost:8085/callback", "Redirect URI for PKCE flow")
	viper.BindPFlag("redirect_uri", tokenCmd.Flags().Lookup("redirect-uri"))

	// Local‑mode flags
	tokenCmd.Flags().String("private-key", "", "Path to RSA PKCS#8 private key (local mode)")
	viper.BindPFlag("jwt_private_key", tokenCmd.Flags().Lookup("private-key"))

	tokenCmd.Flags().String("user", "", "Subject (sub) claim / user ID (local mode)")
	viper.BindPFlag("token_user", tokenCmd.Flags().Lookup("user"))

	tokenCmd.Flags().String("roles", "", "Comma-separated roles claim (local mode)")
	viper.BindPFlag("token_roles", tokenCmd.Flags().Lookup("roles"))

	tokenCmd.Flags().Duration("ttl", time.Hour, "Token lifetime (local mode)")
	viper.BindPFlag("token_ttl", tokenCmd.Flags().Lookup("ttl"))
}

func runToken(cmd *cobra.Command, args []string) error {
	// Config has already been loaded in RootCmd.PersistentPreRun → initConfig.
	// Just pass it down.
	return utils.RunTokenFlows(config.Config.Token)
}
