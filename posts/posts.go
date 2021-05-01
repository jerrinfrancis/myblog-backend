package posts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

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

type Message struct {
	SenderName  string `json:"senderName"`
	SenderEmail string `json:"senderEmail"`
	Message     string `json:"message"`
}

func SendMessage(w http.ResponseWriter, r *http.Request) {
	c := new(Message)
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Println(c)
	defer r.Body.Close()
	postBody, _ := json.Marshal(map[string]string{
		"senderEmail": c.SenderEmail,
		"senderName":  c.SenderName,
		"message":     c.Message,
	})
	body := bytes.NewBuffer(postBody)
	_, err := http.Post(os.Getenv("MY_BLOG_CONTACTSELF_API"), "application/json", body)
	if err != nil {
		// error to send email should not result in backend service terminating
		log.Printf("An Error Occured %v", err)
	}
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
	w.Header().Add("Content-Type", "image/png")
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
func Options(w http.ResponseWriter, r *http.Request) {
	return
}
func Delete(w http.ResponseWriter, r *http.Request) {
	fmt.Println("testdelete")
	mn := mongo.New()
	var err error
	slug, slugExists := r.Context().Value("param1").(string)
	fmt.Println("Delete", slug)
	if slugExists {
		err = mn.Posts().DeleteBySlug(slug)
	}
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)

	} else {
		w.WriteHeader(http.StatusOK)
	}
}
func Get(w http.ResponseWriter, r *http.Request) {
	mn := mongo.New()
	var posts *[]db.Post
	var post *db.Post
	var err error
	slug, slugExists := r.Context().Value("param1").(string)
	//fmt.Println(slug)
	filters, filterExists := r.Context().Value("filter").(string)

	if filterExists {
		fmt.Println("filtereists")
		posts, err = mn.Posts().FindByFilter(filters)
	} else if slugExists {
		fmt.Println("sligexists")
		post, err = mn.Posts().FindBySlug(slug)
	} else {
		fmt.Println("find all")
		posts, err = mn.Posts().FindAll()
	}
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var bytes []byte
	if slugExists {
		bytes, err = json.Marshal(post)
	} else {
		bytes, err = json.Marshal(posts)
	}
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
	loc, _ := time.LoadLocation("Asia/Kolkata")
	now := time.Now().In(loc)
	fmt.Println(now.Format("11-11-1991"))

	log.Println("Post handler reached")
	mn := mongo.New()
	p := db.Post{}
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()
	if p.Slug == "" {
		p.Slug = strings.ReplaceAll(p.Title, " ", "-")
	}
	err := mn.Posts().Create(p)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
}

func Update(w http.ResponseWriter, r *http.Request) {
	mn := mongo.New()
	p := db.Post{}
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()
	if p.Slug == "" {
		w.Write([]byte("Missign slug"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Println("Reached Update")
	_, err := mn.Posts().UpdateContentBySlug(p.Slug, p.Content, p.ContentPreview)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
