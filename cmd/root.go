package cmd

import (
	"github.com/patrickdappollonio/vsc-replacer/processor"
	"github.com/spf13/cobra"
)

func mainCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "vsc-replacer",
		Short: "A tool to replace regex matches in files",
		RunE: func(cmd *cobra.Command, args []string) error {
			regexInput, _ := cmd.Flags().GetString("regex")
			replacement, _ := cmd.Flags().GetString("replacement")
			dir, _ := cmd.Flags().GetString("dir")
			dryRun, _ := cmd.Flags().GetBool("dry-run")

			return processor.ReplaceFiles(processor.ReplaceOptions{
				Regex:       regexInput,
				Replacement: replacement,
				Dir:         dir,
				DryRun:      dryRun,
			})
		},
	}

	// Add flags to the root command
	rootCmd.PersistentFlags().String("regex", "", "Regular expression with capture groups")
	rootCmd.PersistentFlags().String("replacement", "", "Replacement string")
	rootCmd.PersistentFlags().String("dir", "", "Directory with files")
	rootCmd.PersistentFlags().Bool("dry-run", false, "Dry run mode")
	rootCmd.MarkPersistentFlagRequired("regex")
	rootCmd.MarkPersistentFlagRequired("replacement")
	rootCmd.MarkPersistentFlagRequired("dir")

	// Add subcommands
	rootCmd.AddCommand(grepCommand())

	return rootCmd
}

// MainCommand returns the configured root command
func MainCommand() *cobra.Command {
	return mainCommand()
}
