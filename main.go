package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/spf13/cobra"
)

func lineDiff(original, replaced []string, filename string) []string {
	var diffLines []string
	dmp := diffmatchpatch.New()

	for i := range original {
		if original[i] != replaced[i] {
			diffs := dmp.DiffMain(original[i], replaced[i], false)
			diffLines = append(diffLines, fmt.Sprintf("--- a/%s:%d\n+++ b/%s:%d\n%s", filename, i+1, filename, i+1, dmp.DiffPrettyText(diffs)))
		}
	}

	return diffLines
}

var grepCmd = &cobra.Command{
	Use:   "grep",
	Short: "Search for regex matches in files without replacing them",
	RunE: func(cmd *cobra.Command, args []string) error {
		regexInput, _ := cmd.Flags().GetString("regex")
		dir, _ := cmd.Flags().GetString("dir")

		regex, err := regexp.Compile(regexInput)
		if err != nil {
			return fmt.Errorf("failed to compile regular expression: %w", err)
		}

		err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// Process only regular files (not directories)
			if !info.IsDir() {
				content, err := os.ReadFile(path)
				if err != nil {
					return fmt.Errorf("failed to read file %s: %w", path, err)
				}

				scanner := bufio.NewScanner(strings.NewReader(string(content)))
				lineNum := 0
				for scanner.Scan() {
					lineNum++
					line := scanner.Text()
					if regex.MatchString(line) {
						fmt.Printf("%s:%d:%s\n", path, lineNum, line)
					}
				}
			}

			return nil
		})
		if err != nil {
			return err
		}

		return nil
	},
}

var rootCmd = &cobra.Command{
	Use:   "vsc-replacer",
	Short: "A tool to replace regex matches in files",
	RunE: func(cmd *cobra.Command, args []string) error {
		regexInput, _ := cmd.Flags().GetString("regex")
		replacement, _ := cmd.Flags().GetString("replacement")
		dir, _ := cmd.Flags().GetString("dir")
		dryRun, _ := cmd.Flags().GetBool("dry-run")

		regex, err := regexp.Compile(regexInput)
		if err != nil {
			return fmt.Errorf("failed to compile regular expression: %w", err)
		}

		err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// Process only regular files (not directories)
			if !info.IsDir() {
				content, err := os.ReadFile(path)
				if err != nil {
					return fmt.Errorf("failed to read file %s: %w", path, err)
				}

				replaced := regex.ReplaceAllString(string(content), replacement)

				if dryRun {
					scannerO := bufio.NewScanner(strings.NewReader(string(content)))
					scannerR := bufio.NewScanner(strings.NewReader(replaced))

					var originalLines, replacedLines []string
					for scannerO.Scan() {
						originalLines = append(originalLines, scannerO.Text())
					}
					for scannerR.Scan() {
						replacedLines = append(replacedLines, scannerR.Text())
					}

					diffLines := lineDiff(originalLines, replacedLines, path)
					for _, diff := range diffLines {
						fmt.Println(diff)
						fmt.Println(strings.Repeat("-", 3))
					}
				} else {
					err = os.WriteFile(path, []byte(replaced), info.Mode())
					if err != nil {
						return fmt.Errorf("failed to write to file %s: %w", path, err)
					}
				}
			}

			return nil
		})
		if err != nil {
			return err
		}

		return nil
	},
}

func main() {
	// Add flags to the root command
	rootCmd.PersistentFlags().String("regex", "", "Regular expression with capture groups")
	rootCmd.PersistentFlags().String("replacement", "", "Replacement string")
	rootCmd.PersistentFlags().String("dir", "", "Directory with files")
	rootCmd.PersistentFlags().Bool("dry-run", false, "Dry run mode")
	rootCmd.MarkPersistentFlagRequired("regex")
	rootCmd.MarkPersistentFlagRequired("replacement")
	rootCmd.MarkPersistentFlagRequired("dir")

	// Add flags to the grep command (same as root but replacement is optional)
	grepCmd.Flags().String("regex", "", "Regular expression with capture groups")
	grepCmd.Flags().String("replacement", "", "Replacement string (ignored in grep mode)")
	grepCmd.Flags().String("dir", "", "Directory with files")
	grepCmd.Flags().Bool("dry-run", false, "Dry run mode (ignored in grep mode)")
	grepCmd.MarkFlagRequired("regex")
	grepCmd.MarkFlagRequired("dir")

	// Add the grep subcommand to the root command
	rootCmd.AddCommand(grepCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
