package server

type Channel struct {
	ID       int64
	Title    string
	Password string
}

type Post struct {
	Channel_id       int
	Title            string
	Rating           int
	Description      string
	Tags             string
	Video_id         string
	Video_path       string
	Thumbnail_id     string
	Thumbnail_path   string
	G_thumbnail_id   string
	G_thumbnail_path string
}
