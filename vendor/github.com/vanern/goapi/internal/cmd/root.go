// internal/cmd/root.go
package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/vanern/goapi/config"
	"github.com/vanern/goapi/types"
)

var hooks = make(map[types.Hook]types.ConfigHook)

func AddRunHook(name types.Hook, hook types.ConfigHook) {
	if _, exists := hooks[name]; exists {
		logrus.Warnf("Hook %q already exists, overwriting", name)
	}
	hooks[name] = hook
}

// RootCmd is the base command for goapi.
var RootCmd = &cobra.Command{
	Use:   "goapi",
	Short: "goapi is the API server application",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// 1) Load config (flags > env > file > defaults)
		initConfig()

		// 2) Configure log level
		level := viper.GetString("log_level")
		lvl, err := logrus.ParseLevel(level)
		if err != nil {
			logrus.Fatalf("invalid log level %q: %v", level, err)
		}
		logrus.SetLevel(lvl)

		// 3) Run the hook, if provided
		afterConfigHook := hooks[types.GlobalAfterConfigHook]
		if afterConfigHook != nil {
			afterConfigHook()
		}

	},
}

// Execute runs the root command.
func Execute() error {

	return RootCmd.Execute()
}

func init() {
	// Tell Viper to read ENV variables (GOAPI_*) and bind them automatically.
	viper.SetEnvPrefix("goapi")
	viper.AutomaticEnv()

	// Define the log‐level flag and bind to viper
	RootCmd.PersistentFlags().
		String("log-level", "info", "log level (debug, info, warn, error)")
	viper.BindPFlag("log_level", RootCmd.PersistentFlags().Lookup("log-level"))
	viper.SetDefault("log_level", "info")
}

// initConfig reads configuration from (in order):
//  1. Cobra flags
//  2. Environment variables (GOAPI_*)
//     2.1. A mode‐specific .env file if GOAPI_RUN_MODE is set
//  3. Optional config file (config.{json,yaml,toml} in cwd)
//  4. Built‐in defaults
func initConfig() {
	// 3) Optional config file

	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err == nil {
		logrus.Infof("Using config file: %s", viper.ConfigFileUsed())
	}

	// 4) Populate your global Config struct

	// server
	config.Config.Server.Start.Verbose = viper.GetBool("start_verbose")
	config.Config.Server.JWTPublicKey = viper.GetString("jwt_public_key")
	config.Config.Server.JWTRealmURL = viper.GetString("jwt_realm_url")
	config.Config.Server.JWTKeyRefreshInterval = viper.GetDuration("jwt_key_refresh_interval")

	// token (local)
	config.Config.Token.JWTPrivateKey = viper.GetString("jwt_private_key")
	config.Config.Token.UserID = viper.GetString("token_user")
	config.Config.Token.Roles = viper.GetString("token_roles")
	config.Config.Token.TTL = viper.GetDuration("token_ttl")

	// token (remote / oauth2)
	config.Config.Token.JWTRealmURL = viper.GetString("jwt_realm_url") // same key as server if you want
	config.Config.Token.ClientID = viper.GetString("client_id")
	config.Config.Token.ClientSecret = viper.GetString("client_secret")
	config.Config.Token.Username = viper.GetString("username")
	config.Config.Token.Password = viper.GetString("password")
	config.Config.Token.Code = viper.GetString("code")
	config.Config.Token.CodeVerifier = viper.GetString("code_verifier")
	config.Config.Token.RedirectURI = viper.GetString("redirect_uri")
}
