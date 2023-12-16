package main

import (
	"encoding/json"
	"github.com/tidwall/gjson"
	"github.com/yuin/goldmark/util"
	"io/ioutil"
	"log"
	"net/http"
	url2 "net/url"
	"os"
	"strconv"
	"time"
)

type Settings struct {
	Proxy    string `json:"proxy"`
	Account  string `json:"account"`
	Password string `json:"password"`
	Cookie   string `json:"cookie"`
}

var settings Settings

type NotGood struct{}

func (i *NotGood) Error() string {
	return "LikeNotEnough"
}

type ImageData struct {
	URLs struct {
		ThumbMini string `json:"thumb_mini"`
		Small     string `json:"small"`
		Regular   string `json:"regular"`
		Original  string `json:"original"`
	} `json:"urls"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

// TODO ：客户端代理设置 OK
func clinentInit() {
	jfile, _ := os.Open("Config.json")
	defer jfile.Close()
	bytevalue, _ := ioutil.ReadAll(jfile)
	json.Unmarshal(bytevalue, &settings)
	log.Printf(settings.Proxy)
	proxyURL, err := url2.Parse(settings.Proxy)
	if err != nil {
		log.Println(err)
		return
	}
	client = &http.Client{
		Transport: &http.Transport{
			Proxy:                 http.ProxyURL(proxyURL),
			DisableKeepAlives:     true,
			ResponseHeaderTimeout: time.Second * 5,
		},
	}
}

// TODO:下载作品主题信息json OK
func GetWebpageData(url string) ([]byte, error) { //请求得到作品json
	var response *http.Response
	var err error
	for i := 0; i < 5; i++ {
		response, err = client.Get(url)
		if err == nil {
			break
		}
		if i == 4 && err != nil {
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
	jsonmsg := gjson.ParseBytes(data).Get("body") //读取json内作品及作者id信息
	var ageLimit = "all-age"
	i = &Illust{}
	for _, tag := range jsonmsg.Get("tags.tags.#.tag").Array() {
		i.Tags = append(i.Tags, tag.Str)
		if tag.Str == "R-18" {
			ageLimit = "r18"
			break
		}
	}
	i.AgeLimit = ageLimit
	i.Pid = jsonmsg.Get("illustId").Int()
	i.UserID = jsonmsg.Get("userId").Int()
	i.Caption = jsonmsg.Get("alt").Str
	i.CreatedTime = jsonmsg.Get("createDate").Str
	i.Pages = jsonmsg.Get("pageCount").Int()
	i.Title = jsonmsg.Get("illustTitle").Str
	i.UserName = jsonmsg.Get("userName").Str
	i.Likecount = jsonmsg.Get("likeCount").Int()
	if i.Likecount < int64(1500) {
		return nil, &NotGood{}
	}
	imagejson := gjson.ParseBytes(pages).Get("body").String()
	all := gjson.ParseBytes(pages).Get("").Array()
	for _, image := range all {
		log.Printf(image.String())
	}
	log.Println(imagejson)
	var imagedata []ImageData
	err = json.Unmarshal(util.StringToReadOnlyBytes(imagejson), &imagedata)
	if err != nil {
		log.Println("Error decoding", err)
	}

	i.PreviewImageUrl = imagedata[0].URLs.ThumbMini
	for _, image := range imagedata {
		i.ImageUrl = append(i.ImageUrl, image.URLs.Original)
	}
	log.Print(i.ImageUrl)
	return i, nil
}
func GetAuthor(id int64, ss *map[string]gjson.Result) error {
	data, err := GetWebpageData("https://www.pixiv.net/ajax/user/" + strconv.FormatInt(id, 10) + "/profile/all")
	if err != nil {
		return err
	}
	jsonmsg := gjson.ParseBytes(data).Get("body")
	*ss = jsonmsg.Get("illusts").Map()
	return nil
}
