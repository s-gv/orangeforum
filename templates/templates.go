// Copyright (c) 2017 Sagar Gubbi. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package templates

import (
	"html/template"
	"io"
	"log"
)

var tmpls map[string]*template.Template = make(map[string]*template.Template)

func init() {
	tmpls["adminindex.html"] = template.Must(template.New("base").Parse(baseSrc))
	template.Must(tmpls["adminindex.html"].New("adminindex").Parse(adminindexSrc))

	tmpls["changepass.html"] = template.Must(template.New("base").Parse(baseSrc))
	template.Must(tmpls["changepass.html"].New("changepass").Parse(changepassSrc))

	tmpls["commentedit.html"] = template.Must(template.New("base").Parse(baseSrc))
	template.Must(tmpls["commentedit.html"].New("commentedit").Parse(commenteditSrc))

	tmpls["commentindex.html"] = template.Must(template.New("base").Parse(baseSrc))
	template.Must(tmpls["commentindex.html"].New("commentindex").Parse(commentindexSrc))

	tmpls["extranote.html"] = template.Must(template.New("base").Parse(baseSrc))
	template.Must(tmpls["extranote.html"].New("extranote").Parse(extranoteSrc))

	tmpls["forgotpass.html"] = template.Must(template.New("base").Parse(baseSrc))
	template.Must(tmpls["forgotpass.html"].New("forgotpass").Parse(forgotpassSrc))

	tmpls["groupindex.html"] = template.Must(template.New("base").Parse(baseSrc))
	template.Must(tmpls["groupindex.html"].New("groupindex").Parse(groupindexSrc))

	tmpls["groupedit.html"] = template.Must(template.New("base").Parse(baseSrc))
	template.Must(tmpls["groupedit.html"].New("groupedit").Parse(groupeditSrc))

	tmpls["groups.html"] = template.Must(template.New("base").Parse(baseSrc))
	template.Must(tmpls["groups.html"].New("groups").Parse(groupindexSrc))

	tmpls["index.html"] = template.Must(template.New("base").Parse(baseSrc))
	template.Must(tmpls["index.html"].New("index").Parse(indexSrc))

	tmpls["login.html"] = template.Must(template.New("base").Parse(baseSrc))
	template.Must(tmpls["login.html"].New("login").Parse(loginSrc))

	tmpls["profile.html"] = template.Must(template.New("base").Parse(baseSrc))
	template.Must(tmpls["profile.html"].New("profile").Parse(profileSrc))

	tmpls["profilecomments.html"] = template.Must(template.New("base").Parse(baseSrc))
	template.Must(tmpls["profilecomments.html"].New("profilecomments").Parse(profilecommentsSrc))

	tmpls["profiletopics.html"] = template.Must(template.New("base").Parse(baseSrc))
	template.Must(tmpls["profiletopics.html"].New("profiletopics").Parse(profiletopicsSrc))

	tmpls["profilegroups.html"] = template.Must(template.New("base").Parse(baseSrc))
	template.Must(tmpls["profilegroups.html"].New("profilegroups").Parse(profilegroupsSrc))

	tmpls["resetpass.html"] = template.Must(template.New("base").Parse(baseSrc))
	template.Must(tmpls["resetpass.html"].New("resetpass").Parse(resetpassSrc))

	tmpls["signup.html"] = template.Must(template.New("base").Parse(baseSrc))
	template.Must(tmpls["signup.html"].New("signup").Parse(signupSrc))

	tmpls["topicedit.html"] = template.Must(template.New("base").Parse(baseSrc))
	template.Must(tmpls["topicedit.html"].New("topicedit").Parse(topiceditSrc))

	tmpls["topicindex.html"] = template.Must(template.New("base").Parse(baseSrc))
	template.Must(tmpls["topicindex.html"].New("topicindex").Parse(topicindexSrc))

	tmpls["pm.html"] = template.Must(template.New("base").Parse(baseSrc))
	template.Must(tmpls["pm.html"].New("pm").Parse(pmSrc))
}

func Render(wr io.Writer, template string, data interface{}) {
	err := tmpls[template].Execute(wr, data)
	if err != nil {
		log.Panicf("[ERROR] Error rendering %s: %s\n", template, err)
	}
}
