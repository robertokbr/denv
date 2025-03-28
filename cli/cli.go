package cli

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/robertokbr/denv/bucket"
	"github.com/robertokbr/denv/config"
)

type CLI struct {
	s3bucket            *bucket.S3Bucket
	flagUpload          string
	flagName            string
	flagOutput          string
	flagList            bool
	flagDelete          string
	flagHelp            bool
	flagConfig          bool
	flagRename          string
	flagCompletionFiles bool
	flagSetupCompletion bool
	flagRecursive       bool
	commands            map[string]Command
}

func New() *CLI {
	cli := &CLI{
		commands: make(map[string]Command),
	}

	flag.BoolVar(&cli.flagHelp, "help", false, "See how to use the CLI")
	flag.BoolVar(&cli.flagConfig, "config", false, "Start the app config")
	flag.StringVar(&cli.flagUpload, "up", "", "Upload some file")
	flag.StringVar(&cli.flagName, "name", "", "Nickname to the file you will upload")
	flag.StringVar(&cli.flagOutput, "out", "", "Optional flag to specify the output such as: .env.example")
	flag.BoolVar(&cli.flagList, "list", false, "List all files in the bucket")
	flag.StringVar(&cli.flagDelete, "del", "", "Delete some file in the bucket")
	flag.StringVar(&cli.flagRename, "rename", "", "Rename a file in the bucket")
	flag.BoolVar(&cli.flagCompletionFiles, "completion-files", false, "List files for shell completion (internal use)")
	flag.BoolVar(&cli.flagSetupCompletion, "setup-completion", false, "Setup shell completion for denv commands")
	flag.BoolVar(&cli.flagRecursive, "r", false, "Upload a directory recursively (will be zipped)")

	flag.Parse()

	// Register commands first so we can handle special commands
	cli.registerCommands()

	// Skip initialization for completion-related commands
	if !cli.flagCompletionFiles && !cli.flagSetupCompletion {
		// Initialize configuration
		err := initializeApp()
		if err != nil {
			log.Fatalf("Failed to initialize application: %v", err)
		}

		// Create S3 bucket instance
		cli.initializeS3Bucket()
	}

	return cli
}

func initializeApp() error {
	// Initialize paths
	if err := config.InitPaths(); err != nil {
		return fmt.Errorf("failed to initialize paths: %v", err)
	}

	// Setup environment
	if err := config.SetupEnvironment(); err != nil {
		return fmt.Errorf("failed to setup environment: %v", err)
	}

	return nil
}

func (cli *CLI) initializeS3Bucket() {
	creds := config.GetAWSCredentials()
	cli.s3bucket = bucket.NewS3Bucket(
		creds.AccessKey,
		creds.SecretKey,
		creds.BucketName,
		creds.BucketRegion,
	)
}

func (cli *CLI) registerCommands() {
	commands := []Command{
		newHelpCommand(cli),
		newConfigCommand(cli),
		newUploadCommand(cli),
		newDownloadCommand(cli),
		newListCommand(cli),
		newDeleteCommand(cli),
		newRenameCommand(cli),
		newCompletionFilesCommand(cli),
		newSetupCompletionCommand(cli),
	}

	for _, cmd := range commands {
		cli.commands[cmd.Name] = cmd
	}
}

func (cli *CLI) validateEnvironment() bool {
	err := config.ValidateEnvironment()
	if err != nil {
		PrintSetupMessage()
		return false
	}
	return true
}

func (cli *CLI) executeWithValidation(operation func()) {
	if !cli.validateEnvironment() {
		return
	}
	operation()
}

