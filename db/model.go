package db

import (
	"database/sql/driver"
	"gorm.io/gorm"
	"time"
)

type Model struct {
	CreatedAt time.Time      `gorm:"autoCreateTime;default:now()" cru:"-"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" cru:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-" cru:"-"`
}

type Time time.Time

func (ct Time) MarshalJSON() ([]byte, error) {
	if t := time.Time(ct); t.IsZero() {
		return []byte(`""`), nil
	} else {
		formattedTime := t.Format("2006-01-02 15:04:05")
		return []byte(`"` + formattedTime + `"`), nil
	}
}

func (ct Time) Value() (driver.Value, error) {
	if t := time.Time(ct); t.IsZero() {
		return nil, nil
	} else {
		return t, nil
	}
}

func ColumnName(name string) string {
	return namer.ColumnName("", name)
}
