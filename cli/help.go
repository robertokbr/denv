package cli

import (
	"fmt"
)

func PrintHelp() {
	fmt.Println("denv --config to start the CLI configuration")
	fmt.Println("denv --up [file path] --name [file nickname] to upload some env file")
	fmt.Println("denv --name [file nickname] to download some env file you have uploaded")
	fmt.Println("denv --name [file nickname] --out [file name] to download some env file you have uploaded with some specific name")
	fmt.Println("denv --list to list all files in the bucket")
	fmt.Println("denv --del [file nickname] to delete some file in the bucket")
	fmt.Println("denv --rename [file nickname] --name [new nickname] to rename a file in the bucket")
	fmt.Println("denv --setup-completion to install tab completion for commands (bash)")
}

func PrintSetupMessage() {
	fmt.Println("🤔 Hello! Type 'denv --setup' to start setting up the application")
}

func PrintSuccessConfig() {
	fmt.Println("🔥 Thank you! Everything is right!")
	fmt.Println("🤓 Type denv --help if you want to see how to use the CLI.")
	fmt.Println("🫢 Type denv --config again if the CLI is not working properly.")
}