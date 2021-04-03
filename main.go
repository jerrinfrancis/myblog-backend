package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	//	"strings"

	"github.com/jerrinfrancis/myblog/posts"
	"github.com/jerrinfrancis/myblog/router"
)

func main() {
	fmt.Println(os.Getenv("MGDBURL"), "test")
	log.Println(os.Getenv("MGDBURL"), "test")
	//os.Setenv("MGDBURL", "mongodb://mongo:27017")

	os.Setenv("MGDBURL", "mongodb://127.0.0.1:27017")

	router := router.NewRouter()
	router.SetHandlerFunc("POST", "/post", posts.Post)
	router.SetHandlerFunc("GET", "/posts", posts.Get)
	router.SetHandlerFunc("DELETE", "/posts", posts.Delete)
	router.SetHandlerFunc("OPTIONS", "/posts", posts.Options)
	router.SetHandlerFunc("OPTIONS", "/post", posts.Options)
	router.SetHandlerFunc("POST", "/uploadfile", posts.UploadFile)
	router.SetHandlerFunc("GET", "/images", posts.GetFile)
	router.SetHandlerFunc("POST", "/category", posts.PostCategory)
	router.SetHandlerFunc("OPTIONS", "/category", posts.Options)
	router.SetHandlerFunc("GET", "/categories", posts.GetCategories)
	router.SetHandlerFunc("PATCH", "/editpost", posts.Update)
	router.SetHandlerFunc("OPTIONS", "/editpost", posts.Options)
	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	//test
	log.Println("Server listening at ", server.Addr)
	server.ListenAndServe()
}
