package main

import (
	"strconv"
)

// TODO: 可视化gui 查找作者、插图   gui 5%
// TODO： client代理 100% json网络请求 100%  header下载请求 0%

type Illust struct {
	Pid             int64    `db:"pid"`
	Title           string   `db:"title"`
	Caption         string   `db:"caption"`
	Tags            []string `db:"tags"`
	ImageUrl        []string `db:"image_url"`
	PreviewImageUrl string   `db:"preview_image"`
	AgeLimit        string   `db:"age_limit"`
	CreatedTime     string   `db:"created_time"`
	UserID          int64    `db:"user_id"`
	UserName        string   `db:"user_name"`
	Pages           int64    `db:"pages"`
}

func (i *Illust) msg() string {
	var tags string
	for _, value := range i.Tags {
		tags = tags + value + "\n"
	}
	return strconv.FormatInt(i.Pid, 10) + "\n  " + i.Title + "\n  " + i.Caption + "\n  " + tags + "\n " + i.AgeLimit + "\n  " + strconv.FormatInt(i.UserID, 10) + "\n " + i.UserName

}

func main() {
	LogInit()     //日志打印
	windowInit()  //gui面板
	clinentInit() //服务端请求设置
	appwindow.ShowAndRun()
}
