package main

import (
	"fmt"
	"net/http"

	"practice.blog.com/models"

	"github.com/jinzhu/gorm"

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
	db, err := gorm.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	db.LogMode(true)

	blogC := controllers.NewBlog(db)
	userC := controllers.NewUser(db)
	if err := controllers.AutoMigrate(db, models.BlogModel{}, models.UserModel{}); err != nil {
		panic(err)
	}

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
