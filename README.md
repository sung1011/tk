# tk

tk is a CLI for GO that tickles applications

## INSTALL

`go get -v -u github.com/sung1011/tk`

## Config

`$HOME/.tk.yaml`   (参照 common/.tk.yaml.default)

## Usage

```txt
tk [command]

Available Commands:
  ascii       ASCII码表，可传入字符以搜索
  build       (for dev)修改tk库后直接编译
  info        计算机基础信息
  new         (for dev)新建一个命令
  scanport    扫描端口
  ssh         ssh连接远端机器，可传入想执行的命令
  time        转化时间，可传入字符串或时间戳
```

### TODO

```txt
   httpCode表  
   time细化 如:+7day
   同步文件
   网页收藏夹
   加密解密 AES DES RSA SHA
   crontab验证 输出未来10次的执行时间  
```
