package forum

import (
	"database/sql"

	. "forum/pkg/models"

	_ "github.com/mattn/go-sqlite3"
)

type App_db struct {
	DB   *sql.DB
	Data Data
}

func InitDB(db *sql.DB) *App_db {
	return &App_db{DB: db}
}

func (app *App_db) Migrate() error {
	query := `
		CREATE TABLE IF NOT EXISTS users_type(
			id INTEGER PRIMARY KEY AUTOINCREMENT, 
			rank INTEGER NOT NULL,
			label TEXT NOT NULL);

		CREATE TABLE IF NOT EXISTS categories(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			descriptions TEXT NOT NULL,
			creation DATETIME NOT NULL   
		);

		CREATE TABLE IF NOT EXISTS users(
			id INTEGER PRIMARY KEY AUTOINCREMENT, 
			user_type_id INTEGER NOT NULL,
			username TEXT NOT NULL, 
			pwd TEXT NOT NULL,
			email TEXT NOT NULL,
			valid INTEGER NOT NULL,
			asked_mod INTEGER DEFAULT 0,
			creation DATETIME NOT NULL,
			session_token TEXT,
			FOREIGN KEY(user_type_id)REFERENCES users_type(id) ON DELETE CASCADE);
		
		CREATE TABLE IF NOT EXISTS post(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			author_id INTEGER NOT NULL,
			author TEXT NOT NULL,
    		img VARCHAR(255) NOT NULL,
			title TEXT NOT NULL UNIQUE,
			content TEXT NOT NULL,
			creation CURRENT_TIMESTAMP,
			flagged INTEGER DEFAULT 0,
			FOREIGN KEY(author_id) REFERENCES users(id) ON DELETE CASCADE
		);

		CREATE TABLE IF NOT EXISTS link_cat_post(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			category_id INTEGER NOT NULL,
			post_id INTEGER NOT NULL,
			FOREIGN KEY(category_id) REFERENCES categories(id) ON DELETE CASCADE
			FOREIGN KEY(post_id) REFERENCES post(id) ON DELETE CASCADE
		);
		
		CREATE TABLE IF NOT EXISTS link_post(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			post_id INTEGER NOT NULL,
			likes BOOLEAN NOT NULL,
			FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
			FOREIGN KEY(post_id) REFERENCES post(id) ON DELETE CASCADE
		);

		CREATE TABLE IF NOT EXISTS comment(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			author_id INTEGER NOT NULL,
			post_id INTEGER NOT NULL,
			content TEXT NOT NULL,
			img VARCHAR(255) NOT NULL,
			creation CURRENT_TIMESTAMP,
			flagged INTEGER DEFAULT 0,
			FOREIGN KEY(author_id) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY(post_id) REFERENCES post(id) ON DELETE CASCADE
		);
		
		CREATE TABLE IF NOT EXISTS link_comment(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			comment_id INTEGER NOT NULL,
			likes BOOLEAN NOT NULL,
			FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY(comment_id) REFERENCES comment(id) ON DELETE CASCADE
		);
	`
	_, err := app.DB.Exec(query)

	// creation user_type
	var count int
	errChecker := app.DB.QueryRow("SELECT COUNT(*) FROM users_type").Scan(&count)
	if errChecker == sql.ErrNoRows || count == 0 {
		_, err = app.DB.Exec("INSERT INTO users_type(rank, label) VALUES (?,?)", 1, "user")
		if err != nil {
			return err
		}
		_, err = app.DB.Exec("INSERT INTO users_type(rank, label) VALUES (?,?)", 2, "moderator")
		if err != nil {
			return err
		}
		_, err = app.DB.Exec("INSERT INTO users_type(rank, label) VALUES (?,?)", 3, "admin")
		if err != nil {
			return err
		}
		_, err = app.DB.Exec("INSERT INTO users_type(rank, label) VALUES (?,?)", 4, "mod_light")
		if err != nil {
			return err
		}
	}
	return err
}