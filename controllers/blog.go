package controllers

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/jinzhu/gorm"

	"practice.blog.com/models"

	"practice.blog.com/views"
)

// Blog defines the shape of the blog
type Blog struct {
	HomeView       views.View
	PostView       views.View
	SinglePostView views.View
	UpdateView     views.View
	bs             models.BlogService
}

// BlogForm defines the shape of the blogform
type BlogForm struct {
	Title string        `schema:"title"`
	Body  template.HTML `schema:"body"`
}

// NewBlog returns the Blog struct
func NewBlog(db *gorm.DB) Blog {
	return Blog{
		HomeView:       views.NewView("bootstrap", "blog/home"),
		PostView:       views.NewView("bootstrap", "blog/new"),
		SinglePostView: views.NewView("bootstrap", "blog/post"),
		UpdateView:     views.NewView("bootstrap", "blog/update"),
		bs:             models.NewBlogService(db),
	}
}

// Home handles / PATH
func (b Blog) Home(w http.ResponseWriter, r *http.Request) {
	var blogs []models.BlogModel
	if err := b.bs.All(&blogs); err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", blogs)
	b.HomeView.Render(w, blogs)
}

// Post handles /post PATH
func (b Blog) Post(w http.ResponseWriter, r *http.Request) {
	b.PostView.Render(w, nil)
}

// HandlePost handles /post POST PATH
func (b Blog) HandlePost(w http.ResponseWriter, r *http.Request) {
	var form BlogForm
	if err := ParseForm(r, &form); err != nil {
		panic(err)
	}
	blog := models.BlogModel{
		Title: form.Title,
		Body:  form.Body,
	}
	if err := b.bs.Create(&blog); err != nil {
		panic(err)
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

// PostBYID finds a single post by ID
func (b Blog) PostBYID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	data, err := b.bs.ByID(id)
	if err != nil {
		panic(err)
	}
	b.SinglePostView.Render(w, data)
}

// Update renders the update form
func (b Blog) Update(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	blog, err := b.bs.ByID(id)
	if err != nil {
		panic(err)
	}
	b.UpdateView.Render(w, blog)
}

// HandleUpdate updates a blog
func (b Blog) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var form BlogForm
	blog, err := b.bs.ByID(id)
	if err != nil {
		panic(err)
	}
	if err := ParseForm(r, &form); err != nil {
		panic(err)
	}
	blog.Title = form.Title
	blog.Body = form.Body
	if err := b.bs.Update(blog); err != nil {
		panic(err)
	}
	http.Redirect(w, r, "/post/"+strconv.FormatUint(uint64(blog.ID), 10), http.StatusFound)
}

// Delete removes a blog from the database
func (b Blog) Delete(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if err := b.bs.Delete(id); err != nil {
		panic(err)
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

// AutoMigrate calls the automigrate func on blogservice
func (b Blog) AutoMigrate() error {
	return b.bs.AutoMigrate()
}
