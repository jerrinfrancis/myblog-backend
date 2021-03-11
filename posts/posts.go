package posts

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"strings"

	"github.com/jerrinfrancis/myblog/db"
	"github.com/jerrinfrancis/myblog/db/mongo"
)

type Image struct {
	ImageUrl string `json:"imageUrl"`
}

type Filter struct {
	Category []string
	Tag      []string
}

func GetFile(w http.ResponseWriter, r *http.Request) {
	filename, ok := r.Context().Value("param1").(string)
	if !ok {
		//no filename error
		return
	}

	file, err := ioutil.ReadFile("temp-images/" + filename)
	if err != nil {
		log.Println(err.Error())
		return
	}
	w.Write(file)

}

func UploadFile(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Writing file")
	r.ParseMultipartForm(10 << 20)
	file, _, err := r.FormFile("myFile")
	if err != nil {
		log.Println(err.Error())
		return
	}

	defer file.Close()
	// check if the directory exists else create it
	_, err = os.Stat("temp-images")
	if err != nil {
		log.Println(err)
		_ = os.Mkdir("temp-images", 0777)
	}

	tempFile, err := ioutil.TempFile("temp-images", "upload-*.png")
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer tempFile.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println(err.Error())
	}
	tempFile.Write(fileBytes)
	filename := strings.Split(tempFile.Name(), "temp-images/")[1]
	imageurl := "http://" + r.Host + "/images/" + filename
	fmt.Println(imageurl)
	fd := &Image{
		ImageUrl: imageurl,
	}
	//log.Println(fd)
	w.Header().Add("Content-Type", "application/json")
	bytes, _ := json.Marshal(fd)
	fmt.Println("test" + r.Host)

	w.Write(bytes)

}

func GetCategories(w http.ResponseWriter, r *http.Request) {
	var categories *[]db.Category
	mn := mongo.New()
	categories, err := mn.Categories().FindAll()
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	bytes, err := json.Marshal(categories)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(bytes)
	return
}

func Get(w http.ResponseWriter, r *http.Request) {
	mn := mongo.New()
	var posts *[]db.Post
	var err error
	filters, ok := r.Context().Value("filter").(string)
	if ok {
		posts, err = mn.Posts().FindByFilter(filters)
	} else {
		posts, err = mn.Posts().FindAll()
	}
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	bytes, err := json.Marshal(posts)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return

	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(bytes)

	return
}
func PostCategory(w http.ResponseWriter, r *http.Request) {
	mn := mongo.New()
	c := db.Category{}
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()
	err := mn.Categories().Create(c)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}
func Post(w http.ResponseWriter, r *http.Request) {
	log.Println("Post handler reached")
	mn := mongo.New()
	p := db.Post{}
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()
	err := mn.Posts().Create(p)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
}
