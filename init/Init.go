package init

import (
	"context"
	"errors"
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
	. "main/DAO"
	"os"
)

var (
	ymlfile *os.File
	Setting Settings
	is      = false
)

func UpdateSettings() {
	out, _ := yaml.Marshal(&Setting)
	ioutil.WriteFile("settings.yml", out, 0644)
}

type FyneLogWriter struct {
	LogText *widget.Entry
}

func (w *FyneLogWriter) Write(p []byte) (n int, err error) {
	message := string(p)
	w.LogText.SetText(w.LogText.Text + message)
	return len(p), nil
}

var (
	Appwindow fyne.Window
)

func WindowInit() {
	app := app.New()
	Appwindow = app.NewWindow("GO Pixiv")
	authorId := widget.NewEntry()
	illustId := widget.NewEntry()
	illustLabel := widget.NewLabel("Download by IllustId")
	authorLabel := widget.NewLabel("Download all Illusts by AuthorId")
	button1 := widget.NewButton("Download", func() {})
	button1.OnTapped = func() {
		text := illustId.Text
		text = statics.CatchNumber(text)
		if text == "" {
			return
		}
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
		text = statics.CatchNumber(text)
		if text == "" {
			return
		}
		authorId.SetText("")
		button2.Disabled()
		InfoLog.Println(text + " pushed TaskQueue")
		TaskPool.Add(func() {
			if IsClosed {
				return
			}
			c := make(chan string, 2000)
			all, err := GetAuthor(statics.StringToInt64(text))
			if err != nil {
				DebugLog.Println("Error getting author", err)
				button2.Enable()

				return
			}
			InfoLog.Println(text + "'s artworks Start download")
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
					Download(illust)
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
			InfoLog.Println(text+"'s artworks -> Satisfied and Successfully downloaded illusts: ", satisfy, "in all: ", len(all))
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
	Appwindow.SetIcon(icon)
	Appwindow.SetContent(all)
	Appwindow.Resize(fyne.Size{300, 250})
}

func init() {
	if is {
		return
	}
	is = true
	Log_init()
	Ctx, Cancel = context.WithCancel(context.Background())
	ymlfile, _ = os.OpenFile("settings.yml", os.O_RDWR, 0644)
	defer ymlfile.Close()
	bytevalue, _ := ioutil.ReadAll(ymlfile)
	yaml.Unmarshal(bytevalue, &Setting)
	Setting.LikeLimit = max(Setting.LikeLimit, 0)
	_, err := os.Stat(Setting.Downloadposition)
	if err != nil {
		Setting.Downloadposition = "Download"
	}
	Setting.Retry429 = max(Setting.Retry429, 3000)
	Setting.Retryinterval = max(Setting.Retryinterval, 200)
	Setting.Downloadinterval = max(Setting.Downloadinterval, 100)
	DebugLog.Println("Check settings:"+Setting.Proxy, "PHPSESSID="+Setting.Cookie, Setting.Downloadposition)
	UpdateSettings()
	TaskPool = goruntine.NewGoPool(200, 1)
	TaskPool.Run()
	SinglePool = gopool.NewGoPool(1, gopool.WithTaskQueueSize(100))
	P = gopool.NewGoPool(4, gopool.WithTaskQueueSize(5000))
	WindowInit()

}
