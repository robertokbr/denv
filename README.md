# 🥸 Denv, the secret files' helper tool
I use Ubuntu to work, and my Mac to create personal projects, but sometimes I have to share lots of env files, tokens, and config files... across my machines, which is super boring! (I used to use google drive, and it really sucks!)
So I have created Denv, a CLI that allows me to easily upload and download those config files.

## 🤩 How to install denv
```bash
    make
```

## 😜 How to config
```bash
    # You will need to get in hands your AWS secret key, access key, and a S3 bucket name
    denv --conf
```

## 🤯 How to use
```bash
    denv --up [filename] --name [nickname to use to download your file]
```

```bash
    denv --name [nickname used to upload your file]
```

That is it! 👋🏻
