package cmd

import (
	"github.com/patrickdappollonio/vsc-replacer/processor"
	"github.com/spf13/cobra"
)

func grepCommand() *cobra.Command {
	grepCmd := &cobra.Command{
		Use:   "grep",
		Short: "Search for regex matches in files without replacing them",
		RunE: func(cmd *cobra.Command, args []string) error {
			regexInput, _ := cmd.Flags().GetString("regex")
			dir, _ := cmd.Flags().GetString("dir")

			return processor.GrepFiles(processor.GrepOptions{
				Regex: regexInput,
				Dir:   dir,
			})
		},
	}

	// Add flags to the grep command (same as root but replacement is optional)
	grepCmd.Flags().String("regex", "", "Regular expression with capture groups")
	grepCmd.Flags().String("replacement", "", "Replacement string (ignored in grep mode)")
	grepCmd.Flags().String("dir", "", "Directory with files")
	grepCmd.Flags().Bool("dry-run", false, "Dry run mode (ignored in grep mode)")
	grepCmd.MarkFlagRequired("regex")
	grepCmd.MarkFlagRequired("dir")

	return grepCmd
}
