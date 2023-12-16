package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	. "github.com/ManInM00N/go-tool/statics"
	"github.com/tidwall/gjson"
	"log"
	"net/http"
	"os"
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

func windowInit() {
	app := app.New()
	appwindow = app.NewWindow("GO Pixiv")
	authorId := widget.NewEntry()
	illustId := widget.NewEntry()
	button1 := widget.NewButton("Download", func() {})
	button1.OnTapped = func() {
		button1.Disable()
		illust, err := work(StringToInt64(illustId.Text))
		if err != nil || illust == nil {
			return
		}
		illust.Download()
		button1.Enable()
	}
	button2 := widget.NewButton("Download", func() {})
	button2.OnTapped = func() {
		button2.Disable()

		var all map[string]gjson.Result
		GetAuthor(StringToInt64(authorId.Text), &all)
		log.Println(len(all))
		for key, _ := range all {
			illust, err := work(StringToInt64(key))
			if err != nil || illust == nil {
				continue
			}
			illust.Download()
		}
		button2.Enable()
	}
	ginLog := widget.NewMultiLineEntry()
	content := container.New(layout.NewGridLayoutWithColumns(2), illustId, button1, container.NewScroll(ginLog), container.NewScroll(ginLog), authorId, button2)
	appwindow.SetContent(content)
	//out := io.MultiWriter(&FyneLogWriter{LogText: ginLog})

	appwindow.Resize(fyne.Size{800, 600})
}
func LogInit() {
	log.SetFlags(log.Ldate | log.Ltime)
	f, _ = os.OpenFile("temp.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	log.SetOutput(f)

}
