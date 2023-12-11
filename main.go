package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"log"
	"net/http"
	url2 "net/url"
	"os"
	"strconv"
	"unsafe"
)

// TODO: 可视化gui 查找作者、插图  gui 20%
// TODO： client代理 100% json网络请求 100%  body下载 30%
const (
	url     string = "https://p1.ssl.qhimg.com/t01d5be5429abbf0de5.png"
	texturl string = "https://p1.ssl.qhimg.com/"
)

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
}

func (i *Illust) msg() string {
	var tags string
	for key, value := range i.Tags {
		tags = tags + value
	}
	return strconv.FormatInt(i.Pid, 10) + "\n  " + i.Title + "\n  " + i.Caption + "\n  " + i.Tags + "\n " + i.AgeLimit + "\n  " + strconv.FormatInt(i.UserID, 10) + "\n " + i.UserName

}

var rg *gin.Engine
var f *os.File
var client *http.Client

func LogInit() {
	log.SetFlags(log.Ldate | log.Ltime)
	f, _ = os.OpenFile("temp.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	log.SetOutput(f)
}

func main() {
	rg = gin.Default()
	app := app.New()
	appwindow := app.NewWindow("GO Pixiv")
	//app.Run()
	LogInit()

	proxyURL, err := url2.Parse("http://127.0.0.1:10809")
	if err != nil {
		log.Println(err)
		return
	}
	client = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		},
	}

	text := widget.NewEntry()
	button := widget.NewButton("click me", func() {
		work(StringToInt64(text.Text))
	})
	content := container.New(layout.NewVBoxLayout(), text, button)
	appwindow.SetContent(content)
	appwindow.Resize(fyne.Size{600, 400})
	appwindow.ShowAndRun()
}
func work(id int64) (i *Illust, err error) {
	data, err := GetWebpageData("https://www.pixiv.net/ajax/illust/" + strconv.FormatInt(id, 10))
	if err != nil {
		log.Fatalln(err)
	}
	json := gjson.ParseBytes(data).Get("body")
	var ageLimit = "all-age"
	for _, tag := range json.Get("tags.tags.#.tag").Array() {
		if tag.Str == "R-18" {
			ageLimit = "r18"
			break
		}
	}
	i = &Illust{}
	i.AgeLimit = ageLimit
	i.Pid = json.Get("illustId").Int()
	i.UserID = json.Get("userID").Int()

	for key, _ := range json.Get("tags").Map() {
		i.Tags = append(i.Tags, key)
	}
	i.Caption = json.Get("alt").Str
	i.CreatedTime = json.Get("createDate").Str
	for key, _ := range json.Get("urls").Map() {
		i.ImageUrls = append(i.ImageUrls, key)
	}
	i.Title = json.Get("illustTitle").Str
	i.UserName = json.Get("userName").Str
	log.Println(i.msg())
	//fmt.Println(json)
	return i, nil
}

func GetWebpageData(url string) (data []byte, err error) {
	response, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	webpageBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return webpageBytes, nil
}
func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
func StringToInt64(s string) int64 {
	var num int64 = 0
	for i := 0; i < len(s); i++ {
		num = num*10 + int64(s[i]-'0')
	}
	return num
}
