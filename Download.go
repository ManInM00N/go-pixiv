package main

import (
	"bufio"
	"github.com/ManInM00N/go-tool/statics"
	"io"
	"log"
	"net/http"
	url2 "net/url"
	"os"
	"strings"
	"time"
)

// TODO: 作者全部作品下载
// TODO: 基础下载 OK   目录管理下载 OK  主要图片全部下载    并发下载
// TODO: 指针内存问题OK
func (i *Illust) Download() {
	//var Request = new(http.Request)
	var err error
	Request, err2 := http.NewRequest("GET", i.PreviewImageUrl, nil)
	println(Request)
	if err2 != nil {
		log.Println("Error creating request", err2)
		return
	}
	Request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36")
	Request.Header.Set("referer", "https://www.pixiv.net/")
	Request.Header.Set("cookie", settings.Cookie)
	//UserID := Int64ToString(i.UserID)
	var Response *http.Response
	defer func() {
		if Response != nil {
			Response.Body.Close()
		}
	}()
	for j := 0; j < 7; j++ {
		Response, err = client.Do(Request)
		if j == 6 && err != nil {
			log.Println("Read Illust Error", err, i.PreviewImageUrl)
			break
		} else if err == nil {
			break
		}
	}

	_, err = os.Stat(settings.Downloadposition)
	if err != nil {
		os.Mkdir(settings.Downloadposition, os.ModePerm)
	}
	AuthorFile := settings.Downloadposition + "/" + statics.Int64ToString(i.UserID)
	_, err = os.Stat(AuthorFile)
	if err != nil {
		os.Mkdir(AuthorFile, os.ModePerm)
	}
	Type := AuthorFile + "/" + i.AgeLimit
	_, err = os.Stat(Type)
	if err != nil {
		os.Mkdir(Type, os.ModePerm)
	}

	//预览图：
	filename := GetFileName(i.PreviewImageUrl)
	filename = Type + "/" + filename
	var f *os.File
	f, err = os.Create(filename)
	if err != nil {
		log.Println("File Create Error", err)
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	buf := make([]byte, 1024)
	for {
		len, err := Response.Body.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Println("Read error", err)
				os.Remove(filename)
				break
			}
			w.Write(buf[:len])
			break
		}
		w.Write(buf[:len])
	}
	for j := int64(0); j < i.Pages; j++ {
		imagefilename := GetFileName(i.ImageUrl[j])
		imagefilepath := Type + "/" + imagefilename

		img, err2 := os.Create(imagefilepath)
		if err2 != nil {
			log.Println("File Create Error", err2)
			//os.Remove(imagefilepath)
			continue
		}
		w = bufio.NewWriter(img)
		Request.URL, _ = url2.Parse(i.ImageUrl[j])
		ok := true
		for k := 0; k < 5; k++ {
			Response, err = client.Do(Request)
			if k == 4 && err != nil {
				log.Println("Error", err, i.ImageUrl[j])
				ok = false
				break
			} else if err == nil {
				break
			}
		}
		if !ok {
			os.Remove(imagefilepath)
			log.Println("Download Failed", imagefilepath)
			continue
		}
		for {
			var len int
			len, err = Response.Body.Read(buf)
			if err != nil {
				if err != io.EOF {
					log.Println("Read error", err)
					os.Remove(filename)
					break
				}
				w.Write(buf[:len])
				break
			}
			w.Write(buf[:len])
		}
		img.Close()
		log.Println("Download Success", imagefilename)
	}
	time.Sleep(50 * time.Millisecond)
}
func GetFileName(path string) string {
	index := strings.LastIndex(path, "/")
	path = path[index+1:]
	return path
}
