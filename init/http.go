package init

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ManInM00N/go-tool/statics"
	"github.com/tidwall/gjson"
	"github.com/yuin/goldmark/util"
	"io"
	"io/ioutil"
	. "main/DAO"
	"net/http"
	url2 "net/url"
	"os"
	"strconv"
	"time"
)

// TODO: 作者全部作品下载OK
// TODO: 基础下载 OK   目录管理下载 OK  主要图片全部下载OK    并发下载OK
// TODO: 指针内存问题OK
// TODO: 图片下载完整  OK
func Download(i *Illust) {
	var err error
	total := 0
	Request, err2 := http.NewRequest("GET", i.PreviewImageUrl, nil)
	clientcopy := GetClient()
	if err2 != nil {
		DebugLog.Println("Error creating request", err2)
		return
	}
	Request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36")
	Request.Header.Set("referer", "https://www.pixiv.net")
	Cookie := &http.Cookie{
		Name:  "PHPSESSID",
		Value: Setting.Cookie,
	}
	Request.AddCookie(Cookie)
	Request.Header.Set("PHPSESSID", Setting.Cookie)
	var Response *http.Response
	defer func() {
		if Response != nil {
			Response.Body.Close()
		}
	}()
	_, err = os.Stat(Setting.Downloadposition)
	if err != nil {
		os.Mkdir(Setting.Downloadposition, os.ModePerm)
	}
	AuthorFile := Setting.Downloadposition + "/" + statics.Int64ToString(i.UserID)
	_, err = os.Stat(AuthorFile)
	if err != nil {
		os.Mkdir(AuthorFile, os.ModePerm)
	}
	Type := AuthorFile + "/" + i.AgeLimit
	_, err = os.Stat(Type)
	if err != nil {
		os.Mkdir(Type, os.ModePerm)
	}
	failtimes := 0
	for j := int64(0); j < i.Pages; j++ {
		imagefilename := statics.GetFileName(i.ImageUrl[j])
		imagefilepath := Type + "/" + imagefilename
		img, err2 := os.Stat(imagefilepath)
		if err2 == nil {
			if img.Size() != 0 {
				time.Sleep(time.Millisecond * time.Duration(Setting.Downloadinterval))
				continue
			}
		}
		Request.URL, _ = url2.Parse(i.ImageUrl[j])
		ok := true
		for k := 0; k < 10; k++ {
			Response, err = clientcopy.Do(Request)
			if k == 9 && err != nil {
				DebugLog.Println("Illust Resouce Request Error", err, Response.Status)
				ok = false
				j--
				failtimes++
				if failtimes > 2 {
					j++
				}
				break
			} else if err == nil {
				break
			}
			time.Sleep(time.Millisecond * time.Duration(Setting.Downloadinterval))
		}
		if !ok {
			os.Remove(imagefilepath)
			continue
		}
		failtimes = 0
		f, err := os.OpenFile(imagefilepath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
		if err != nil {
			DebugLog.Println(i.Pid, "Download Failed", err, "retrying")
			os.Remove(imagefilepath)
			j--
			continue
		}
		bufWriter := bufio.NewWriter(f)
		_, err = io.Copy(bufWriter, Response.Body)
		if err != nil {
			DebugLog.Println(i.Pid, " Write Failed", err)
		}
		f.Close()
		bufWriter.Flush()
		total++
		time.Sleep(time.Millisecond * time.Duration(Setting.Downloadinterval))
	}
	return
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
		DebugLog.Println("Error creating request", err)
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
			DebugLog.Println("Request failed ", err)
			return nil, err
		}
		time.Sleep(time.Duration(Setting.Retryinterval) * time.Millisecond)

	}
	defer response.Body.Close()

	webpageBytes, err3 := ioutil.ReadAll(response.Body)
	if err3 != nil {
		DebugLog.Println("read failed", err3)
		return nil, err3
	}
	if response.StatusCode != http.StatusOK {
		DebugLog.Println("status code ", response.StatusCode)
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
		DebugLog.Println("GetWebpageData error", err2)
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
		DebugLog.Println("get illustpage data error", err2)
		return nil, err
	}
	imagejson := gjson.ParseBytes(pages).Get("body").String()
	var imagedata []ImageData
	err2 = json.Unmarshal(util.StringToReadOnlyBytes(imagejson), &imagedata)
	if err2 != nil {
		err = fmt.Errorf("error decoding %w", err2)
		DebugLog.Println("Error decoding", err2)
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
func JustDownload(pid string, mode bool) (int, bool) {
	illust, err := work(statics.StringToInt64(pid), 1)
	if !mode {
		if !errors.Is(err, &NotGood{}) && !errors.Is(err, &AgeLimit{}) {
			return 0, true
		}
	}
	if illust == nil {
		DebugLog.Println(pid, " Download failed")
		return 0, false
	}
	if mode {
		InfoLog.Println(pid + " Start download")
	}
	Download(illust)
	if mode {
		InfoLog.Println(pid + " Finished download")
	}
	return 1, true
}
