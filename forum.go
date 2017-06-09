package main

import (
	"net/http"
	"log"
	"flag"
	"github.com/s-gv/orangeforum/models/db"
	"time"
	"math/rand"
	"github.com/s-gv/orangeforum/models"
	"github.com/s-gv/orangeforum/views"
	"golang.org/x/crypto/ssh/terminal"
	"fmt"
	"syscall"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	dbFileName := flag.String("dbname", "orangeforum.db", "Database file path (default: orangeforum.db)")
	port := flag.String("port", "9123", "Port to listen on (default: 9123)")
	shouldMigrate := flag.Bool("migrate", false, "Migrate DB (default: false)")
	createSuperUser := flag.Bool("createsuperuser", false, "Create superuser")

	flag.Parse()

	db.Init("sqlite3", *dbFileName)

	if models.IsMigrationNeeded() {
		if *shouldMigrate {
			models.Migrate()
			return
		} else {
			log.Fatalf("[ERROR] DB migration needed.\n")
		}
	} else {
		if *shouldMigrate {
			log.Fatalf("[ERROR] DB migration not needed. DB up-to-date.\n")
			return
		}
	}

	if *createSuperUser {
		var userName string
		fmt.Printf("Username: ")
		fmt.Scan(&userName)

		fmt.Printf("Password: ")
		password, err := terminal.ReadPassword(int(syscall.Stdin))
		fmt.Printf("\n")
		if err != nil {
			log.Fatalf("[ERROR] Error creating super user: %s\n", err)
		}
		if len(password) < 8 {
			fmt.Printf("Password should have at least 8 characters.\n")
			return
		}

		fmt.Printf("Password (again): ")
		password2, err := terminal.ReadPassword(int(syscall.Stdin))
		fmt.Printf("\n")
		if err != nil {
			log.Fatalf("[ERROR] Error creating super user: %s\n", err)
		}

		pass := string(password)
		pass2 := string(password2)
		if pass != pass2 {
			fmt.Printf("The two psasswords do not match.\n")
			return
		}

		if err := models.CreateSuperUser(userName, pass); err != nil {
			fmt.Printf("Error creating superuser: %s\n", err)
		}
		return
	}


	http.HandleFunc("/", views.IndexHandler)
	http.HandleFunc("/test", views.TestHandler)

	http.HandleFunc("/creategroup", views.CreateGroupHandler)

	http.HandleFunc("/signup", views.SignupHandler)
	http.HandleFunc("/login", views.LoginHandler)
	http.HandleFunc("/logout", views.LogoutHandler)
	http.HandleFunc("/changepass", views.ChangePasswdHandler)
	http.HandleFunc("/forgotpass", views.ForgotPasswdHandler)
	http.HandleFunc("/resetpass", views.ResetPasswdHandler)

	http.HandleFunc("/admin", views.AdminIndexHandler)

	log.Println("[INFO] Starting orangeforum on port", *port)
	http.ListenAndServe(":" + *port, nil)
}
