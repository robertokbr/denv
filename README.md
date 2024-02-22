# 🥸 Denv, the configuration files helper tool
I use to use Ubuntu to work, and my Mac to create personal projects, but sometimes I have to share lots of env files, tokens, config files... across my machines, which is super boring! (I use to use google drive, and it really sucks!)
So I have created Denv, a CLI that allow me to upload and download those config files with easy.

## 🤩 How to install dev
```bash
    make
```

## 😜 How to config
```bash
    # you will need to get in hands your AWS secret key, access key and a S3 bucket name
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