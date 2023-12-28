package main

import (
	. "main/init"
)

// TODO: 可视化gui 查找作者、插图   gui 5%
// TODO： client代理 100% json网络请求 100%  header下载请求 100%

func main() {
	Appwindow.ShowAndRun()
	IsClosed = true
	P.Wait()
	defer func() {
		P.Release()
		TaskPool.Close()
		SinglePool.Release()
	}()
}
