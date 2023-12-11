package main

import (
	"bufio"
	"github.com/tidwall/gjson"
	"io"
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
func (i *Illust) Download() error {
	Request, _ := http.NewRequest("GET", i.PreviewImageUrl, nil)
	Request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36")
	Request.Header.Set("referer", "https://www.pixiv.net/")
	Request.Header.Set("cookie", settings.Cookie)
	//UserID := Int64ToString(i.UserID)
	var Response *http.Response
	var err error

	for j := 0; j < 3; j++ {
		Response, err = client.Do(Request)
		if j == 2 && err != nil {
			log.Println("Error", err, i.PreviewImageUrl)
			return err
		} else if err == nil {
			break
		}
	}
	log.Println(i.PreviewImageUrl, Response.StatusCode, Response.Body)
	f, _ := os.Create(i.PreviewImageUrl)
	w := bufio.NewWriter(f)

	buf := make([]byte, 1024)
	for {
		len, err := Response.Body.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Println("Read error", err)
				os.Remove(i.PreviewImageUrl)
				return nil
			}
			break
		}
		w.Write(buf[:len])
	}
	return nil
}

func main() {
	LogInit()     //日志打印
	windowInit()  //gui面板
	clinentInit() //服务端请求设置
	appwindow.ShowAndRun()
}

// TODO: 作品信息json请求   OK
func work(id int64) (i *Illust, err error) { //按作品id查找
	data, err := GetWebpageData("https://www.pixiv.net/ajax/illust/" + strconv.FormatInt(id, 10))

	if err != nil {
		return nil, err
	}
	pages, err := GetWebpageData("https://www.pixiv.net/ajax/illust/" + strconv.FormatInt(id, 10) + "/pages")
	if err != nil {
		return nil, err
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
	i.PreviewImageUrl = json.Get("urls.thumb").Str
	json = gjson.ParseBytes(pages).Get("body")
	for _, url := range json.Get(".original").Array() {
		i.ImageUrl = append(i.ImageUrl, url.Str)
		log.Println(url.Str)
	}
	//log.Print(i.msg())
	return i, nil
}

// TODO ：下载作品json OK
func GetWebpageData(url string) ([]byte, error) { //请求得到作品json
	//response, err := client.Get(url)
	var response *http.Response
	var err error
	for i := 0; i < 3; i++ {
		response, err = client.Get(url)
		if err == nil {
			break
		}
		if i == 2 && err != nil {
			log.Println("Request failed ", err)
			return nil, err
		}
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
