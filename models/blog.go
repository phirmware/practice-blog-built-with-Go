package models

import (
	"errors"
	"html/template"
	"strconv"

	"github.com/jinzhu/gorm"
)

var (
	errInvalidID    = errors.New("Blog: Invalid ID")
	errTitleMissing = errors.New("Blog: A blog title is missing")
	errBodyMissing  = errors.New("Blog: A blog body is missing")
)

// BlogDB defines the blog type
type BlogDB interface {
	All(*[]BlogModel) error
	Create(*BlogModel) error
	ByID(id string) (*BlogModel, error)
	Update(*BlogModel) error
	Delete(id string) error
	AutoMigrate() error
}

var _ BlogDB = BlogService{}

// BlogModel defines the shape of the blog
type BlogModel struct {
	gorm.Model
	Title string        `gorm:"not null"`
	Body  template.HTML `gorm:"not null"`
}

// BlogService defines the shape of the blogservice
type BlogService struct {
	BlogDB
}

type blogGorm struct {
	db *gorm.DB
}

type blogValidation struct {
	BlogDB
}

type blogValidationFn func(blog *BlogModel) error

func runBlogValidationFns(blog *BlogModel, fns ...blogValidationFn) error {
	for _, fn := range fns {
		if err := fn(blog); err != nil {
			return err
		}
	}
	return nil
}

// ************************* Blog Validation Functions ************************** //

func (bv *blogValidation) checkForMissingFields(blog *BlogModel) error {
	if blog.Title == "" {
		return errTitleMissing
	}
	if blog.Body == "" {
		return errBodyMissing
	}
	return nil
}

func (bv *blogValidation) invalidID(blog *BlogModel) error {
	if blog.ID == 0 {
		return errInvalidID
	}
	return nil
}

// ************************* Blog Validation Functions ************************** //

func newBlogGorm(connectionInfo string) (*blogGorm, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	db.LogMode(true)
	if err != nil {
		return nil, err
	}
	return &blogGorm{
		db: db,
	}, nil
}

func newBlogValidation(bg *blogGorm) *blogValidation {
	return &blogValidation{
		BlogDB: bg,
	}
}

// NewBlogService returns the BlogService struct
func NewBlogService(connectionInfo string) (*BlogService, error) {
	bg, err := newBlogGorm(connectionInfo)
	bv := newBlogValidation(bg)
	if err != nil {
		return nil, err
	}

	return &BlogService{
		BlogDB: bv,
	}, nil
}

// All returns all blogs
func (bg blogGorm) All(blogs *[]BlogModel) error {
	if err := bg.db.Find(&blogs).Error; err != nil {
		return err
	}
	return nil
}

// *********** Create Methods ************ //

func (bv *blogValidation) Create(blog *BlogModel) error {
	if err := runBlogValidationFns(blog, bv.checkForMissingFields); err != nil {
		return err
	}
	return bv.BlogDB.Create(blog)
}

// Create creates a resource
func (bg blogGorm) Create(blog *BlogModel) error {
	return bg.db.Create(blog).Error
}

// ********** Create Methods ************ //

// ByID Finds blog by ID
func (bg blogGorm) ByID(id string) (*BlogModel, error) {
	var blog BlogModel
	if err := bg.db.Where("id = ?", id).First(&blog).Error; err != nil {
		return nil, err
	}
	return &blog, nil
}

// ************* Update Methods ***************** //
// Update updates the blog
func (bg blogGorm) Update(update *BlogModel) error {
	return bg.db.Save(update).Error
}

func (bv *blogValidation) Update(update *BlogModel) error {
	if err := runBlogValidationFns(update, bv.checkForMissingFields); err != nil {
		return err
	}
	return bv.BlogDB.Update(update)
}

// ************* Update Methods ***************** //

// ************* Delete Methods ***************** //
// Delete removes a blog
func (bg blogGorm) Delete(id string) error {
	return bg.db.Where("id = ?", id).Delete(BlogModel{}).Error
}

func (bv *blogValidation) Delete(id string) error {
	var blog BlogModel
	int, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return err
	}
	blog.ID = uint(int)
	if err := runBlogValidationFns(&blog, bv.invalidID); err != nil {
		return err
	}
	return bv.BlogDB.Delete(id)
}

// ************* Delete Methods ***************** //

// AutoMigrate automatically creates tables
func (bg blogGorm) AutoMigrate() error {
	return bg.db.AutoMigrate(BlogModel{}).Error
}
