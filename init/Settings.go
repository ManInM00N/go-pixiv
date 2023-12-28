package init

import (
	"context"
	"github.com/ManInM00N/go-tool/goruntine"
	"github.com/devchat-ai/gopool"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
)

type Settings struct {
	Proxy            string `yml:"proxy"`
	Cookie           string `yml:"cookie"`
	Agelimit         bool   `yml:"r-18" `
	Downloadposition string `yml:"downloadposition"`
	LikeLimit        int64  `yml:"minlikelimit"`
	Retry429         int    `yml:"retry429"`
	Downloadinterval int    `yml:"downloadinterval"`
	Retryinterval    int    `yml:"retryinterval"`
}

var (
	ymlfile *os.File
	Setting Settings
	is      = false
)

func UpdateSettings() {
	out, _ := yaml.Marshal(&Setting)
	ioutil.WriteFile("settings.yml", out, 0644)
}

func init() {
	if is {
		return
	}
	is = true
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
	log.Println("Check settings:"+Setting.Proxy, "PHPSESSID="+Setting.Cookie, Setting.Downloadposition)
	UpdateSettings()
	TaskPool = goruntine.NewGoPool(200, 1)
	TaskPool.Run()
	SinglePool = gopool.NewGoPool(1, gopool.WithTaskQueueSize(100))
	P = gopool.NewGoPool(4, gopool.WithTaskQueueSize(5000))
}
