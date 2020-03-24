# fkbro
[![](https://img.shields.io/github/last-commit/zshorz/fkbro)](https://github.com/zshorz/fkbro/)

Telegram Bot [@btcinfo_kirito_bot](https://t.me/btcinfo_kirito_bot) 暴躁老哥 BTC信息、行情查询机器人
## 预览
* 在`Telegram`中搜索 `@btcinfo_kirito_bot` 即可添加机器人 
* `/help`获取操作方法
* 可以把机器人添加到群组
* 链上资金监控请关注频道 `@fkbrolive`

图片可能和最新版本不符合：
<table>
  <tr>
        <td align="center"><img src="https://raw.githubusercontent.com/zshorz/markdownphoto/master/fkbor/quotes.png" width=400 /><br />
        <td align="center"><img src="https://raw.githubusercontent.com/zshorz/markdownphoto/master/fkbor/newest.png" width=400 /><br />
  </tr>
  <tr>
          <td align="center"><img src="https://raw.githubusercontent.com/zshorz/markdownphoto/master/fkbor/market.png" width=400 /><br />
          <td align="center"><img src="https://raw.githubusercontent.com/zshorz/markdownphoto/master/fkbor/q.png" width=400 /><br />
  </tr>
  <tr>
          <td align="center"><img src="https://raw.githubusercontent.com/zshorz/markdownphoto/master/fkbor/live.jpg" width=400 /><br />
          <td align="center"><img src="https://raw.githubusercontent.com/zshorz/markdownphoto/master/fkbor/live2.jpg" width=400 /><br />
  </tr>
</table>


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
  * `bot_owner` - 机器人拥有者,填telegram用户名
  * `whale_apikey` - rss功能需要的apikey,[WhaleAlert](https://docs.whale-alert.io/#introduction),留空为不启用
    * 启用后，要用`setup.sql`初始化mysql数据库
    * 并填写相关设置
    * `rss` - 服务器启动时，这些频道自动订阅
  * 其余保持默认
* 执行 `./fkbro.exe -c config.json` ,也可以不用`-c`参数,默认执行目录下`config.json`
* [release](https://github.com/zshorz/fkbro/releases) 里已经为Linux提供了编译好的二进制文件

## 关于
* 版本号格式 `x.y.z` - `y`为偶数的是稳定版本
* donate - `17MaxVm9Zm8WQfHrrgA2LVzuEDdSL6AZVN`



