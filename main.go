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
	fmt.Println("Mongo URL", os.Getenv("MGDBURL"))

	//os.Setenv("MGDBURL", "mongodb://127.0.0.1:27017")
	//os.Setenv("BLOG_IN_PROD", "X")

	router := router.NewRouter()
	if os.Getenv("BLOG_IN_PROD") != "X" {
		router.SetHandlerFunc("POST", "/post", posts.Post)
		router.SetHandlerFunc("DELETE", "/posts", posts.Delete)
		router.SetHandlerFunc("POST", "/uploadfile", posts.UploadFile)
		router.SetHandlerFunc("POST", "/category", posts.PostCategory)
		router.SetHandlerFunc("PATCH", "/editpost", posts.Update)

	}
	router.SetHandlerFunc("OPTIONS", "/posts", posts.Options)
	router.SetHandlerFunc("OPTIONS", "/post", posts.Options)
	router.SetHandlerFunc("OPTIONS", "/editpost", posts.Options)
	router.SetHandlerFunc("OPTIONS", "/category", posts.Options)
	router.SetHandlerFunc("GET", "/posts", posts.Get)
	router.SetHandlerFunc("GET", "/images", posts.GetFile)
	router.SetHandlerFunc("GET", "/categories", posts.GetCategories)

	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	//test
	log.Println("Server listening at ", server.Addr)
	server.ListenAndServe()
}
