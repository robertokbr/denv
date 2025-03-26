/*
* I don't know how does it works but it works and was completly made by AI :P
 */

package cli

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/robertokbr/denv/bucket"
	"github.com/robertokbr/denv/config"
)

const (
	// ZSH completion script
	completionScript = `#compdef denv

_denv_files() {
  local files
  IFS=' ' read -A files <<< "$(denv --completion-files)"
  _describe -t files "files" files
}

_denv() {
  local curcontext="$curcontext" state line
  typeset -A opt_args
  
  _arguments \
    '--help[See how to use the CLI]' \
    '--config[Start the app config]' \
    '--up[Upload some file]:file:_files' \
    '--name[Nickname for the file]:file:_denv_files' \
    '--out[Optional output name]:filename:' \
    '--list[List all files in the bucket]' \
    '--del[Delete some file in the bucket]:file:_denv_files' \
    '--rename[Rename a file in the bucket]:file:_denv_files' \
    '--setup-completion[Setup shell completion for denv commands]'
}

compdef _denv denv
`
)

// GenerateCompletionScript generates the zsh completion script content
func GenerateCompletionScript() string {
	return completionScript
}

// WriteCompletionScript writes the completion script to the user's zsh completion directory
func WriteCompletionScript() error {
	// Create completion directory if it doesn't exist
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("could not get user home directory: %v", err)
	}

	var completionPath string
	// For ZSH, the completion functions should go in a directory in the fpath
	// Common locations include ~/.zsh/functions or ~/.zsh/completion
	zshFunctionsDir := path.Join(homeDir, ".zsh", "functions")
	if _, err := os.Stat(zshFunctionsDir); err == nil {
		completionPath = path.Join(zshFunctionsDir, "_denv")
	} else {
		// Create directory if it doesn't exist
		err = os.MkdirAll(zshFunctionsDir, 0755)
		if err != nil {
			return fmt.Errorf("could not create zsh functions directory: %v", err)
		}
		completionPath = path.Join(zshFunctionsDir, "_denv")
	}

	// Write completion script
	err = os.WriteFile(completionPath, []byte(completionScript), 0644)
	if err != nil {
		return fmt.Errorf("could not write completion script: %v", err)
	}

	// Check if zshrc exists
	zshrcPath := path.Join(homeDir, ".zshrc")
	if _, err := os.Stat(zshrcPath); err == nil {
		// Check if fpath is already set in zshrc
		zshrcContent, err := os.ReadFile(zshrcPath)
		if err != nil {
			return fmt.Errorf("could not read .zshrc file: %v", err)
		}

		// Only add fpath entry if it isn't already there
		if !strings.Contains(string(zshrcContent), zshFunctionsDir) {
			// Append fpath entry to .zshrc
			f, err := os.OpenFile(zshrcPath, os.O_APPEND|os.O_WRONLY, 0644)
			if err != nil {
				return fmt.Errorf("could not open .zshrc for appending: %v", err)
			}
			defer f.Close()

			// Add the functions directory to fpath and enable compinit
			fpathUpdate := fmt.Sprintf("\n# Added by denv for completion\nfpath=(%s $fpath)\nautoload -Uz compinit\ncompinit\n", zshFunctionsDir)
			if _, err := f.WriteString(fpathUpdate); err != nil {
				return fmt.Errorf("could not update .zshrc: %v", err)
			}
		}
	}

	fmt.Println("🎉 ZSH completion has been set up successfully!")
	fmt.Println("ℹ️  You need to restart your shell or run 'source ~/.zshrc' to enable it.")

	return nil
}

// GetFileList returns a list of files from the S3 bucket for completion
func GetFileList() ([]string, error) {
	// Initialize config
	if err := config.InitPaths(); err != nil {
		return nil, err
	}

	// Setup environment
	if err := config.SetupEnvironment(); err != nil {
		return nil, err
	}

	// Check if environment is properly configured
	err := config.ValidateEnvironment()
	if err != nil {
		return nil, err
	}

	// Get credentials
	creds := config.GetAWSCredentials()

	// Create S3 bucket client
	s3Client := bucket.NewS3Bucket(
		creds.AccessKey,
		creds.SecretKey,
		creds.BucketName,
		creds.BucketRegion,
	)

	// Get the list of files
	return s3Client.ListFileNames()
}

// PrintFileList prints the list of files for command completion
func PrintFileList() {
	files, err := GetFileList()
	if err != nil {
		return // Silent fail for completion
	}

	fmt.Println(strings.Join(files, " "))
}
