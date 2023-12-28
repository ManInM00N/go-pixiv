package main

import (
	"errors"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/ManInM00N/go-tool/statics"
	"log"
	. "main/init"
)

type FyneLogWriter struct {
	LogText *widget.Entry
}

func (w *FyneLogWriter) Write(p []byte) (n int, err error) {
	message := string(p)
	w.LogText.SetText(w.LogText.Text + message)
	return len(p), nil
}

var (
	appwindow fyne.Window
)

func windowInit() {
	app := app.New()
	appwindow = app.NewWindow("GO Pixiv")
	authorId := widget.NewEntry()
	illustId := widget.NewEntry()
	illustLabel := widget.NewLabel("Download by IllustId")
	authorLabel := widget.NewLabel("Download all Illusts by AuthorId")
	button1 := widget.NewButton("Download", func() {})
	button1.OnTapped = func() {
		text := illustId.Text
		SinglePool.AddTask(func() (interface{}, error) {
			JustDownload(text, true)
			return nil, nil
		})
		illustId.SetText("")
	}
	container.New(layout.NewStackLayout())
	button2 := widget.NewButton("Download", func() {})
	button2.OnTapped = func() {
		text := authorId.Text
		authorId.SetText("")
		button2.Disabled()
		log.Println(text + " pushed TaskQueue")
		TaskPool.Add(func() {
			if IsClosed {
				return
			}
			c := make(chan string, 2000)
			all, err := GetAuthor(statics.StringToInt64(text))
			if err != nil {
				log.Println("Error getting author", err)
				button2.Enable()

				return
			}
			log.Println(text + "'s artworks Start download")
			satisfy := 0
			for key, _ := range all {
				k := key
				if IsClosed {
					return
				}
				P.AddTask(func() (interface{}, error) {
					//time.Sleep(1 * time.Second)
					if IsClosed {
						return nil, nil
					}
					temp := k
					illust, err := work(statics.StringToInt64(temp), 1)
					if err != nil {
						//log.Println(key, " Download failed")
						//continue
						if !errors.Is(err, &NotGood{}) && !errors.Is(err, &AgeLimit{}) {
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
			for len(c) > 0 {
				if IsClosed {
					return
				}
				ss := <-c
				//log.Println(ss, " Download failed Now retrying")
				P.AddTask(func() (interface{}, error) {
					if a, b := JustDownload(ss, false); b {
						satisfy += a
					}
					return nil, nil
				})
			}
			P.Wait()
			log.Println(text+"'s artworks -> Satisfied and Successfully downloaded illusts: ", satisfy, "in all: ", len(all))
			satisfy = 0
			close(c)
		})
		button2.Enable()

	}
	r18 := widget.NewCheck("R-18", func(i bool) {
	})
	r18.SetChecked(Setting.Agelimit)
	r18.Refresh()
	likelimit := widget.NewLabel("likelimit")
	readlikelimit := widget.NewEntry()
	cookieLabel := widget.NewLabel("cookie")
	readcookie := widget.NewEntry()
	readlikelimit.SetText(statics.Int64ToString(Setting.LikeLimit))
	readcookie.SetText(Setting.Cookie)
	readcookie.Refresh()
	Likelimit := container.New(layout.NewGridWrapLayout(fyne.Size{Width: 100, Height: 38}), likelimit, readlikelimit)
	Cookie := container.New(layout.NewGridWrapLayout(fyne.Size{Width: 100, Height: 38}), cookieLabel, readcookie)
	save := widget.NewButton("Save Settings", func() {
		Setting.Agelimit = r18.Checked
		to := readlikelimit.Text
		if !statics.AllNum(to) {
			to = "0"
		}
		Setting.LikeLimit = statics.StringToInt64(to)
		UpdateSettings()
	})
	setting := container.New(layout.NewGridWrapLayout(fyne.Size{Width: 400, Height: 50}), r18, Likelimit, Cookie, save)
	content := container.New(layout.NewGridLayoutWithColumns(3), illustLabel, illustId, button1, authorLabel, authorId, button2)
	//stackqueue := container.NewScroll()
	all := container.NewVBox(content, setting)
	icon, _ := fyne.LoadResourceFromPath("assets/icon.ico")
	app.SetIcon(icon)
	appwindow.SetIcon(icon)
	appwindow.SetContent(all)
	appwindow.Resize(fyne.Size{300, 250})
}
