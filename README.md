# 🥸 Denv, the secret and config files' helper tool
I use Ubuntu to work, and my Mac to create personal projects, but sometimes I have to share lots of env files, tokens, and config files... across my machines, which is super boring! (I used to use google drive, and it sucks)
So I have created `denv`, a CLI that allows me to easily upload and download those config files.

## 🤩 How to install denv Mac and Linux
```bash
    make
```

## 🤩 How to install denv Windows
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
    denv --up [filename] --name [nickname]
```

```bash
    denv --name [nickname]

    # or

    denv --name [nickname] --out [filename]

    # ex: denv --name mygitconfig --out .config
```

That is it! 👋🏻
