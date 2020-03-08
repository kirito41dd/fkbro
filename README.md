# fkbro
[![](https://img.shields.io/github/last-commit/zshorz/fkbro)](https://github.com/zshorz/fkbro/)

Telegram Bot [@btcinfo_kirito_bot](https://t.me/btcinfo_kirito_bot) 暴躁老哥 BTC信息查询机器人
## 预览
* 在`Telegram`中搜索 `@btcinfo_kirito_bot` 即可添加机器人 
* `/help`获取操作方法
* 可以把机器人添加到群组
## 编译
* 进入到项目的根目录
* 执行 `go build -o fkbro.exe github.com/zshorz/fkbro`
## 运行
* 在`telegram`中搜索`@BotFather`按步骤申请机器人,得到`token`
* 填写配置文件`config.json`
  * `bot_token`  - 机器人`token`
  * `proxy` - 国内无法连接`telegram`服务器,可配置代理
    * `socks5://127.0.0.1:1080`
    * `http://127.0.0.1:1081`
  * `static_path` - `/static`文件夹的路径,要确保程序运行时可以找到
  * 其余保持默认
* 执行 `./fkbro.exe -c config.json` ,也可以不用`-c`参数,默认执行目录下`config.json`
* [release](https://github.com/zshorz/fkbro/releases) 里已经为Linux提供了编译好的二进制文件



