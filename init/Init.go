package init

import (
	"context"
	"fyne.io/fyne/v2/widget"
	"github.com/ManInM00N/go-tool/goruntine"
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
