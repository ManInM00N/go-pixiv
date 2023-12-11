package main

import (
	"github.com/tidwall/gjson"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"unsafe"
)

// TODO: 可视化gui 查找作者、插图   gui 5%
// TODO： client代理 100% json网络请求 100%  header下载请求 0%

type Illust struct {
	Pid         int64    `db:"pid"`
	Title       string   `db:"title"`
	Caption     string   `db:"caption"`
	Tags        []string `db:"tags"`
	ImageUrls   []string `db:"image_urls"`
	AgeLimit    string   `db:"age_limit"`
	CreatedTime string   `db:"created_time"`
	UserID      int64    `db:"user_id"`
	UserName    string   `db:"user_name"`
	Pages       int64    `db:"pages"`
}

func (i *Illust) msg() string {
	var tags string
	for _, value := range i.Tags {
		tags = tags + value + "\n"
	}
	return strconv.FormatInt(i.Pid, 10) + "\n  " + i.Title + "\n  " + i.Caption + "\n  " + tags + "\n " + i.AgeLimit + "\n  " + strconv.FormatInt(i.UserID, 10) + "\n " + i.UserName

}
func (i *Illust) Download() error {

	UserID := Int64ToString(i.UserID)

	return nil
}
func main() {
	LogInit()     //日志打印
	windowInit()  //gui面板
	clinentInit() //服务端请求设置
	appwindow.ShowAndRun()
}
func work(id int64) (i *Illust, err error) { //按作品id查找
	data, err := GetWebpageData("https://www.pixiv.net/ajax/illust/" + strconv.FormatInt(id, 10))
	if err != nil {
		log.Println("Request failed", err)
		//log.Fatalln(err)
		os.Exit(2)
	}
	json := gjson.ParseBytes(data).Get("body") //读取json内作品及作者id信息
	var ageLimit = "all-age"
	i = &Illust{}
	for _, tag := range json.Get("tags.tags.#.tag").Array() {
		i.Tags = append(i.Tags, tag.Str)
		if tag.Str == "R-18" {
			ageLimit = "r18"
			break
		}
	}
	i.AgeLimit = ageLimit
	i.Pid = json.Get("illustId").Int()
	i.UserID = json.Get("userId").Int()
	i.Caption = json.Get("alt").Str
	i.CreatedTime = json.Get("createDate").Str
	i.Pages = json.Get("pageCount").Int()
	i.Title = json.Get("illustTitle").Str
	i.UserName = json.Get("userName").Str
	for key, _ := range json.Get("urls").Map() {
		i.ImageUrls = append(i.ImageUrls, key)
	}
	log.Print(i.msg())
	return i, nil
}

func GetWebpageData(url string) ([]byte, error) { //请求得到作品json
	response, err := client.Get(url)
	if err != nil {
		log.Println("Request failed ", err)
		log.Fatalln(err)

		os.Exit(3)
		return nil, err
	}
	defer response.Body.Close()

	webpageBytes, err3 := ioutil.ReadAll(response.Body)
	if err3 != nil {
		log.Println("read failed", err)
		os.Exit(4)
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		log.Println("status code ", response.StatusCode)
	}
	return webpageBytes, nil
}
func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
