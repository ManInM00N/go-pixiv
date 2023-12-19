<h1>使用须知 Read Me First!!!</h1>
<h1>
！！！不提供代理！！！
<br>
梯子代理问题需要自己解决
<br>
需要有自己的pixiv账号,若个人行为导致账号封禁概不负责<br>
保证代理无误的情况下，不添加cookie也可以下载<br>
<br>
下载速度取决于你的代理，理论上参数可以调很快但是会给429，而且太快会被封443端口，所以不建议下太快，默认设置间隔时间80ms
</h1>
<h2>
目前没有设计gui，暂时使用fyne做简陋的读入，作者现在大二前端写的依托，所以等寒假摸完了vue，再换成wails在做成包发布<br>
现存在大型图片容易读取不完整，小图片也存在这种问题，故临时更换读取方法，下载速度变慢<br>
使用参考:https://github.com/daydreamer-json/pixiv-ajax-api-docs/tree/main<br>
</h2>
<h2>配置设定</h2>
<h3>在settings.yml中</h3>
<p>
proxy:你本地梯子的代理ip后面的端口，可以从你梯子的配置中得到，这个不会配的话我无能为力,拿v2ray的配置方法举例，端口就是http后面的数字：<br>
<img src="https://github.com/ManInM00N/go-pixiv/blob/master/img/proxy.png"><br>
cookie:打开登录后的pixiv网页，在电脑网页按F12，从应用程序一栏中Cookie的PHPSESSID的值<br>
<img src="https://github.com/ManInM00N/go-pixiv/blob/master/img/cookie1.png"><br>
<img src="https://github.com/ManInM00N/go-pixiv/blob/master/img/cookie2.png"><br>
<img src="https://github.com/ManInM00N/go-pixiv/blob/master/img/cookie3.png"><br>
r-18:true启用，false禁用   懂得都懂<br>
minlikelimit:下载图片的点赞数限制 小于的不下载<br>
downloadposition:图片储存位置，如果目标位置没有文件夹则会改成此目录下的Download文件夹(自动创建)<br>
queuelimit:下载队列并发数，多出来的阻塞等待<br>
illustqueuelimit:每张插画的并发数，多出来的阻塞等待<br>
</p>