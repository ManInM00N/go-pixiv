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
)

// TODO: 基础下载 OK   目录管理下载 OK  主要图片全部下载    并发下载
func (i *Illust) Download() error {
	Request, _ := http.NewRequest("GET", i.PreviewImageUrl, nil)
	Request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36")
	Request.Header.Set("referer", "https://www.pixiv.net/")
	Request.Header.Set("cookie", settings.Cookie)
	Request.Header.Set("account", settings.Account)
	Request.Header.Set("password", settings.Password)
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
	AuthorFile := "Download/" + statics.Int64ToString(i.UserID)
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
	f, err := os.Create(filename)
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
				return nil
			}
			break
		}
		w.Write(buf[:len])
	}
	for j := int64(0); j < i.Pages; j++ {
		imgfilename := Type + "/" + GetFileName(i.ImageUrl[j])
		img, err := os.Create(imgfilename)
		if err != nil {
			log.Println("File Create Error", err)
		}
		w = bufio.NewWriter(img)
		Request.URL, _ = url2.Parse(i.ImageUrl[j])
		for k := 0; k < 3; k++ {
			Response, err = client.Do(Request)
			if k == 2 && err != nil {
				log.Println("Error", err, i.ImageUrl[j])
				return err
			} else if err == nil {
				break
			}
		}
		for {
			len, err := Response.Body.Read(buf)
			if err != nil {
				if err != io.EOF {
					log.Println("Read error", err)
					os.Remove(filename)
					return nil
				}
				break
			}
			w.Write(buf[:len])
		}
		img.Close()

	}
	return nil
}
func GetFileName(path string) string {
	index := strings.LastIndex(path, "/")
	path = path[index+1:]
	return path
}
