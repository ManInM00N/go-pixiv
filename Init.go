package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	. "github.com/ManInM00N/go-tool/goruntine"
	"github.com/ManInM00N/go-tool/statics"
	"github.com/tidwall/gjson"
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
var client *http.Client

type FyneLogWriter struct {
	LogText *widget.Entry
}

func (w *FyneLogWriter) Write(p []byte) (n int, err error) {
	message := string(p)
	w.LogText.SetText(w.LogText.Text + message)
	return len(p), nil
}

var pool *GoPool

func windowInit() {
	pool = NewGoPool(settings.Illustqueuelimit)
	app := app.New()
	Taskpool := NewGoPool(settings.Queuelimit)
	appwindow = app.NewWindow("GO Pixiv")
	authorId := widget.NewEntry()
	illustId := widget.NewEntry()
	illustLabel := widget.NewLabel("Download by IllustId")
	authorLabel := widget.NewLabel("Download all Illusts by AuthorId")
	button1 := widget.NewButton("Download", func() {})
	button1.OnTapped = func() {
		text := illustId.Text
		go pool.Run(func() {
			JustDownload(text)
		})
		illustId.SetText("")
	}
	container.New(layout.NewStackLayout())
	button2 := widget.NewButton("Download", func() {})
	button2.OnTapped = func() {
		text := authorId.Text
		go Taskpool.Run(func() {
			var all map[string]gjson.Result
			GetAuthor(statics.StringToInt64(text), &all)
			satisfy := 0
			log.Println(text + "'s artworks Start download")
			for key, _ := range all {
				illust, err := work(statics.StringToInt64(key))
				if err != nil {
					continue
				}
				illust.Download()
				satisfy++
			}
			log.Println(text+"'s artworks -> Satisfied illusts: ", satisfy, "in all: ", len(all))

		})

		authorId.SetText("")
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
	log.SetFlags(log.Ldate | log.Ltime)
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
	proxyURL, err := url2.Parse(settings.Proxy)
	log.Println("Check settings:"+settings.Proxy, "PHPSESSID="+settings.Cookie, settings.Downloadposition)
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
func JustDownload(pid string) {
	illust, _ := work(statics.StringToInt64(pid))
	log.Println(pid + "Start download")
	illust.Download()
	log.Println(pid + "Finished download")
}
