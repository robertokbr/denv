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
	completionScript = `
		_denv_completion() {
			local cur prev opts
			COMPREPLY=()
			cur="${COMP_WORDS[COMP_CWORD]}"
			prev="${COMP_WORDS[COMP_CWORD-1]}"
			
			# List of options
			opts="--help --config --up --name --out --list --del --rename"
			
			# Check if we need file completion
			case "${prev}" in
				--del|--rename)
					# Get the list of files from bucket
					files=$(denv --completion-files)
					COMPREPLY=( $(compgen -W "${files}" -- ${cur}) )
					return 0
					;;
				--name)
					# Check if previous command was --rename to complete file names
					if [[ "${COMP_WORDS[COMP_CWORD-2]}" == "--rename" ]]; then
						files=$(denv --completion-files)
						COMPREPLY=( $(compgen -W "${files}" -- ${cur}) )
						return 0
					fi
					return 0
					;;
				--up)
					# Complete with local files
					COMPREPLY=( $(compgen -f ${cur}) )
					return 0
					;;
				*)
					;;
			esac
			
			# Complete options
			if [[ ${cur} == -* ]] ; then
				COMPREPLY=( $(compgen -W "${opts}" -- ${cur}) )
				return 0
			fi
		}

		complete -F _denv_completion denv
	`
)

// GenerateCompletionScript generates the bash completion script content
func GenerateCompletionScript() string {
	return completionScript
}

// WriteCompletionScript writes the completion script to the user's bash completion directory
func WriteCompletionScript() error {
	// Create completion directory if it doesn't exist
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("could not get user home directory: %v", err)
	}

	var completionPath string
	// Check if we're using bash_completion.d or completions directory
	// macOS usually uses /usr/local/etc/bash_completion.d/
	// Linux often uses /etc/bash_completion.d/ or ~/.local/share/bash-completion/completions/

	// First try user's home directory
	bashCompletionDir := path.Join(homeDir, ".bash_completion.d")
	if _, err := os.Stat(bashCompletionDir); err == nil {
		completionPath = path.Join(bashCompletionDir, "denv")
	} else {
		// Create directory if it doesn't exist
		err = os.MkdirAll(bashCompletionDir, 0755)
		if err != nil {
			return fmt.Errorf("could not create bash completion directory: %v", err)
		}
		completionPath = path.Join(bashCompletionDir, "denv")
	}

	// Write completion script
	err = os.WriteFile(completionPath, []byte(completionScript), 0644)
	if err != nil {
		return fmt.Errorf("could not write completion script: %v", err)
	}

	// Create .bash_completion file in home directory if it doesn't exist
	bashCompletionFile := path.Join(homeDir, ".bash_completion")
	if _, err := os.Stat(bashCompletionFile); os.IsNotExist(err) {
		content := "for bcfile in ~/.bash_completion.d/* ; do\n  [ -f \"$bcfile\" ] && . $bcfile\ndone\n"
		err = os.WriteFile(bashCompletionFile, []byte(content), 0644)
		if err != nil {
			return fmt.Errorf("could not write .bash_completion file: %v", err)
		}
	}

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
