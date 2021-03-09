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

// var mux map[string]map[string]func(http.ResponseWriter, *http.Request)

// type myHandler struct{}

// func (h *myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	if hm, ok := mux[r.Method]; ok {
// 		if h, ok := hm[r.URL.String()]; ok {
// 			h(w, r)
// 			return
// 		} else {
// 			if strings.HasPrefix(r.URL.String(), "/images/") {
// 				if h, ok = hm["/getfile"]; ok {
// 					h(w, r)
// 					log.Println(r.URL.String())
// 					return
// 				}
// 			}
// 		}
// 	}
// }

func main() {
	fmt.Println(os.Getenv("JERTEST"))
	//os.Setenv("MGDBURL", "mongodb://mongo:27017")

	//	os.Setenv("MGDBURL", "mongodb://127.0.0.1:27017")
	router := router.NewRouter()
	router.SetHandlerFunc("POST", "/post", posts.Post)
	router.SetHandlerFunc("GET", "/posts", posts.Get)
	router.SetHandlerFunc("POST", "/uploadfile", posts.UploadFile)
	router.SetHandlerFunc("GET", "/images", posts.GetFile)

	server := http.Server{
		Addr: ":8080",
		//	Handler: &myHandler{},
		Handler: router,
	}

	// mux = make(map[string]map[string]func(http.ResponseWriter, *http.Request))
	// mux["GET"] = make(map[string]func(http.ResponseWriter, *http.Request))
	// mux["POST"] = make(map[string]func(http.ResponseWriter, *http.Request))
	// mux["GET"]["/posts"] = posts.Get
	// mux["GET"]["/getfile"] = posts.GetFile
	// mux["POST"]["/post"] = posts.Post
	// mux["POST"]["/uploadfile"] = posts.UploadFile
	log.Println("Server listening at ", server.Addr)
	server.ListenAndServe()
}
