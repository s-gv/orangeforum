// Copyright (c) 2017 Sagar Gubbi. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
    "github.com/eyedeekay/sam-forwarder/config"
	"github.com/s-gv/orangeforum/models"
	"github.com/s-gv/orangeforum/models/db"
	"github.com/s-gv/orangeforum/views"
	"golang.org/x/crypto/ssh/terminal"
	"log"
	"math/rand"
	"net/http"
	"net/http/fcgi"
	"syscall"
	"time"
)

func getCreds() (string, string) {
	var userName string
	fmt.Printf("Username: ")
	fmt.Scanf("%s\n", &userName)

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

	dsn := flag.String("dsn", "orangeforum.db", "Data source name")
	dbDriver := flag.String("dbdriver", "sqlite3", "DB driver name")
	addr := flag.String("addr", ":9123", "Port to listen on")
	shouldMigrate := flag.Bool("migrate", false, "Migrate DB")
	createSuperUser := flag.Bool("createsuperuser", false, "Create superuser (interactive)")
	createUser := flag.Bool("createuser", false, "Create user. Optional arguments: <username> <password> <email>")
	changePasswd := flag.Bool("changepasswd", false, "Change password")
	deleteSessions := flag.Bool("deletesessions", false, "Delete all sessions (logout all users)")
	fcgiMode := flag.Bool("fcgi", false, "Fast CGI rather than listening on a port")
    usei2p := flag.Bool("usei2p", false, "Forward the service to the i2p network as an eepSite")
    i2pconf := flag.String("i2pini", "./contrib/tunnels.orangeforum.conf", "i2p tunnel configuration file to use")

	flag.Parse()

    if *usei2p {
        if i2pforwarder, i2perr := i2ptunconf.NewSAMForwarderFromConfig(*i2pconf, "127.0.0.1", "7656"); i2perr != nil {
            fmt.Printf("Error creating i2p tunnel from config, %s", i2perr.Error())
            return
        }else{
            *addr = i2pforwarder.Target()
            fmt.Printf("Serving eepSite on, %s", i2pforwarder.Base32())
            go i2pforwarder.Serve()
        }
    }

	db.Init(*dbDriver, *dsn)

	if *shouldMigrate {
		models.Migrate()
		return
	}

	if models.IsMigrationNeeded() {
		log.Panicf("[ERROR] DB migration needed.\n")
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
		args := flag.Args()
		var username, passwd, email string
		if len(args) >= 2 {
			username = args[0]
			passwd = args[1]
		} else {
			username, passwd = getCreds()
		}
		if len(args) >= 3 {
			email = args[2]
		}
		if username != "" && passwd != "" {
			if err := models.CreateUser(username, passwd, email); err != nil {
				fmt.Printf("Error creating user: %s\n", err)
			}
		} else {
			fmt.Printf("Error: Username and password cannot be blank.\n")
		}
		return
	}

	if *changePasswd {
		userName, pass := getCreds()
		if userName != "" && pass != "" {
			if err := models.UpdateUserPasswd(userName, pass); err != nil {
				fmt.Printf("Error changing password: %s\n", err)
			}
		}
		return
	}

	if *deleteSessions {
		db.Exec(`DELETE FROM sessions;`)
		return
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", views.IndexHandler)

	mux.HandleFunc("/static/css/orangeforum.css", views.StyleHandler)

	mux.HandleFunc("/static/js/orangeforum.js", views.ScriptHandler)

	mux.HandleFunc("/favicon.ico", views.FaviconHandler)

	mux.HandleFunc("/img", views.ImageHandler)

	mux.HandleFunc("/note", views.NoteHandler)

	mux.HandleFunc("/admin", views.AdminIndexHandler)

	mux.HandleFunc("/pm", views.PrivateMessageHandler)
	mux.HandleFunc("/pm/new", views.PrivateMessageCreateHandler)
	mux.HandleFunc("/pm/delete", views.PrivateMessageDeleteHandler)

	mux.HandleFunc("/groups/edit", views.GroupEditHandler)
	mux.HandleFunc("/groups/subscribe", views.GroupSubscribeHandler)
	mux.HandleFunc("/groups/unsubscribe", views.GroupUnsubscribeHandler)
	mux.HandleFunc("/groups", views.GroupIndexHandler)

	mux.HandleFunc("/topics/new", views.TopicCreateHandler)
	mux.HandleFunc("/topics/edit", views.TopicUpdateHandler)
	mux.HandleFunc("/topics/subscribe", views.TopicSubscribeHandler)
	mux.HandleFunc("/topics/unsubscribe", views.TopicUnsubscribeHandler)
	mux.HandleFunc("/topics", views.TopicIndexHandler)

	mux.HandleFunc("/comments/new", views.CommentCreateHandler)
	mux.HandleFunc("/comments/edit", views.CommentUpdateHandler)
	mux.HandleFunc("/comments", views.CommentIndexHandler)

	mux.HandleFunc("/signup", views.SignupHandler)
	mux.HandleFunc("/login", views.LoginHandler)
	mux.HandleFunc("/logout", views.LogoutHandler)
	mux.HandleFunc("/changepass", views.ChangePasswdHandler)
	mux.HandleFunc("/forgotpass", views.ForgotPasswdHandler)
	mux.HandleFunc("/resetpass", views.ResetPasswdHandler)

	mux.HandleFunc("/users", views.UserProfileHandler)
	mux.HandleFunc("/users/update", views.UserProfileUpdateHandler)
	mux.HandleFunc("/users/comments", views.UserCommentsHandler)
	mux.HandleFunc("/users/topics", views.UserTopicsHandler)
	mux.HandleFunc("/users/groups", views.UserGroupsHandler)

	if *fcgiMode {
		fcgi.Serve(nil, mux)
		return
	}

	srv := &http.Server{
		Handler:      mux,
		Addr:         *addr,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Println("[INFO] Starting orangeforum at", *addr)
	err := srv.ListenAndServe()
	if err != nil {
		log.Panicf("[ERROR] %s\n", err)
	}

}
