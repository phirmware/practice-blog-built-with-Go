package main

import (
	"fmt"
	"net/http"

	"practice.blog.com/models"

	"practice.blog.com/controllers"

	"github.com/gorilla/mux"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

const serverPort = ":8080"
const (
	host   = "localhost"
	port   = 5432
	user   = "postgres"
	dbname = "practiceblog"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable", host, port, user, dbname)
	bs, err := models.NewBlogService(psqlInfo)
	if err != nil {
		panic(err)
	}
	us, err := models.NewUserService(psqlInfo)
	if err != nil {
		panic(err)
	}

	if err := bs.AutoMigrate(); err != nil {
		panic(err)
	}
	if err := us.AutoMigrate(); err != nil {
		panic(err)
	}

	blogC := controllers.NewBlog(bs)
	userC := controllers.NewUser(us)

	r := mux.NewRouter()
	r.HandleFunc("/", blogC.Home)
	r.HandleFunc("/post", blogC.Post).Methods("GET")
	r.HandleFunc("/post", blogC.HandlePost).Methods("POST")
	r.HandleFunc("/post/{id}", blogC.PostBYID).Methods("GET")
	r.HandleFunc("/post/delete/{id}", blogC.Delete).Methods("GET")
	r.HandleFunc("/signup", userC.New).Methods("GET")
	r.HandleFunc("/signup", userC.Register).Methods("POST")
	r.HandleFunc("/update/{id}", blogC.Update).Methods("GET")
	r.HandleFunc("/update/{id}", blogC.HandleUpdate).Methods("POST")
	r.HandleFunc("/users", userC.FindAll).Methods("GET")

	fmt.Printf("Server listening at port %s\n", serverPort)
	http.ListenAndServe(serverPort, r)
}
