package cli

import (
	"fmt"
)

type Command struct {
	Name        string
	Description string
	Execute     func() error
}

func printCommandError(format string, args ...interface{}) error {
	errorMsg := fmt.Sprintf(format, args...)
	fmt.Println(errorMsg)
	return fmt.Errorf(errorMsg)
}

func newUploadCommand(cli *CLI) Command {
	return Command{
		Name:        "upload",
		Description: "Upload a file to S3 bucket",
		Execute: func() error {
			if cli.flagName == "" {
				return printCommandError("🌝 Please, provide a nickname to your file using --name flag")
			}
			
			cli.handleUpload()
			return nil
		},
	}
}

func newDownloadCommand(cli *CLI) Command {
	return Command{
		Name:        "download",
		Description: "Download a file from S3 bucket",
		Execute: func() error {
			cli.handleDownload()
			return nil
		},
	}
}

func newListCommand(cli *CLI) Command {
	return Command{
		Name:        "list",
		Description: "List all files in S3 bucket",
		Execute: func() error {
			cli.handleList()
			return nil
		},
	}
}

func newDeleteCommand(cli *CLI) Command {
	return Command{
		Name:        "delete",
		Description: "Delete a file from S3 bucket",
		Execute: func() error {
			cli.handleDelete()
			return nil
		},
	}
}

func newRenameCommand(cli *CLI) Command {
	return Command{
		Name:        "rename",
		Description: "Rename a file in S3 bucket",
		Execute: func() error {
			if cli.flagName == "" {
				return printCommandError("🌝 Please, provide a new name for the file using --name flag")
			}
			cli.handleRename()
			return nil
		},
	}
}

func newConfigCommand(cli *CLI) Command {
	return Command{
		Name:        "config",
		Description: "Configure the application",
		Execute: func() error {
			cli.handleConfig()
			return nil
		},
	}
}

func newCompletionFilesCommand(cli *CLI) Command {
	return Command{
		Name:        "completion-files",
		Description: "List files for shell completion",
		Execute: func() error {
			cli.handleCompletionFiles()
			return nil
		},
	}
}

func newSetupCompletionCommand(cli *CLI) Command {
	return Command{
		Name:        "setup-completion",
		Description: "Setup shell completion for denv commands",
		Execute: func() error {
			cli.handleSetupCompletion()
			return nil
		},
	}
}

func newHelpCommand(cli *CLI) Command {
	return Command{
		Name:        "help",
		Description: "Show help information",
		Execute: func() error {
			cli.handleHelp()
			return nil
		},
	}
}