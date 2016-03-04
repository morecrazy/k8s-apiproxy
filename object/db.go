package object

import (
	"third/gorm"
)

type DB struct {
	*gorm.DB
}

func (db *DB) RecordNotFoundError() error {
	return gorm.RecordNotFound
}

func (db *DB) ParseTagSetting(str string) map[string]string {
	return gorm.ParseTagSetting(str)
}
