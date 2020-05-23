package views

import (
	"html/template"
	"net/http"
	"path/filepath"
)

const (
	fileExt  = ".html"
	filePath = "views/"
)

// View defines the shape of the view struct
type View struct {
	t      *template.Template
	layout string
}

func layoutFiles() []string {
	files, err := filepath.Glob("views/layout/*")
	if err != nil {
		panic(err)
	}
	return files
}

func appendExt(files []string) {
	for i, f := range files {
		files[i] = filePath + f + fileExt
	}
}

// NewView returns the new view
func NewView(layout string, files ...string) View {
	appendExt(files)
	files = append(files, layoutFiles()...)
	t, err := template.ParseFiles(files...)
	if err != nil {
		panic(err)
	}
	return View{
		t:      t,
		layout: layout,
	}
}

// Render renders the html page specified
func (v View) Render(w http.ResponseWriter, data interface{}) error {
	if _, b := data.(Data); b {
	} else {
		data = &Data{
			Yield: data,
		}
	}
	return v.t.ExecuteTemplate(w, v.layout, data)
}
