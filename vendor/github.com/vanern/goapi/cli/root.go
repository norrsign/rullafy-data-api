package cli

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/vanern/goapi/internal/cmd"
	"github.com/vanern/goapi/types"
)

func GetRootCmd() *cobra.Command {
	return cmd.RootCmd
}

func AddRunHook(name types.Hook, hook types.ConfigHook) {
	cmd.AddRunHook(name, hook)
}

func Run() error {
	return cmd.Execute()
}

func init() {
	if err := godotenv.Load(); err == nil {
		logrus.Infof("loaded environment variables from .env file")
	}

	if mode := os.Getenv("GOAPI_RUN_MODE"); mode != "" {
		envFile := fmt.Sprintf(".env.%s", mode)
		if err := godotenv.Overload(envFile); err != nil {
			logrus.Warnf("could not load %q: %v", envFile, err)
		} else {
			logrus.Infof("loaded environment variables from %q", envFile)
		}
	}

}
