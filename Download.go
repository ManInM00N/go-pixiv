package main

import (
	"github.com/ManInM00N/go-tool/statics"
	"io/ioutil"
	"log"
	"net/http"
	url2 "net/url"
	"os"
	"strings"
	"time"
)

// TODO: 作者全部作品下载OK
// TODO: 基础下载 OK   目录管理下载 OK  主要图片全部下载OK    并发下载OK
// TODO: 指针内存问题OK
// TODO: 图片下载完整  ????
func (i *Illust) Download() {
	//var Request = new(http.Request)
	var err error
	total := 0
	Request, err2 := http.NewRequest("GET", i.PreviewImageUrl, nil)
	clientcopy := client
	if err2 != nil {
		log.Println("Error creating request", err2)
		return
	}
	Request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36")
	Request.Header.Set("referer", "https://www.pixiv.net")
	Request.Header.Set("cookie", settings.Cookie)
	//UserID := Int64ToString(i.UserID)
	var Response *http.Response
	defer func() {
		if Response != nil {
			Response.Body.Close()
		}
	}()
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
	failtimes := 0
	for j := int64(0); j < i.Pages; j++ {
		imagefilename := GetFileName(i.ImageUrl[j])
		imagefilepath := Type + "/" + imagefilename

		//img, err2 := os.Create(imagefilepath)
		//if err2 != nil {
		//	log.Println("File Create Error", err2)
		//	os.Remove(imagefilepath)
		//	continue
		//}
		//w := bufio.NewWriter(img)
		Request.URL, _ = url2.Parse(i.ImageUrl[j])
		ok := true
		for k := 0; k < 10; k++ {
			Response, err = clientcopy.Do(Request)
			if k == 9 && err != nil {
				log.Println("Illust Resouce Request Error", err, Response.Status)
				//log.Println("Retry Download", i.ImageUrl[j])
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
			time.Sleep(time.Millisecond * 100)

		}
		if !ok {
			os.Remove(imagefilepath)
			continue
		}
		failtimes = 0
		buf, err := ioutil.ReadAll(Response.Body)
		if err != nil {
			log.Println(i.Pid, "Download Failed", err, "retrying")
			os.Remove(imagefilepath)

			j--
			continue
		}
		err = ioutil.WriteFile(imagefilepath, buf, 666)
		total++
		if err != nil {
			log.Println("Write Failed", err)
			//buf := make([]byte, 1024)
			//ok = true
			//for {
			//	var len int
			//	len, err = Response.Body.Read(buf)
			//	if err != nil {
			//		if err != io.EOF {
			//			log.Println("Read error", err)
			//			os.Remove(imagefilepath)
			//			break
			//		}
			//		lr += int64(len)
			//		_, err := w.Write(buf[:len])
			//		if err != nil {
			//			log.Println("Read bytes error", err)
			//			ok = false
			//			break
			//		}
			//		break
			//	}
			//	lr += int64(len)
			//	_, err := w.Write(buf[:len])
			//	if err != nil {
			//		ok = false
			//		log.Println("Read bytes error", err)
			//		break
			//	}
			//}
			//if !ok {
			//	j--
			//	os.Remove(imagefilepath)
			//	continue
			//}
			//log.Println(lr)
			//img.Close()

		}
	}
	time.Sleep(100 * time.Millisecond)
	//log.Println(i.Pid, "Total pictures:", len(i.ImageUrl), "Actually download", total)
	return
}
func GetFileName(path string) string {
	index := strings.LastIndex(path, "/")
	path = path[index+1:]
	return path
}
