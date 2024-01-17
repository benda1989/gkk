package logger

import (
	"fmt"
	"os"
	"time"

	"gkk/config"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"github.com/weekface/mgorus"
)

var Log *Logger

type Logger struct {
	*logrus.Logger
	Heads []string
}

func (l *Logger) AddHead(args ...string) {
	for _, v := range args {
		l.Heads = append(l.Heads, v)
	}
}

func Init() *Logger {
	L := logrus.New()
	path := "./log/log"
	file, _ := rotatelogs.New(
		path+".%Y%m%d",
		rotatelogs.WithLinkName(path),
		rotatelogs.WithMaxAge(time.Duration(100*24)*time.Hour),        //自动删除
		rotatelogs.WithRotationTime(time.Duration(24*60)*time.Minute), //分割时间
	)
	L.SetOutput(file)
	L.SetFormatter(&logrus.JSONFormatter{})
	mongo(L)
	Log = &Logger{L, []string{}}
	return Log
}

func mongo(l *logrus.Logger) {
	if c := config.GetDefault().Log; c != nil {
		if c.Host != "" {
			hooker, err := mgorus.NewHookerWithAuth(c.Host, c.Db, c.Collection, c.User, c.Password)
			if err == nil {
				l.Hooks.Add(hooker)
			} else {
				fmt.Println("Mongo log err: " + err.Error())
			}
		}
		if c.Mode == "debug" {
			l.SetOutput(os.Stdout)
		}
	}
}
