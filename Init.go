package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/ManInM00N/go-tool/goruntine"
	"github.com/ManInM00N/go-tool/statics"
	"github.com/devchat-ai/gopool"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"net/http"
	url2 "net/url"
	"os"
	"time"
)

type Settings struct {
	Proxy            string `yml:"proxy"`
	Cookie           string `yml:"cookie"`
	Agelimit         bool   `yml:"r-18" `
	Downloadposition string `yml:"downloadposition"`
	LikeLimit        int64  `yml:"minlikelimit"`
	Queuelimit       int    `yml:"queuelimit"`
	Illustqueuelimit int    `yml:"illustqueuelimit"`
}

var settings Settings

var appwindow fyne.Window
var f *os.File

type FyneLogWriter struct {
	LogText *widget.Entry
}

func (w *FyneLogWriter) Write(p []byte) (n int, err error) {
	message := string(p)
	w.LogText.SetText(w.LogText.Text + message)
	return len(p), nil
}

var TaskPool *goruntine.GoPool

// var P *ants.Pool
var P gopool.GoPool
var satisfy int

func windowInit() {
	app := app.New()
	//defer P.Release()
	//defer taskpool.Release()
	TaskPool = goruntine.NewGoPool(1)
	P = gopool.NewGoPool(4, gopool.WithTaskQueueSize(5000))
	//P, _ = ants.NewPool(4, ants.WithLogger(log.Default()))
	//Taskpool := NewGoPool(settings.Queuelimit)
	appwindow = app.NewWindow("GO Pixiv")
	authorId := widget.NewEntry()
	illustId := widget.NewEntry()
	illustLabel := widget.NewLabel("Download by IllustId")
	authorLabel := widget.NewLabel("Download all Illusts by AuthorId")
	button1 := widget.NewButton("Download", func() {})
	button1.OnTapped = func() {
		text := illustId.Text
		P.AddTask(func() (interface{}, error) {
			JustDownload(text)
			return nil, nil
		})
		//P.Submit(func() {
		//	JustDownload(text)
		//	return
		//})
		//go pool.Run(func() {
		//	JustDownload(text)
		//})
		illustId.SetText("")
	}
	container.New(layout.NewStackLayout())
	button2 := widget.NewButton("Download", func() {})
	button2.OnTapped = func() {
		text := authorId.Text
		authorId.SetText("")

		button2.Disabled()
		c := make(chan string, 2000)
		//go TaskPool.Run(func() {
		go func() {
			all, err := GetAuthor(statics.StringToInt64(text))
			if err != nil {
				log.Println("Error getting author", err)
				button2.Enable()

				return
			}
			log.Println(text + " pushed TaskQueue")
			log.Println(text + "'s artworks Start download")
			satisfy = 0
			for key, _ := range all {
				k := key
				P.AddTask(func() (interface{}, error) {
					//time.Sleep(1 * time.Second)
					temp := k
					illust, err := work(statics.StringToInt64(temp))
					if err != nil {
						//log.Println(key, " Download failed")
						//continue
						if (err != &NotGood{}) && err != (&AgeLimit{}) {
							c <- temp

						}
						return nil, nil
					}
					illust.Download()
					satisfy++
					return nil, nil
				})
			}
			P.Wait()
			log.Println(text+"'s artworks -> Satisfied and Successfully downloaded illusts: ", satisfy, "in all: ", len(all))
			for len(c) > 0 {
				ss := <-c
				log.Println(ss, " Download failed Now retrying")
				P.AddTask(func() (interface{}, error) {
					JustDownload(ss)
					return nil, nil
				})
			}
			P.Wait()
			close(c)
		}()

		button2.Enable()

	}
	r18 := widget.NewCheck("R-18", func(i bool) {
	})
	r18.SetChecked(settings.Agelimit)
	r18.Refresh()
	likelimit := widget.NewLabel("likelimit")
	readlikelimit := widget.NewEntry()
	cookieLabel := widget.NewLabel("cookie")
	readcookie := widget.NewEntry()
	readlikelimit.SetText(statics.Int64ToString(settings.LikeLimit))
	readcookie.SetText(settings.Cookie)
	readcookie.Refresh()
	Likelimit := container.New(layout.NewGridWrapLayout(fyne.Size{Width: 100, Height: 38}), likelimit, readlikelimit)
	Cookie := container.New(layout.NewGridWrapLayout(fyne.Size{Width: 100, Height: 38}), cookieLabel, readcookie)
	save := widget.NewButton("Save Settings", func() {
		settings.Agelimit = r18.Checked
		to := readlikelimit.Text
		if !statics.AllNum(to) {
			to = "0"
		}
		settings.LikeLimit = statics.StringToInt64(to)
		out, _ := yaml.Marshal(&settings)
		ioutil.WriteFile("settings.yml", out, 0644)
	})
	setting := container.New(layout.NewGridWrapLayout(fyne.Size{Width: 400, Height: 50}), r18, Likelimit, Cookie, save)
	content := container.New(layout.NewGridLayoutWithColumns(3), illustLabel, illustId, button1, authorLabel, authorId, button2)
	//stackqueue := container.NewScroll()
	all := container.NewVBox(content, setting)
	icon, _ := fyne.LoadResourceFromPath("img/icon.ico")
	app.SetIcon(icon)
	appwindow.SetIcon(icon)
	appwindow.SetContent(all)
	appwindow.Resize(fyne.Size{300, 250})
}
func LogInit() {
	T := time.Now()
	logfile := fmt.Sprintf("errorlog/%4d-%2d-%2d.log", T.Year(), T.Month(), T.Day())
	log.SetFlags(log.Ltime)
	f, _ = os.OpenFile(logfile, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	log.SetOutput(f)

}

// TODO ：客户端代理设置 OK
func clinentInit() {
	jfile, _ := os.OpenFile("settings.yml", os.O_RDWR, 0644)
	defer jfile.Close()
	bytevalue, _ := ioutil.ReadAll(jfile)
	yaml.Unmarshal(bytevalue, &settings)
	settings.LikeLimit = max(settings.LikeLimit, 0)
	_, err := os.Stat(settings.Downloadposition)
	if err != nil {
		settings.Downloadposition = "Download"
	}
	settings.Queuelimit = max(settings.Queuelimit, 1)

	out, _ := yaml.Marshal(&settings)
	ioutil.WriteFile("settings.yml", out, 0644)
	_, err = url2.Parse(settings.Proxy)
	log.Println("Check settings:"+settings.Proxy, "PHPSESSID="+settings.Cookie, settings.Downloadposition)
	if err != nil {
		log.Println(err)
		return
	}
}
func GetClient() *http.Client {
	proxyURL, _ := url2.Parse(settings.Proxy)
	return &http.Client{
		Transport: &http.Transport{
			Proxy:                 http.ProxyURL(proxyURL),
			DisableKeepAlives:     true,
			ResponseHeaderTimeout: time.Second * 5,
		},
	}
}
func JustDownload(pid string) {
	illust, _ := work(statics.StringToInt64(pid))
	if illust == nil {
		log.Println(pid, " Download failed")
		return
	}
	log.Println(pid + " Start download")
	illust.Download()
	log.Println(pid + " Finished download")
}
func Get(pid string) {

}
