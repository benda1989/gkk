package logger

import (
	"fmt"
	"runtime"
	"strings"
	"sync"
	"time"

	"gkk/config"
)

var errs = &errFunc{lock: new(sync.Mutex)}

func ErrorRegister(f func(string)) {
	errs.register(f, false)
}
func ErrorRegisterOnly1(f func(string)) {
	errs.register(f, true)
}
func ErrorWx(e ...string) {
	errs.run(strings.Join(e, " "))
	Log.Error(e)
}

type errFunc struct {
	f        []func(string)
	lastTime int64
	lock     *sync.Mutex
}

func (e *errFunc) register(f func(string), only bool) {
	if only {
		e.f = []func(string){f}
	} else {
		e.f = append(e.f, f)
	}
}

func (e *errFunc) run(logs string) {
	if len(e.f) == 0 {
		return
	}
	if time.Now().Unix()-e.lastTime > 1 {
		e.lock.Lock()
		e.lastTime = time.Now().Unix()
		defer e.lock.Unlock()
		for _, v := range e.f {
			v(logs)
		}
	}
}

func StackSend(skip int, e string) {
	e += "\n" + Stack(skip)
	if config.IsDebug {
		fmt.Print(e)
	}
	errs.run(e)
}

func Stack(skip int) (re string) {
	for i := skip; ; i++ {
		_, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		if !strings.Contains(file, "local/go/src") && !strings.Contains(file, "/go/pkg") {
			logs := fmt.Sprintf("%s:%d\n", file, line)
			re += logs
		}
	}
	return
}
