# chain-util
chain util


## 1. 介绍
目前支持的功能有：
- BSC、TRON的链
- 创建钱包地址、导入地址 ：根据助记词生成私钥和地址、根据私钥生成地址

注意要用ssh方式 配置好ssh key

引入私有库：
```
go get github.com/EntySquare/chain-util
```

配置git不以http方式拉取
```
git config --global url."git@github.com:".insteadOf https://github.com/
```
配置GoMod私有仓库
```
go env -w GOPRIVATE=github.com
```

go mod tidy
