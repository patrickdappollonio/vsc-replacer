package processor

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/sergi/go-diff/diffmatchpatch"
)

// GrepOptions contains configuration for grep operations
type GrepOptions struct {
	Regex string
	Dir   string
}

// ReplaceOptions contains configuration for replacement operations
type ReplaceOptions struct {
	Regex       string
	Replacement string
	Dir         string
	DryRun      bool
}

// GrepFiles searches for regex matches in files and prints them
func GrepFiles(opts GrepOptions) error {
	regex, err := regexp.Compile(opts.Regex)
	if err != nil {
		return fmt.Errorf("failed to compile regular expression: %w", err)
	}

	return filepath.Walk(opts.Dir, func(path string, info os.FileInfo, err error) error {
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
			if err := scanner.Err(); err != nil {
				return fmt.Errorf("failed to scan file %s: %w", path, err)
			}
		}

		return nil
	})
}

// ReplaceFiles performs regex replacement in files
func ReplaceFiles(opts ReplaceOptions) error {
	regex, err := regexp.Compile(opts.Regex)
	if err != nil {
		return fmt.Errorf("failed to compile regular expression: %w", err)
	}

	return filepath.Walk(opts.Dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Process only regular files (not directories)
		if !info.IsDir() {
			content, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("failed to read file %s: %w", path, err)
			}

			replaced := regex.ReplaceAllString(string(content), opts.Replacement)

			if opts.DryRun {
				return showDiff(string(content), replaced, path)
			} else {
				err = os.WriteFile(path, []byte(replaced), info.Mode())
				if err != nil {
					return fmt.Errorf("failed to write to file %s: %w", path, err)
				}
			}
		}

		return nil
	})
}

// showDiff displays the differences between original and replaced content
func showDiff(original, replaced, filename string) error {
	scannerO := bufio.NewScanner(strings.NewReader(original))
	scannerR := bufio.NewScanner(strings.NewReader(replaced))

	var originalLines, replacedLines []string
	for scannerO.Scan() {
		originalLines = append(originalLines, scannerO.Text())
	}
	if err := scannerO.Err(); err != nil {
		return fmt.Errorf("failed to scan original content for file %s: %w", filename, err)
	}
	for scannerR.Scan() {
		replacedLines = append(replacedLines, scannerR.Text())
	}
	if err := scannerR.Err(); err != nil {
		return fmt.Errorf("failed to scan replaced content for file %s: %w", filename, err)
	}

	diffLines := generateLineDiff(originalLines, replacedLines, filename)
	for _, diff := range diffLines {
		fmt.Println(diff)
		fmt.Println(strings.Repeat("-", 3))
	}

	return nil
}

// generateLineDiff creates a diff between original and replaced lines
func generateLineDiff(original, replaced []string, filename string) []string {
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
