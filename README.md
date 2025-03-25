# 🥸 Denv, the secret and config files' helper tool

I use Ubuntu for work and my Mac for personal projects. However, I often find myself needing to share numerous environment files, tokens, and config files between my machines, which can be quite tedious. In the past, I used Google Drive, but it was not an ideal solution.

To address this problem, I created `denv`, a command-line interface (CLI) tool that simplifies the process of uploading and downloading these config files.

## 🤩 How to install denv - Mac and Linux
```bash
make
```

## 🤩 How to install denv - Windows
```bash
mkdir C:\bin
```
```bash
go build -o denv.exe cmd\denv\main.go
```
```bash
mv denv.exe C:\bin\denv.exe
```
```bash
setx PATH "C:\bin;%PATH%"
```
```bash
# ensure it is right placed
where.exe denv.exe
```

## 😜 How to configure
```bash
# You will need to have your AWS secret key, access key, and a S3 bucket name ready
denv --config
```

## 🤯 Commands

### Upload files
```bash
# To upload a file with a specific nickname
denv --up [filename] --name [nickname]

# Example: Upload .env file with the nickname "myproject"
denv --up .env --name myproject
```

### Download files
```bash
# To download a file using its nickname
denv --name [nickname]

# To download a file with a custom output filename
denv --name [nickname] --out [output-filename]

# Example: Download a file nicknamed "myproject" and save it as .env.production
denv --name myproject --out .env.production
```

### List files
```bash
# To list all files stored in your bucket
denv --list
```

### Delete files
```bash
# To delete a file from the bucket
denv --del [nickname]

# Example: Delete a file nicknamed "old-config"
denv --del old-config
```

### Rename files
```bash
# To rename a file in the bucket
denv --rename [old-nickname] --name [new-nickname]

# Example: Rename "config-dev" to "config-development"
denv --rename config-dev --name config-development
```

### Tab Completion

Denv supports tab completion for file names when using commands like `--del`, `--rename`, and `--name`. To set up tab completion:

```bash
# Set up bash completion for denv commands
denv --setup-completion

# After installation, restart your shell or run:
source ~/.bash_completion
```

Once set up, you can use Tab to complete file names:
```bash
denv --del [TAB]               # Shows all available files
denv --rename config-[TAB]     # Shows files starting with "config-"
```

### Help
```bash
# To display help information about all commands
denv --help
```

## 📝 Usage Examples

### Complete workflow example:
```bash
# Upload your environment file
denv --up .env.development --name dev-env

# List all your stored files
denv --list

# Download the file on another machine
denv --name dev-env --out .env.development

# Rename the file to something more descriptive
denv --rename dev-env --name project-development-env

# When you don't need it anymore
denv --del project-development-env
```

That is it! 👋🏻