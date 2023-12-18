<h1>使用须知</h1>
<h1>
！！！不提供代理！！！
<br>
梯子代理问题需要自己解决
<br>
需要有自己的pixiv账号,若个人行为导致账号封禁概不负责
<br>
压力测试200张/s可行，默认设置间隔时间50ms
</h1>
<h2>
目前没有设计gui，暂时使用fyne做简陋的读入，作者现在大二前端写的依托，所以等摸完了vue，再换成wails在做成包发布<br>
且存在大型图片容易读取不完整，小图片也存在这种问题
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

</p>