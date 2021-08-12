// Copyright (c) 2021 Orange Forum authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package templates

import (
	_ "embed"
	"html/template"
)

var Signin, Signup, ChangePass, OneTimeSignin *template.Template // Auth related templates
var Index, Admin *template.Template                              // Other templates

//go:embed base.html
var baseTmplStr string

//go:embed signin.html
var signinTmplStr string

//go:embed signup.html
var signupTmplStr string

//go:embed ot_signin.html
var otSigninTmplStr string

//go:embed changepass.html
var changePassTmplStr string

//go:embed index.html
var indexTmplStr string

//go:embed admin.html
var adminTmplStr string

func init() {
	Signin = template.Must(template.Must(template.New("base").Parse(baseTmplStr)).Parse(signinTmplStr))
	Signup = template.Must(template.Must(template.New("base").Parse(baseTmplStr)).Parse(signupTmplStr))
	ChangePass = template.Must(template.Must(template.New("base").Parse(baseTmplStr)).Parse(changePassTmplStr))
	OneTimeSignin = template.Must(template.Must(template.New("base").Parse(baseTmplStr)).Parse(otSigninTmplStr))

	Index = template.Must(template.Must(template.New("base").Parse(baseTmplStr)).Parse(indexTmplStr))
	Admin = template.Must(template.Must(template.New("base").Parse(baseTmplStr)).Parse(adminTmplStr))
}
