package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strings"

	"github.com/robertokbr/denv/bucket"
	"github.com/robertokbr/denv/helpers"
)

func main() {
	helpers.StartProject()
	s3bucket := bucket.NewS3Bucket()

	help := flag.Bool("help", false, "See how to use the CLI")
	conf := flag.Bool("conf", false, "Start the app config")
	up := flag.String("up", "", "Upload some file")
	name := flag.String("name", "", "Nickname to the file you will upload")

	flag.Parse()

	if *help {
		fmt.Println("denv --conf to start the CLI configuration")
		fmt.Println("denv --up [file path] --name [file nickname] to upload some env file")
		fmt.Println("denv --name [file nickname] to download some env file you have uploaded")
	}

	if *conf {
		var AWS_ACCESS_KEY string
		var AWS_SECRET_KEY string
		var AWS_BUCKET_NAME string

		fmt.Println("🚧 Insert your AWS Access key")
		fmt.Scan(&AWS_ACCESS_KEY)
		fmt.Println("🚧 Insert your AWS Secret key")
		fmt.Scan(&AWS_SECRET_KEY)
		fmt.Println("🚧 Insert your AWS Bucket name")
		fmt.Scan(&AWS_BUCKET_NAME)
		fmt.Println("🔥 Thank you! Everything is right!")
		fmt.Println("🤓 Type denv --help if you want to see how to use the CLI.")
		fmt.Println("🫢 Type denv --conf again if the CLI is not working properly.")

		file, err := os.Create(helpers.ENV_PATH)

		if err != nil {
			log.Fatalf("Failed to write env file: %v", err)
		}

		data := "AWS_ACCESS_KEY=%s\nAWS_SECRET_KEY=%s\nAWS_BUCKET_NAME=%s\n"

		data = fmt.Sprintf(data, AWS_ACCESS_KEY, AWS_SECRET_KEY, AWS_BUCKET_NAME)

		sdata := strings.NewReader(data)

		_, err = io.Copy(file, sdata)

		if err != nil {
			log.Fatalf("Failed to copy data to env file: %v", err)
		}

		fmt.Println("\n⭐️ denv is right configured and ready to be used!")

		return
	}

	if *up != "" && *name == "" {
		fmt.Println("🌝 Please, provide a nickname to your file using --name flag")
		return
	}

	if *up != "" && *name != "" {
		isEnvSet := helpers.CheckEnvs()

		if !isEnvSet {
			fmt.Println("🤔 Hello! Type 'denv --setup' to start setting up the application")
			return
		}

		currentPath, err := os.Getwd()

		if err != nil {
			log.Fatalf("Failed to get the current path %v", err)
		}

		filePath := path.Join(currentPath, *up)

		fmt.Println(filePath)

		s3bucket.UploadFile(filePath, *name)

		return
	}

	if *name != "" {
		isEnvSet := helpers.CheckEnvs()

		if !isEnvSet {
			fmt.Println("🤔 Hello! Type 'denv --setup' to start setting up the application")
			return
		}

		s3bucket.DownloadFile(*name)
	}
}
