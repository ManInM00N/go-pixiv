package main

import (
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"github.com/yuin/goldmark/util"
	"io/ioutil"
	"log"
	. "main/init"
	"net/http"
	"strconv"
	"time"
)

type NotGood struct {
	S string
}
type AgeLimit struct {
	S string
}

func (i *NotGood) Error() string {
	return i.S
}
func (i *AgeLimit) Error() string {
	return i.S
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

func CheckMode(url, id string, num int) (string, string) {
	if num == 1 { //illust page
		return "https://www.pixiv.net/ajax/illust/" + url, "https://www.pixiv.net/artworks/" + id
	} else if num == 2 { // author page
		return "https://www.pixiv.net/ajax/user/" + url + "/profile/all", "https://www.pixiv.net/member.php?id=" + id
	} else if num == 3 {

	}
	return "", ""
}

// TODO:下载作品主题信息json OK
func GetWebpageData(url, id string, num int) ([]byte, error) { //请求得到作品json

	var response *http.Response
	var err error
	s1, s2 := CheckMode(url, id, num)
	Request, err := http.NewRequest("GET", s1, nil)
	if err != nil {
		log.Println("Error creating request", err)
		return nil, err
	}
	Request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36")
	Request.Header.Set("referer", s2)
	Cookie := &http.Cookie{
		Name:  "PHPSESSID",
		Value: Setting.Cookie,
	}
	Request.AddCookie(Cookie)
	Request.Header.Set("PHPSESSID", Setting.Cookie)

	clientcopy := GetClient()
	for i := 0; i < 10; i++ {
		response, err = clientcopy.Do(Request)
		if err == nil {
			if response.StatusCode == 429 {
				time.Sleep(time.Duration(Setting.Retry429) * time.Millisecond)
				//i--
				continue
			}
			break
		}
		if i == 9 && err != nil {
			log.Println("Request failed ", err)
			return nil, err
		}
		time.Sleep(time.Duration(Setting.Retryinterval) * time.Millisecond)

	}
	defer response.Body.Close()

	webpageBytes, err3 := ioutil.ReadAll(response.Body)
	if err3 != nil {
		log.Println("read failed", err3)
		return nil, err3
	}
	if response.StatusCode != http.StatusOK {
		log.Println("status code ", response.StatusCode)
		if response.StatusCode == 429 {
			time.Sleep(time.Duration(Setting.Retry429) * time.Millisecond)
		}
	}
	return webpageBytes, nil
}

// TODO: 作品信息json请求   OK
// TODO: 多页下载 OK
func work(id int64, mode int) (i *Illust, err error) { //按作品id查找
	urltail := strconv.FormatInt(id, 10)
	strid := urltail
	err = nil
	data, err2 := GetWebpageData(urltail, strid, mode)

	if err2 != nil {
		err = fmt.Errorf("GetWebpageData error %w", err2)
		log.Println("GetWebpageData error", err2)
		return nil, err
	}
	jsonmsg := gjson.ParseBytes(data).Get("body") //读取json内作品及作者id信息
	i = &Illust{
		AgeLimit:    "all-age",
		Pid:         jsonmsg.Get("illustId").Int(),
		UserID:      jsonmsg.Get("userId").Int(),
		Caption:     jsonmsg.Get("alt").Str,
		CreatedTime: jsonmsg.Get("createDate").Str,
		Pages:       jsonmsg.Get("pageCount").Int(),
		Title:       jsonmsg.Get("illustTitle").Str,
		UserName:    jsonmsg.Get("userName").Str,
		Likecount:   jsonmsg.Get("bookmarkCount").Int(),
	}
	for _, tag := range jsonmsg.Get("tags.tags.#.tag").Array() {
		i.Tags = append(i.Tags, tag.Str)
		if tag.Str == "R-18" {
			i.AgeLimit = "r18"
			break
		}
	}
	if i.Likecount < Setting.LikeLimit {
		err = fmt.Errorf("%w", &NotGood{"LikeNotEnough"})
	}
	if i.AgeLimit == "r18" && !Setting.Agelimit {
		err = fmt.Errorf("%w", &AgeLimit{"AgeLimitExceed"})
	}
	pages, err2 := GetWebpageData(urltail+"/pages", strid, mode)
	if err2 != nil {
		err = fmt.Errorf("Get illustpage data error %w", err2)
		log.Println("get illustpage data error", err2)
		return nil, err
	}
	imagejson := gjson.ParseBytes(pages).Get("body").String()
	var imagedata []ImageData
	err2 = json.Unmarshal(util.StringToReadOnlyBytes(imagejson), &imagedata)
	if err2 != nil {
		err = fmt.Errorf("error decoding %w", err2)
		log.Println("Error decoding", err2)
		return nil, err
	}

	i.PreviewImageUrl = imagedata[0].URLs.ThumbMini
	for _, image := range imagedata {
		i.ImageUrl = append(i.ImageUrl, image.URLs.Original)
	}

	return i, err
}
func GetAuthor(id int64) (map[string]gjson.Result, error) {
	data, err := GetWebpageData(strconv.FormatInt(id, 10), strconv.FormatInt(id, 10), 2)
	if err != nil {
		return nil, err
	}
	jsonmsg := gjson.ParseBytes(data).Get("body")
	ss := jsonmsg.Get("illusts").Map()
	return ss, nil
}
