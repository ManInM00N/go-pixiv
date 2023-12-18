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
	"os"
	"time"
)

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
	pool = NewGoPool(8)
	app := app.New()
	appwindow = app.NewWindow("GO Pixiv")
	authorId := widget.NewEntry()
	illustId := widget.NewEntry()
	illustLabel := widget.NewLabel("Download by IllustId")
	authorLabel := widget.NewLabel("Download all Illusts by AuthorId")
	button1 := widget.NewButton("Download", func() {})
	button1.OnTapped = func() {
		button1.Disable()
		illust, err := work(statics.StringToInt64(illustId.Text))
		if err != nil || illust == nil {
			return
		}
		pool.Run(illust.Download)
		illustId.SetText("")
		button1.Enable()
	}

	button2 := widget.NewButton("Download", func() {})
	button2.OnTapped = func() {
		button2.Disable()

		var all map[string]gjson.Result
		GetAuthor(statics.StringToInt64(authorId.Text), &all)
		log.Println(len(all))
		for key, _ := range all {
			illust, err := work(statics.StringToInt64(key))
			if err != nil || illust == nil {
				continue
			}
			pool.Run(illust.Download)
		}
		authorId.SetText("")
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
