// Copyright (c) 2021 Orange Forum authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package main

import (
	"fmt"
	"syscall"

	"github.com/s-gv/orangeforum/models"
	"golang.org/x/crypto/ssh/terminal"
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

func commandCreateSuperUser() {
	fmt.Printf("Creating superuser...\n")
	userName, pass := getCreds()
	if userName != "" && pass != "" {
		if err := models.CreateSuperUser(userName, pass); err != nil {
			fmt.Printf("Error creating superuser: %s\n", err)
		}
	}
}

func commandCreateUser() {
}

func commandChangePasswd() {

}
