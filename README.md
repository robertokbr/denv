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

## 😜 How to config
```bash
    # You will need to get in hands your AWS secret key, access key, and a S3 bucket name
    denv --conf
```

## 🤯 How to use
```bash
    # To upload a file
    denv --up [filename] --name [nickname]
```

```bash
    # To download a file
    denv --name [nickname]

    # or

    denv --name [nickname] --out [filename]

    # ex: denv --name mygitconfig --out .config
```

```bash
    # To list all files
    denv --list
```

```bash
    # To delete a file
    denv --del [nickname]
```

That is it! 👋🏻
