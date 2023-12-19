package main

import (
	"encoding/json"
	"github.com/tidwall/gjson"
	"github.com/yuin/goldmark/util"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

type NotGood struct{}
type AgeLimit struct{}

func (i *NotGood) Error() string {
	return "LikeNotEnough"
}
func (i *AgeLimit) Error() string {
	return "AgeLimitExceed"
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

// TODO:下载作品主题信息json OK
func GetWebpageData(url, id string) ([]byte, error) { //请求得到作品json

	var response *http.Response
	var err error
	Request, err := http.NewRequest("GET", "https://www.pixiv.net/ajax/illust/"+url, nil)
	Request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36")
	Request.Header.Set("referer", "https://www.pixiv.net/artworks/"+id)
	Request.Header.Set("cookie", "PHPSESSID="+settings.Cookie)
	clientcopy := client
	for i := 0; i < 10; i++ {
		response, err = clientcopy.Do(Request)
		if err == nil {
			break
		}
		if i == 9 && err != nil {
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

func GetAuthorWebpage(url, id string) ([]byte, error) {

	var response *http.Response
	var err error
	Request, err := http.NewRequest("GET", "https://www.pixiv.net/ajax/user/"+url+"/profile/all", nil)
	Request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36")
	Request.Header.Set("referer", "https://www.pixiv.net/member.php?id="+id)
	Request.Header.Set("cookie", "PHPSESSID="+settings.Cookie)
	clientcopy := client

	for i := 0; i < 10; i++ {
		response, err = clientcopy.Do(Request)
		if err == nil {
			break
		}
		if i == 9 && err != nil {
			log.Println("Request failed ", err)
			return nil, err
		}
	}
	defer response.Body.Close()
	webpageBytes, err3 := ioutil.ReadAll(response.Body)
	if err3 != nil {
		log.Println("read failed", err3)
		os.Exit(4)
		return nil, err3
	}
	if response.StatusCode != http.StatusOK {
		log.Println("status code ", response.Status)
	}
	return webpageBytes, nil
}

// TODO: 作品信息json请求   OK
// TODO: 多页下载
func work(id int64) (i *Illust, err error) { //按作品id查找
	urltail := strconv.FormatInt(id, 10)
	strid := urltail
	data, err := GetWebpageData(urltail, strid)
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

	pages, err := GetWebpageData(urltail+"/pages", strid)
	if err != nil {
		return nil, err
	}
	imagejson := gjson.ParseBytes(pages).Get("body").String()
	var imagedata []ImageData
	err = json.Unmarshal(util.StringToReadOnlyBytes(imagejson), &imagedata)
	if err != nil {
		log.Println("Error decoding", err)
	}

	i.PreviewImageUrl = imagedata[0].URLs.ThumbMini
	for _, image := range imagedata {
		i.ImageUrl = append(i.ImageUrl, image.URLs.Original)
	}
	if i.Likecount < settings.LikeLimit {
		return i, &NotGood{}
	}
	if i.AgeLimit == "r18" && !settings.Agelimit {
		return i, &AgeLimit{}
	}

	return i, nil
}
func GetAuthor(id int64, ss *map[string]gjson.Result) error {
	data, err := GetAuthorWebpage(strconv.FormatInt(id, 10), strconv.FormatInt(id, 10))
	if err != nil {
		return err
	}
	jsonmsg := gjson.ParseBytes(data).Get("body")
	*ss = jsonmsg.Get("illusts").Map()
	return nil
}
