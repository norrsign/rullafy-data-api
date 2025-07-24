// internal/cmd/completion.go
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// completionCmd generates shell-completion scripts for Cobra.
var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate shell completion script",
	Long: `To load completions:

Bash:

  $ source <(goapi completion bash)

  # To load completions for each session, add to ~/.bashrc:
  #   source <(goapi completion bash)

Zsh:

  # If shell completion isn't already enabled in your environment you will need
  # to enable it. You can execute once:
  #   echo "autoload -U compinit; compinit" >> ~/.zshrc
  #
  # To load completions for each session, add to ~/.zshrc:
  #   source <(goapi completion zsh)

Fish:

  $ goapi completion fish | source

  # To load completions for each session, add to ~/.config/fish/completions/goapi.fish:
  #   goapi completion fish > ~/.config/fish/completions/goapi.fish

PowerShell:

  PS> goapi completion powershell > goapi.ps1
  # Then source it from your PowerShell profile.`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.ExactValidArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		switch args[0] {
		case "bash":
			return RootCmd.GenBashCompletion(os.Stdout)
		case "zsh":
			return RootCmd.GenZshCompletion(os.Stdout)
		case "fish":
			return RootCmd.GenFishCompletion(os.Stdout, true)
		case "powershell":
			return RootCmd.GenPowerShellCompletionWithDesc(os.Stdout)
		}
		return nil
	},
}

func init() {
	RootCmd.AddCommand(completionCmd)
}
