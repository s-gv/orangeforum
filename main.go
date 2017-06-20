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

func getCreds() (string, string) {
	var userName string
	fmt.Printf("Username: ")
	fmt.Scan(&userName)

	fmt.Printf("Password: ")
	password, err := terminal.ReadPassword(int(syscall.Stdin))
	fmt.Printf("\n")
	if err != nil {
		fmt.Printf("[ERROR] %s\n", err)
		return "", ""
	}
	if len(password) < 8 {
		fmt.Printf("[ERROR] Password should have at least 8 characters.\n")
		return "", ""
	}

	fmt.Printf("Password (again): ")
	password2, err := terminal.ReadPassword(int(syscall.Stdin))
	fmt.Printf("\n")
	if err != nil {
		fmt.Printf("[ERROR] %s\n", err)
	}

	pass := string(password)
	pass2 := string(password2)
	if pass != pass2 {
		fmt.Printf("[ERROR] The two psasswords do not match.\n")
		return "", ""
	}

	return userName, pass
}

func main() {
	rand.Seed(time.Now().UnixNano())
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	dbFileName := flag.String("dbname", "orangeforum.db", "Database file path (default: orangeforum.db)")
	port := flag.String("port", "9123", "Port to listen on (default: 9123)")
	shouldMigrate := flag.Bool("migrate", false, "Migrate DB (default: false)")
	createSuperUser := flag.Bool("createsuperuser", false, "Create superuser")
	createUser := flag.Bool("createuser", false, "Create user")

	flag.Parse()

	db.Init("sqlite3", *dbFileName)

	if models.IsMigrationNeeded() {
		if *shouldMigrate {
			models.Migrate()
			return
		} else {
			log.Panicf("[ERROR] DB migration needed.\n")
		}
	} else {
		if *shouldMigrate {
			log.Panicf("[ERROR] DB migration not needed. DB up-to-date.\n")
			return
		}
	}

	if *createSuperUser {
		fmt.Printf("Creating superuser...\n")
		userName, pass := getCreds()
		if userName != "" && pass != "" {
			if err := models.CreateSuperUser(userName, pass); err != nil {
				fmt.Printf("Error creating superuser: %s\n", err)
			}
		}
		return
	}

	if *createUser {
		fmt.Printf("Creating user...\n")
		userName, pass := getCreds()
		if userName != "" && pass != "" {
			if err := models.CreateUser(userName, pass, ""); err != nil {
				fmt.Printf("Error creating superuser: %s\n", err)
			}
		}
		return
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/favicon.ico", views.FaviconHandler)

	mux.HandleFunc("/", views.IndexHandler)

	mux.HandleFunc("/groups/edit", views.GroupEditHandler)
	mux.HandleFunc("/groups/subscribe", views.GroupSubscribeHandler)
	mux.HandleFunc("/groups/unsubscribe", views.GroupUnsubscribeHandler)
	mux.HandleFunc("/groups", views.GroupHandler)

	mux.HandleFunc("/topics/new", views.TopicCreateHandler)
	mux.HandleFunc("/topics/edit", views.TopicUpdateHandler)
	mux.HandleFunc("/topics/subscribe", views.TopicSubscribeHandler)
	mux.HandleFunc("/topics/unsubscribe", views.TopicUnsubscribeHandler)
	mux.HandleFunc("/topics", views.TopicHandler)

	mux.HandleFunc("/comments/new", views.CommentCreateHandler)
	mux.HandleFunc("/comments/edit", views.CommentUpdateHandler)
	mux.HandleFunc("/comments", views.CommentHandler)

	mux.HandleFunc("/test", views.TestHandler)

	mux.HandleFunc("/creategroup", views.CreateGroupHandler)

	mux.HandleFunc("/signup", views.SignupHandler)
	mux.HandleFunc("/login", views.LoginHandler)
	mux.HandleFunc("/logout", views.LogoutHandler)
	mux.HandleFunc("/changepass", views.ChangePasswdHandler)
	mux.HandleFunc("/forgotpass", views.ForgotPasswdHandler)
	mux.HandleFunc("/resetpass", views.ResetPasswdHandler)

	mux.HandleFunc("/note", views.NoteHandler)

	mux.HandleFunc("/admin", views.AdminIndexHandler)

	mux.HandleFunc("/img", views.ImageHandler)

	mux.HandleFunc("/users", views.UserProfileHandler)
	mux.HandleFunc("/users/comments", views.UserCommentsHandler)
	mux.HandleFunc("/users/topics", views.UserTopicsHandler)
	mux.HandleFunc("/users/groups", views.UserGroupsHandler)


	srv := &http.Server{
		Handler: mux,
		Addr: ":" + *port,
		WriteTimeout: 45 * time.Second,
		ReadTimeout:  45 * time.Second,
	}

	log.Println("[INFO] Starting orangeforum on port", *port)
	err := srv.ListenAndServe()
	if err != nil {
		log.Panicf("[ERROR] %s\n", err)
	}

}
