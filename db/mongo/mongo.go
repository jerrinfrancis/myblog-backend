package mongo

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"time"

	"log"

	"github.com/jerrinfrancis/myblog/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
)

type mgDB struct {
	client *mongo.Client
}
type Filter struct {
	Category []string
	Tag      []string
	Limit    []string
}

func (m mgDB) Posts() db.PostsDB {
	return postsDB{col: m.client.Database("blogs").Collection("posts")}

}
func (m mgDB) Categories() db.CategoryDB {
	return categoryDB{col: m.client.Database("blogs").Collection("categories")}

}

type categoryDB struct {
	col *mongo.Collection
}

func (c categoryDB) Create(ca db.Category) error {
	_, err := c.col.InsertOne(context.TODO(), ca)
	return err
}

func (c categoryDB) FindAll() (*[]db.Category, error) {
	var allCategories []db.Category
	curr, err := c.col.Find(context.Background(), bson.D{})
	if err != nil {
		return nil, err
	}
	err = curr.All(context.Background(), &allCategories)
	if err != nil {
		return nil, err
	}
	return &allCategories, nil

}

type postsDB struct {
	col *mongo.Collection
}

func (p postsDB) Create(po db.Post) error {
	_, err := p.col.InsertOne(context.TODO(), po)
	return err

}
func (p postsDB) UpdateContentBySlug(s, content, contentPreview string) (int64, error) {
	var bdoc, set bson.D
	bdoc = append(bdoc, bson.E{Key: "$set", Value: append(set, bson.E{Key: "content", Value: content}, bson.E{Key: "contentPreview", Value: contentPreview})})
	result, error := p.col.UpdateOne(context.Background(), bson.M{"slug": s}, bdoc)
	if error != nil {
		return 0, error
	}
	return result.MatchedCount, nil

}
func (p postsDB) DeleteBySlug(s string) error {
	var bdoc bson.D
	bdoc = append(bdoc, bson.E{Key: "slug", Value: s})
	_, error := p.col.DeleteOne(context.Background(), bdoc)
	if error != nil {
		return error
	}
	return nil
}
func (p postsDB) FindBySlug(s string) (*db.Post, error) {
	var bdoc bson.D
	bdoc = append(bdoc, bson.E{Key: "slug", Value: s})
	var post db.Post
	error := p.col.FindOne(context.Background(), bdoc).Decode(&post)
	if error != nil {
		return nil, error
	}

	return &post, nil
}
func (p postsDB) FindByFilter(f string) (*[]db.Post, error) {

	fmt.Println("f:" + f)
	flts := &Filter{}
	var bdoc = bson.D{}
	//var filterExists bool

	json.Unmarshal([]byte(f), flts)

	for _, a := range flts.Category {
		bdoc = append(bdoc, bson.E{Key: "category", Value: a})
		//filterExists = true
	}
	for _, a := range flts.Tag {
		bdoc = append(bdoc, bson.E{Key: "tags", Value: a})
		//filterExists = true
	}
	options := options.Find()

	for _, a := range flts.Limit {
		log.Println("Limit", a)
		options.SetSort(bson.D{{"_id", -1}})
		limit, _ := strconv.ParseInt(a, 10, 64)
		options.SetLimit(limit)
	}
	var allpost []db.Post

	curr, err := p.col.Find(context.Background(), bdoc, options)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	for curr.Next(context.Background()) {
		var result db.Post
		err = curr.Decode(&result)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		allpost = append(allpost, result)
	}

	return &allpost, err

}
func (p postsDB) FindAll() (*[]db.Post, error) {
	var allpost []db.Post

	//	curr, err := p.col.Find(context.Background(), bson.D{{"category", "politics"}})

	curr, err := p.col.Find(context.Background(), bson.D{})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	for curr.Next(context.Background()) {
		var result db.Post
		err = curr.Decode(&result)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		allpost = append(allpost, result)

	}

	return &allpost, err

}

var client *mongo.Client

func New() db.DB {
	if client != nil {
		return mgDB{client: client}
	}
	blogDBURL := os.Getenv("MY_BLOG_DB_URL")
	if len(blogDBURL) == 0 {
		blogDBURL = "mongodb://127.0.0.1:27017"

	}
	log.Println("creating client :", blogDBURL)
	client, err := mongo.NewClient(options.Client().ApplyURI(blogDBURL))
	if err != nil {
		log.Fatalln(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	if err = client.Connect(ctx); err != nil {
		cancel()
		log.Fatalln("unable to connect to DB", err)
	}
	defer cancel()

	return mgDB{client: client}
}
