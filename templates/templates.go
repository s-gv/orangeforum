package templates

import (
	"html/template"
	"io"
)

var tmpls map[string]*template.Template = make(map[string]*template.Template)

func init() {
	tmpls["changepass.html"] = template.Must(template.New("base").Parse(baseSrc))
	template.Must(tmpls["changepass.html"].New("changepass").Parse(changepassSrc))

	tmpls["creategroup.html"] = template.Must(template.New("base").Parse(baseSrc))
	template.Must(tmpls["creategroup.html"].New("creategroup").Parse(creategroupSrc))

	tmpls["forgotpass.html"] = template.Must(template.New("base").Parse(baseSrc))
	template.Must(tmpls["forgotpass.html"].New("forgotpass").Parse(forgotpassSrc))

	tmpls["groupindex.html"] = template.Must(template.New("base").Parse(baseSrc))
	template.Must(tmpls["groupindex.html"].New("groupindex").Parse(groupindexSrc))

	tmpls["groupnew.html"] = template.Must(template.New("base").Parse(baseSrc))
	template.Must(tmpls["groupnew.html"].New("groupnew").Parse(groupnewSrc))

	tmpls["groups.html"] = template.Must(template.New("base").Parse(baseSrc))
	template.Must(tmpls["groups.html"].New("groups").Parse(groupindexSrc))

	tmpls["index.html"] = template.Must(template.New("base").Parse(baseSrc))
	template.Must(tmpls["index.html"].New("index").Parse(indexSrc))

	tmpls["login.html"] = template.Must(template.New("base").Parse(baseSrc))
	template.Must(tmpls["login.html"].New("login").Parse(loginSrc))

	tmpls["mycomments.html"] = template.Must(template.New("base").Parse(baseSrc))
	template.Must(tmpls["mycomments.html"].New("mycomments").Parse(mycommentsSrc))

	tmpls["mygroups.html"] = template.Must(template.New("base").Parse(baseSrc))
	template.Must(tmpls["mygroups.html"].New("mygroups").Parse(mygroupsSrc))

	tmpls["mytopics.html"] = template.Must(template.New("base").Parse(baseSrc))
	template.Must(tmpls["mytopics.html"].New("mytopics").Parse(mytopicsSrc))

	tmpls["myupvotes.html"] = template.Must(template.New("base").Parse(baseSrc))
	template.Must(tmpls["myupvotes.html"].New("myupvotes").Parse(myupvotesSrc))

	tmpls["new.html"] = template.Must(template.New("base").Parse(baseSrc))
	template.Must(tmpls["new.html"].New("new").Parse(newSrc))

	tmpls["profile.html"] = template.Must(template.New("base").Parse(baseSrc))
	template.Must(tmpls["profile.html"].New("profile").Parse(profileSrc))

	tmpls["reply.html"] = template.Must(template.New("base").Parse(baseSrc))
	template.Must(tmpls["reply.html"].New("reply").Parse(replySrc))

	tmpls["resetpass.html"] = template.Must(template.New("base").Parse(baseSrc))
	template.Must(tmpls["resetpass.html"].New("resetpass").Parse(resetpassSrc))

	tmpls["signup.html"] = template.Must(template.New("base").Parse(baseSrc))
	template.Must(tmpls["signup.html"].New("signup").Parse(signupSrc))

	tmpls["submit.html"] = template.Must(template.New("base").Parse(baseSrc))
	template.Must(tmpls["submit.html"].New("submut").Parse(submitSrc))
}

func Render(wr io.Writer, template string, data interface{}) error {
	return tmpls[template].Execute(wr, data)
}