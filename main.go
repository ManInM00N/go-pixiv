package main

import (
	"fmt"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"unsafe"
)

const (
	url     string = "https://p1.ssl.qhimg.com/t01d5be5429abbf0de5.png"
	texturl string = "https://p1.ssl.qhimg.com/"
)

type Illust struct {
	Pid         int64    `db:"pid"`
	Title       string   `db:"title"`
	Caption     string   `db:"caption"`
	Tags        string   `db:"tags"`
	ImageUrls   []string `db:"image_urls"`
	AgeLimit    string   `db:"age_limit"`
	CreatedTime string   `db:"created_time"`
	UserID      int64    `db:"user_id"`
	UserName    string   `db:"user_name"`
}

var rg *gin.Engine

func main() {
	rg = gin.Default()
	app := app.New()
	appwindow := app.NewWindow("GO Pixiv")
	app.Run()

	text := widget.NewTextGrid()
	content := widget.NewButton("click me", func() {
		work(StringToInt64(text.Text()))
	})
	appwindow.SetContent(content)
	appwindow.SetContent(text)
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
	fmt.Println(ageLimit)
	fmt.Println(json)
	return i, nil
}

func GetWebpageData(url string) (data []byte, err error) {
	response, err := http.Get(url)
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
