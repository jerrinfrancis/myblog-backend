package main

import (
	"log"
	"net/http"
	"os"

	//	"strings"

	"github.com/jerrinfrancis/myblog/posts"
	"github.com/jerrinfrancis/myblog/router"
)

func main() {
	log.Println("Mongo URL: ", os.Getenv("MY_BLOG_DB_URL"))
	log.Println("IN_PROD", os.Getenv("MY_BLOG_IN_PROD"))
	//os.Setenv("MGDBURL", "mongodb://127.0.0.1:27017")
	//os.Setenv("BLOG_IN_PROD", "X")

	router := router.NewRouter()
	// Default setting in production
	var blogInProd bool
	if os.Getenv("MY_BLOG_IN_PROD") == "NO" {
		blogInProd = false
	} else {
		blogInProd = true
	}

	if !blogInProd {
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
	router.SetHandlerFunc("OPTIONS", "/contactJerrin", posts.Options)

	router.SetHandlerFunc("GET", "/posts", posts.Get)
	router.SetHandlerFunc("GET", "/images", posts.GetFile)
	router.SetHandlerFunc("GET", "/categories", posts.GetCategories)
	router.SetHandlerFunc("POST", "/contactJerrin", posts.SendMessage)

	blogPort := os.Getenv("MY_BLOG_PORT")
	var PORT string
	if len(blogPort) > 0 {
		PORT = ":" + blogPort
	} else {
		PORT = ":" + "8085"
	}

	server := http.Server{
		Addr:    PORT,
		Handler: router,
	}
	//test
	log.Println("Server listening at ", server.Addr)
	server.ListenAndServe()
}
