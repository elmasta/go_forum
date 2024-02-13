package main

import (
	"database/sql"
	"fmt"
	. "forum/internal/db"
	middle "forum/pkg/middleware"
	models "forum/pkg/models"
	"log"
	"time"
)

func main() {

	db, err := sql.Open("sqlite3", "config/db/forum.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	app := InitDB(db)
	if err := app.Migrate(); err != nil {
		log.Fatal(err)
	}

	var answer string
	fmt.Println("Do you want to create a superuser? y + enter to confirm, enter or any key(s) + enter to cancel.")
	fmt.Scanln(&answer)
	if answer == "y" || answer == "Y" {
		user := &models.User{}
		user.UserType = 3
		user.Validation = 1
		user.CreationDate = time.Now()

		fmt.Println("Enter a name:")
		fmt.Scanln(&answer)
		user.Username = answer

		fmt.Println("Enter an email:")
		fmt.Scanln(&answer)
		user.Email = answer

		fmt.Println("Enter a password:")
		fmt.Scanln(&answer)
		user.Password = answer

		fmt.Println("Confirm the password:")
		fmt.Scanln(&answer)

		if err := middle.CheckAdminRegister(app.DB, answer, user); err != nil {
			if err.Error() == "email already exist" {
				log.Fatal("Email already exists!")
			}
			if err.Error() == "passwords do not match" {
				log.Fatal("Passwords do not match!")
			}
		}

		if err := app.CreateUser(user); err != nil {
			log.Fatal(err)
		}
	}
}
