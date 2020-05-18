package models

import (
	"errors"
	"html/template"

	"github.com/jinzhu/gorm"
)

var (
	errInvalidID = errors.New("Blog: Invalid ID")
)

// BlogModel defines the shape of the blog
type BlogModel struct {
	gorm.Model
	Title string        `gorm:"not null"`
	Body  template.HTML `gorm:"not null"`
}

// BlogService defines the shape of the blogservice
type BlogService struct {
	db *gorm.DB
}

// NewBlogService returns the BlogService struct
func NewBlogService(db *gorm.DB) BlogService {
	return BlogService{
		db: db,
	}
}

// All returns all blogs
func (bs BlogService) All(blogs *[]BlogModel) error {
	if err := bs.db.Find(&blogs).Error; err != nil {
		return err
	}
	return nil
}

// Create creates a resource
func (bs BlogService) Create(blog *BlogModel) error {
	return bs.db.Create(blog).Error
}

// ByID Finds blog by ID
func (bs BlogService) ByID(id string) (*BlogModel, error) {
	var blog BlogModel
	if err := bs.db.Where("id = ?", id).First(&blog).Error; err != nil {
		return nil, err
	}
	return &blog, nil
}

// Update updates the blog
func (bs BlogService) Update(update *BlogModel) error {
	return bs.db.Save(update).Error
}

// Delete removes a blog
func (bs BlogService) Delete(id string) error {
	if id == string(0) {
		return errInvalidID
	}
	return bs.db.Where("id = ?", id).Delete(BlogModel{}).Error
}

// AutoMigrate automatically creates tables
func (bs BlogService) AutoMigrate() error {
	if err := bs.db.AutoMigrate(BlogModel{}).Error; err != nil {
		return err
	}
	return nil
}
