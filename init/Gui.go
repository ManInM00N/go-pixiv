package init

import (
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/ManInM00N/go-tool/statics"
	. "main/DAO"
)

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
	process := widget.NewProgressBar()
	process.Min = 0
	waitingtasks := 0
	gridLayout := layout.NewGridLayoutWithColumns(3)
	waitingtasksLabel := widget.NewLabel("There is no tasks waiting")
	waitingtasksLabel.TextStyle.Bold = true
	waitingtasksLabel.TextStyle.TabWidth = 16
	TasknameLabel := widget.NewLabel("No Task in queue")
	TasknameLabel.TextStyle.Bold = true
	TasknameLabel.TextStyle.TabWidth = 16
	Process := container.New(gridLayout,
		TasknameLabel,
		process,
		waitingtasksLabel,

		//widget.NewLabel(""),
	)
	//gridWithColumns := container.NewGridWithColumns(2, Process.Objects...)
	// 将 "Control 2" 添加到具有两列的网格容器
	//gridWithColumns.AddObject(Process.Objects[2])
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
		waitingtasks++

		TaskPool.Add(func() {
			if IsClosed {
				return
			}
			c := make(chan string, 2000)
			all, err := GetAuthor(statics.StringToInt64(text))
			waitingtasks--
			if err != nil {
				DebugLog.Println("Error getting author", err)
				return
			}
			if waitingtasks > 0 {
				waitingtasksLabel.SetText("There are " + fmt.Sprintf("%d", waitingtasks) + " waiting tasks")
			} else {
				waitingtasksLabel.SetText("There is no tasks waiting")
			}
			waitingtasksLabel.Refresh()
			TasknameLabel.SetText(text + " are downloading:")
			TasknameLabel.Refresh()
			process.Max = float64(len(all))
			process.Value = 0
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
						process.Value++
						process.Refresh()

						return nil, nil
					}
					Download(illust)
					satisfy++
					process.Value++
					process.Refresh()

					return nil, nil
				})
			}
			P.Wait()
			TasknameLabel.SetText("Now Recheck " + text)
			TasknameLabel.Refresh()
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
			process.SetValue(0)
			process.Refresh()
		})
		if waitingtasks > 0 {
			waitingtasksLabel.SetText("There are " + fmt.Sprintf("%d", waitingtasks) + " waiting tasks")
		} else {
			waitingtasksLabel.SetText("There is no tasks waiting")
		}
		waitingtasksLabel.Refresh()

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
	all := container.NewVBox(content, Process, setting)
	icon, _ := fyne.LoadResourceFromPath("assets/icon.ico")
	app.SetIcon(icon)
	Appwindow.SetIcon(icon)
	Appwindow.SetContent(all)
	Appwindow.Resize(fyne.Size{Width: 300, Height: 250})
	//Appwindow.SetCloseIntercept(func() {
	//
	//	//app.Quit()
	//})
}