func (cli *CLI) handleUpload() {
	cli.executeWithValidation(func() {
		currentPath, err := os.Getwd()
		if err != nil {
			log.Fatalf("Failed to get the current path %v", err)
		}

		fullPath := path.Join(currentPath, cli.flagUpload)

		if cli.flagRecursive {
			// Create a temporary zip file with a unique name that doesn't conflict
			tempDir, err := os.MkdirTemp("", "denv")
			if err != nil {
				log.Fatalf("Failed to create temporary directory: %v", err)
			}
			defer os.RemoveAll(tempDir) // Clean up temp directory

			tempZipPath := path.Join(tempDir, "temp_archive")
			err = createZipArchive(fullPath, tempZipPath)
			if err != nil {
				log.Fatalf("Failed to create zip archive: %v", err)
			}

			// Check if the name already ends with .zip
			bucketName := cli.flagName
			if !strings.HasSuffix(bucketName, ".zip") {
				bucketName += ".zip"
			}

			// Upload the zip file
			cli.s3bucket.UploadFile(tempZipPath, bucketName)
		} else {
			// For regular files, preserve the original file extension if the user hasn't specified one
			originalExt := path.Ext(fullPath)
			targetName := cli.flagName

			// If the original file has an extension and the target name doesn't have any extension
			if originalExt != "" && path.Ext(targetName) == "" {
				targetName += originalExt
			}

			cli.s3bucket.UploadFile(fullPath, targetName)
		}
	})
}

func (cli *CLI) handleDownload() {
	cli.executeWithValidation(func() {
		outputPath := cli.flagOutput
		if outputPath == "" {
			outputPath = cli.flagName
		}

		// Download the file
		cli.s3bucket.DownloadFile(cli.flagName, outputPath)

		// Check if the file is a zip (ends with .zip)
		if strings.HasSuffix(outputPath, ".zip") {
			// Extract the zip file
			extractDir := strings.TrimSuffix(outputPath, ".zip")
			err := extractZipArchive(outputPath, extractDir)
			if err != nil {
				log.Printf("Warning: Failed to extract zip file: %v", err)
			} else {
				// Remove the zip file after extraction
				os.Remove(outputPath)
			}
		}
	})
}

func (cli *CLI) handleList() {
	cli.executeWithValidation(func() {
		cli.s3bucket.ListFiles()
	})
}

func (cli *CLI) handleDelete() {
	cli.executeWithValidation(func() {
		cli.s3bucket.DeleteFile(cli.flagDelete)
	})
}

func (cli *CLI) handleRename() {
	cli.executeWithValidation(func() {
		if cli.flagName == "" {
			fmt.Println("🌝 Please, provide a new name for the file using --name flag")
			return
		}
		cli.s3bucket.RenameFile(cli.flagRename, cli.flagName)
	})
}

func (cli *CLI) handleConfig() {
	ConfigureApplication()
}

func (cli *CLI) handleHelp() {
	PrintHelp()
}

func (cli *CLI) handleCompletionFiles() {
	PrintFileList()
}

func (cli *CLI) handleSetupCompletion() {
	err := WriteCompletionScript()
	if err != nil {
		fmt.Printf("Failed to setup completion: %v\n", err)
		return
	}
}

func (cli *CLI) executeCommand(name string) bool {
	if cmd, exists := cli.commands[name]; exists {
		cmd.Execute()
		return true
	}
	return false
}

func (cli *CLI) Run() {
	// Handle completion commands first as they don't require full initialization
	if cli.flagCompletionFiles && cli.executeCommand("completion-files") {
		return
	}

	if cli.flagSetupCompletion && cli.executeCommand("setup-completion") {
		return
	}

	if cli.flagHelp && cli.executeCommand("help") {
		return
	}

	if cli.flagConfig && cli.executeCommand("config") {
		return
	}

	if cli.flagUpload != "" && cli.flagName != "" && cli.executeCommand("upload") {
		return
	}

	if cli.flagRename != "" && cli.flagName != "" && cli.executeCommand("rename") {
		return
	}

	if cli.flagName != "" && cli.flagUpload == "" && cli.flagRename == "" && cli.executeCommand("download") {
		return
	}

	if cli.flagList && cli.executeCommand("list") {
		return
	}

	if cli.flagDelete != "" && cli.executeCommand("delete") {
		return
	}

	if cli.flagUpload != "" && cli.flagName == "" {
		fmt.Println("🌝 Please, provide a nickname to your file using --name flag")
		return
	}

	if cli.flagRename != "" && cli.flagName == "" {
		fmt.Println("🌝 Please, provide a new name for the file using --name flag")
		return
	}

	fmt.Println("🤓 Type denv --help if you want to see how to use the CLI.")
}
