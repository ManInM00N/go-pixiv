package main

import (
	"bufio"
	"errors"
	"github.com/ManInM00N/go-tool/statics"
	"io"
	"log"
	. "main/init"
	"net/http"
	url2 "net/url"
	"os"
	"strings"
	"time"
)

// TODO: 作者全部作品下载OK
// TODO: 基础下载 OK   目录管理下载 OK  主要图片全部下载OK    并发下载OK
// TODO: 指针内存问题OK
// TODO: 图片下载完整  OK
func (i *Illust) Download() {
	var err error
	total := 0
	Request, err2 := http.NewRequest("GET", i.PreviewImageUrl, nil)
	clientcopy := GetClient()
	if err2 != nil {
		log.Println("Error creating request", err2)
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
	//UserID := Int64ToString(i.UserID)
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
		imagefilename := GetFileName(i.ImageUrl[j])
		imagefilepath := Type + "/" + imagefilename
		img, err2 := os.Stat(imagefilepath)
		if err2 == nil {
			if img.Size() != 0 {
				time.Sleep(time.Millisecond * time.Duration(Setting.Downloadinterval))
				continue
			}
		}
		//w := bufio.NewWriter(assets)
		Request.URL, _ = url2.Parse(i.ImageUrl[j])
		ok := true
		for k := 0; k < 10; k++ {
			Response, err = clientcopy.Do(Request)
			if k == 9 && err != nil {
				log.Println("Illust Resouce Request Error", err, Response.Status)
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
			log.Println(i.Pid, "Download Failed", err, "retrying")
			os.Remove(imagefilepath)
			j--
			continue
		}
		defer f.Close()

		bufWriter := bufio.NewWriter(f)
		_, err = io.Copy(bufWriter, Response.Body)
		if err != nil {
			log.Println(i.Pid, " Write Failed", err)
		}
		bufWriter.Flush()
		total++
		time.Sleep(time.Millisecond * time.Duration(Setting.Downloadinterval))
	}
	//log.Println(i.Pid, "Total pictures:", len(i.ImageUrl), "Actually download", total)
	return
}
func JustDownload(pid string, mode bool) (int, bool) {
	illust, err := work(statics.StringToInt64(pid), 1)
	if !mode {
		if !errors.Is(err, &NotGood{}) && !errors.Is(err, &AgeLimit{}) {
			return 0, true
		}
	}
	if illust == nil {
		log.Println(pid, " Download failed")
		return 0, false
	}
	if mode {
		log.Println(pid + " Start download")
	}
	illust.Download()
	if mode {
		log.Println(pid + " Finished download")
	}
	return 1, true
}
func GetFileName(path string) string {
	index := strings.LastIndex(path, "/")
	path = path[index+1:]
	return path
}
