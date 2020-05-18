package controllers

import (
	"net/http"

	"github.com/jinzhu/gorm"

	"github.com/gorilla/schema"
)

// ParseForm parses the form values using schema package
func ParseForm(r *http.Request, dst interface{}) error {
	if err := r.ParseForm(); err != nil {
		return err
	}
	dec := schema.NewDecoder()
	return dec.Decode(dst, r.PostForm)
}

// AutoMigrate automatically adds table to db
func AutoMigrate(db *gorm.DB, model ...interface{}) error {
	return db.AutoMigrate(model...).Error
}
