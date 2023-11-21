package captcha

import (
	"gkk/cache"
	"strings"
	"time"
)

var c *DigitCaptcha

func init() {
	c = DefaultDigit()
}

func DrawCaptcha() (string, *itemDigit) {
	return c.DrawCaptcha()
}

type CaptStore struct {
	Key    string        `json:"key"`
	Base64 string        `json:"base64"`
	Value  string        `json:"-"`
	Expire time.Duration `json:"-"` //Expire  time.Duration
}

func (c *CaptStore) Set() error {
	return cache.Get().Set(c.Key, c, c.Expire)
}
func (c *CaptStore) Get() error {
	tmp := &CaptStore{}
	if er := cache.Get().Get(c.Key, &tmp); er != nil {
		return er
	}
	c.Value = tmp.Value
	c.Base64 = tmp.Base64
	return nil
}
func (c *CaptStore) Delete() error {
	cache.Get().Delete(c.Key)
	return cache.Get().Delete(strings.Replace(c.Key, "code:", "limit:", 1))
}
