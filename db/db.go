package db

type DB interface {
	Posts() PostsDB
	Categories() CategoryDB
}
type CategoryDB interface {
	Create(c Category) error
	FindAll() (*[]Category, error)
}
type Category struct {
	CategoryName        string `json:"categoryName"`
	CategoryDescription string `json:"categoryDescription"`
}
type PostsDB interface {
	FindAll() (*[]Post, error)
	Create(p Post) error
	FindByFilter(f string) (*[]Post, error)
	FindBySlug(s string) (*Post, error)
	DeleteBySlug(s string) error
	UpdateContentBySlug(s, c string) (int64, error)
}

type Post struct {
	Slug      string   `json:"slug"`
	Title     string   `json:"title"`
	Content   string   `json:"content"`
	Category  []string `json:"category"`
	Tags      []string `json:"tags"`
	Author    string   `json:"author"`
	Date      string   `json:"date"`
	Published bool     `json:"published"`
}
