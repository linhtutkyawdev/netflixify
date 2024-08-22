package server

import (
	"database/sql"
	"os"
	"path/filepath"
	"time"

	"github.com/tursodatabase/go-libsql"
)

func getPosts() ([]Post, error) {
	// Set up the database
	primaryUrl := os.Getenv("TURSO_DATABASE_URL")
	authToken := os.Getenv("TURSO_AUTH_TOKEN")

	dbName := "local.db"
	dir, err := os.MkdirTemp("", "libsql-*")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(dir)

	dbPath := filepath.Join(dir, dbName)
	syncInterval := time.Minute

	connector, err := libsql.NewEmbeddedReplicaConnector(dbPath, primaryUrl,
		libsql.WithAuthToken(authToken),
		libsql.WithSyncInterval(syncInterval),
	)

	if err != nil {
		return nil, err
	}
	defer connector.Close()

	db := sql.OpenDB(connector)
	defer db.Close()

	// Do something with the database
	rows, err := db.Query("SELECT * FROM posts order by video_path desc;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post

	for rows.Next() {
		var post Post

		if err := rows.Scan(&post.Channel_id, &post.Title, &post.Rating, &post.Description, &post.Tags, &post.Video_id, &post.Video_path, &post.Thumbnail_id, &post.Thumbnail_path, &post.G_thumbnail_id, &post.G_thumbnail_path); err != nil {
			return nil, err
		}

		post.Video_path = os.Getenv("TG_API_URL") + "/file/bot" + os.Getenv("BOT_TOKEN") + "/" + post.Video_path
		post.Thumbnail_path = os.Getenv("TG_API_URL") + "/file/bot" + os.Getenv("BOT_TOKEN") + "/" + post.Thumbnail_path
		post.G_thumbnail_path = os.Getenv("TG_API_URL") + "/file/bot" + os.Getenv("BOT_TOKEN") + "/" + post.G_thumbnail_path

		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}
