package mongo

import (
	"encoding/json"
	"fmt"
	"os"

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
}

func (m mgDB) Posts() db.PostsDB {
	return postsDB{col: m.client.Database("blogs").Collection("posts")}

}

type postsDB struct {
	col *mongo.Collection
}

func (p postsDB) Create(po db.Post) error {
	_, err := p.col.InsertOne(context.TODO(), po)
	return err

}
func (p postsDB) FindByFilter(f string) (*[]db.Post, error) {

	fmt.Println("f:" + f)
	flts := &Filter{}
	json.Unmarshal([]byte(f), flts)

	var bdoc bson.D

	for _, a := range flts.Category {
		bdoc = append(bdoc, bson.E{Key: "category", Value: a})
	}
	for _, a := range flts.Tag {
		bdoc = append(bdoc, bson.E{Key: "tags", Value: a})
	}

	var allpost []db.Post

	curr, err := p.col.Find(context.Background(), bdoc)
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
	client, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("MGDBURL")))
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
