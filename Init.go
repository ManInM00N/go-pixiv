package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	. "github.com/ManInM00N/go-tool/statics"
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
	text := widget.NewEntry()
	button := widget.NewButton("Download", func() {

	})
	button.OnTapped = func() {
		button.Disable()
		illust, err := work(StringToInt64(text.Text))
		if err != nil {
			return
		}
		illust.Download()
		button.Enable()
	}
	ginLog := widget.NewMultiLineEntry()
	content := container.New(layout.NewVBoxLayout(), text, button, container.NewScroll(ginLog))
	appwindow.SetContent(content)
	//out := io.MultiWriter(&FyneLogWriter{LogText: ginLog})

	appwindow.Resize(fyne.Size{1200, 800})
}
func LogInit() {
	log.SetFlags(log.Ldate | log.Ltime)
	f, _ = os.OpenFile("temp.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	log.SetOutput(f)

}
