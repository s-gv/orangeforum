package templates

const adminindexSrc = `
{{ define "head" }}
<style>
.row {
	margin-top: 15px;
}
input[type="text"], input[type="number"], textarea {
	width: 90%;
}
.ccol1 {
	float: left;
	max-width: 80%;
}
.ccol2 {
	float: right;
	margin-left: 15px;
}
@media screen and (min-width:600px) {
	.tcol1 {
		float: left;
		text-align: right;
		width: 300px;
	}
	.tcol2 {
		float: left;
		text-align: left;
		margin-left: 15px;
		width: 300px;
	}
	.ccol1 {
		text-align: right;
		width: 300px;
	}
	.ccol2 {
		float: left;
	}

}
</style>
{{ end }}

{{ define "content" }}

<h1>Config</h1>

<form action="/admin" method="POST">
<input type="hidden" name="csrf" value="{{ .Common.CSRF }}">

<div class="row clearfix">
	<div class="tcol1">Forum Name</div>
	<div class="tcol2"><input type="text" name="forum_name" value="{{ index .Config "forum_name" }}" required></div>
</div>
<div class="row clearfix">
	<div class="tcol1">Announcement</div>
	<div class="tcol2"><input type="text" name="header_msg" value="{{ index .Config "header_msg" }}"></div>
</div>
<div class="row clearfix">
	<div class="tcol1">Data Directory</div>
	<div class="tcol2"><input type="text" name="data_dir" value="{{ index .Config "data_dir" }}"></div>
</div>
<div class="row clearfix">
	<div class="tcol1">FROM E-mail</div>
	<div class="tcol2"><input type="text" name="default_from_mail" value="{{ index .Config "default_from_mail" }}"></div>
</div>
<div class="row clearfix">
	<div class="tcol1">SMTP Host</div>
	<div class="tcol2"><input type="text" name="smtp_host" value="{{ index .Config "smtp_host" }}"></div>
</div>
<div class="row clearfix">
	<div class="tcol1">SMTP Port</div>
	<div class="tcol2"><input type="number" name="smtp_port" value="{{ index .Config "smtp_port" }}"></div>
</div>
<div class="row clearfix">
	<div class="tcol1">SMTP Username</div>
	<div class="tcol2"><input type="text" name="smtp_user" value="{{ index .Config "smtp_user" }}"></div>
</div>
<div class="row clearfix">
	<div class="tcol1">SMTP Password</div>
	<div class="tcol2"><input type="text" name="smtp_pass" value="{{ index .Config "smtp_pass" }}"></div>
</div>
<div class="row clearfix">
	<div class="ccol1">Signup disabled</div>
	<div class="ccol2"><input type="checkbox" name="signup_disabled" value="1"{{ if index .Config "signup_disabled" }} checked{{ end }}></div>
</div>
<div class="row clearfix">
	<div class="ccol1">Group creation disabled</div>
	<div class="ccol2"><input type="checkbox" name="group_creation_disabled" value="1"{{ if index .Config "group_creation_disabled" }} checked{{ end }}></div>
</div>
<div class="row clearfix">
	<div class="ccol1">Allow image upload</div>
	<div class="ccol2"><input type="checkbox" name="image_upload_enabled" value="1"{{ if index .Config "image_upload_enabled" }} checked{{ end }}></div>
</div>
<div class="row clearfix">
	<div class="ccol1">Allow file upload</div>
	<div class="ccol2"><input type="checkbox" name="file_upload_enabled" value="1"{{ if index .Config "file_upload_enabled" }} checked{{ end }}></div>
</div>
<div class="row clearfix">
	<div class="ccol1">Allow e-mail subscriptions to groups</div>
	<div class="ccol2"><input type="checkbox" name="allow_group_subscription" value="1"{{ if index .Config "allow_group_subscription" }} checked{{ end }}></div>
</div>
<div class="row clearfix">
	<div class="ccol1">Allow e-mail subscriptions to topics</div>
	<div class="ccol2"><input type="checkbox" name="allow_topic_subscription" value="1"{{ if index .Config "allow_topic_subscription" }} checked{{ end }}></div>
</div>
<div class="row clearfix">
	<div class="ccol1">Save changes?</div>
	<div class="ccol2"><input type="submit" value="Update"></div>
</div>

</form>

<h1>Stats</h1>

<div class="row clearfix">
	<div class="ccol1">Number of users</div>
	<div class="ccol2">{{ .NumUsers }}</div>
</div>
<div class="row clearfix">
	<div class="ccol1">Number of groups</div>
	<div class="ccol2">{{ .NumGroups }}</div>
</div>
<div class="row clearfix">
	<div class="ccol1">Number of topics</div>
	<div class="ccol2">{{ .NumTopics }}</div>
</div>
<div class="row clearfix">
	<div class="ccol1">Number of comments</div>
	<div class="ccol2">{{ .NumComments }}</div>
</div>

<h1>Footer links</h1>

{{ range .ExtraNotes }}
<form action="/admin" method="POST">
<input type="hidden" name="csrf" value="{{ $.Common.CSRF }}">
<input type="hidden" name="linkid" value="{{ .ID }}">
<div class="row clearfix">
	<div class="tcol1">Link name</div>
	<div class="tcol2"><input type="text" name="name" value="{{ .Name }}"></div>
</div>
<div class="clearfix">
	<div class="tcol1">External URL / Content</div>
	<div class="tcol2">
		<input type="text" name="url" value="{{ .URL }}">
		<textarea name="content" rows="6">{{ .Content }}</textarea>
		<input type="submit" name="submit" value="Update">
		<input type="submit" name="submit" value="Delete">
	</div>
</div>
</form>
{{ end }}

<form action="/admin" method="POST">
<input type="hidden" name="csrf" value="{{ .Common.CSRF }}">
<input type="hidden" name="linkid" value="new">
<div class="row clearfix">
	<div class="tcol1">New Link</div>
	<div class="tcol2"><input type="text" name="name" placeholder="Privacy Policy"></div>
</div>
<div class="clearfix">
	<div class="tcol1">External URL / Content</div>
	<div class="tcol2">
		<input type="text" name="url" placeholder="https://...">
		<textarea name="content" rows="6" placeholder="Lorem ipsum dolor sit amet, consectetur adipiscing elit. Pellentesque mollis elit hendrerit mattis vulputate. Vivamus congue convallis urna. Ut mollis ligula velit, vitae feugiat nulla laoreet sit amet. Mauris hendrerit arcu ut quam vestibulum tincidunt ac id erat. Integer vehicula congue orci a sagittis. Cras arcu nibh, scelerisque et ultricies luctus, mollis sit amet ex. Donec a volutpat nibh, ac venenatis nisi. Aliquam erat volutpat. Vivamus ut ex rutrum, tristique neque et, pulvinar velit. Duis facilisis tincidunt arcu nec imperdiet. Cras semper metus nec quam ornare, eget tempor magna tincidunt. Integer fringilla nisl ligula, vel iaculis orci commodo vitae. Duis a volutpat mauris. Aenean ac felis at metus tristique hendrerit cursus a velit. Curabitur a est tellus. Morbi cursus nisi porta nisi congue iaculis."></textarea>
		<input type="submit" value="Create new link">
	</div>
</div>
</form>

{{ .Msg }}

{{ end }}`