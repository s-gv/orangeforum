// Copyright (c) 2021 Orange Forum authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"syscall"

	"github.com/s-gv/orangeforum/models"
	"golang.org/x/crypto/ssh/terminal"
)

func getCreds() (int, string, string, string) {
	var domainName string
	fmt.Printf("Domain name: ")
	fmt.Scanf("%s\n", &domainName)

	domainID := models.GetDomainIDByName(domainName)
	if domainID == nil {
		fmt.Printf("[ERROR] Invalid domain\n")
		return -1, "", "", ""
	}

	var email string
	fmt.Printf("E-mail: ")
	fmt.Scanf("%s\n", &email)

	if strings.Contains(email, " ") || !strings.Contains(email, "@") {
		fmt.Printf("[ERROR] Invalid email\n")
		return -1, "", "", ""
	}

	userName := strings.Split(email, "@")[0] + strconv.Itoa(rand.Intn(100000000))

	fmt.Printf("Password: ")
	password, err := terminal.ReadPassword(int(syscall.Stdin))
	fmt.Printf("\n")
	if err != nil {
		fmt.Printf("[ERROR] %s\n", err)
		return -1, "", "", ""
	}
	if len(password) < 8 {
		fmt.Printf("[ERROR] Password should have at least 8 characters.\n")
		return -1, "", "", ""
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
		return -1, "", "", ""
	}

	return *domainID, email, userName, pass
}

func commandCreateSuperUser() {
	fmt.Printf("Creating superuser...\n")
	domainID, email, userName, passwd := getCreds()
	if userName != "" && passwd != "" {
		if err := models.CreateSuperUser(domainID, email, userName, passwd); err != nil {
			fmt.Printf("Error creating superuser: %s\n", err)
		}
	}
}

func commandCreateUser() {
	fmt.Printf("Creating user...\n")
	domainID, email, userName, passwd := getCreds()
	if userName != "" && passwd != "" {
		if err := models.CreateUser(domainID, email, userName, passwd); err != nil {
			fmt.Printf("Error creating superuser: %s\n", err)
		}
	}
}

func commandChangePasswd() {
	fmt.Printf("Enter new credentials...\n")
	domainID, email, _, passwd := getCreds()
	if passwd != "" {
		if err := models.ChangePasswd(domainID, email, passwd); err != nil {
			fmt.Printf("Error changing password: %s\n", err)
		}
	}
}

func commandCreateDomain() {
	fmt.Printf("Creating domain...\n")

	var domainName string
	fmt.Printf("Domain name: ")
	fmt.Scanf("%s\n", &domainName)

	err := models.CreateDomain(domainName)
	if err != nil {
		fmt.Printf("Error creating domain: %s\n", err)
	}
}
