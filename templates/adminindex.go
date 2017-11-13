// Copyright (c) 2017 Sagar Gubbi. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package templates

const adminindexSrc = `
{{ define "content" }}

<h1>Config</h1>

<form action="/admin" method="POST">
<input type="hidden" name="csrf" value="{{ .Common.CSRF }}">
<table class="form">
	<tr>
		<th><label for="forum_name">Forum Name:</label></th>
		<td><input type="text" name="forum_name" id="forum_name" value="{{ index .Config "forum_name" }}" required></td>
	</tr>
	<tr>
		<th><label for="header_msg">Announcement:</label></th>
		<td><input type="text" name="header_msg" id="header_msg" value="{{ index .Config "header_msg" }}"></td>
	</tr>
	<tr>
		<th><label for="login_msg">Login message:</label></th>
		<td><input type="text" name="login_msg" id="login_msg" value="{{ index .Config "login_msg" }}"></td>
	</tr>
	<tr>
		<th><label for="signup_msg">Signup message:</label></th>
		<td><input type="text" name="signup_msg" id="signup_msg" value="{{ index .Config "signup_msg" }}"></td>
	</tr>
	<tr>
		<th><label for="censored_words"><div class="col-label">Censored words:</label></th>
		<td><textarea name="censored_words" id="censored_words" rows="4" placeholder="shit, bitch, poop">{{ index .Config "censored_words" }}</textarea></td>
	</tr>
	<tr>
		<th><label for="body_appendage"><div class="col-label">Body Appendage:</label></th>
		<td><textarea name="body_appendage" id="body_appendage" rows="4" placeholder="<script>Analytics or something</script>">{{ index .Config "body_appendage" }}</textarea></td>
	</tr>
	<tr>
		<th><label for="data_dir"><div class="col-label">Data Directory:</label></th>
		<td><input type="text" name="data_dir" id="data_dir" value="{{ index .Config "data_dir" }}"></td>
	</tr>
	<tr>
		<th><label for="default_from_mail"><div class="col-label">FROM E-mail:</label></th>
		<td><input type="text" name="default_from_mail" id="default_from_mail" value="{{ index .Config "default_from_mail" }}"></td>
	</tr>
	<tr>
		<th><label for="smtp_host">SMTP Host:</label></th>
		<td><input type="text" name="smtp_host" id="smtp_host" value="{{ index .Config "smtp_host" }}"></td>
	</tr>
	<tr>
		<th><label for="smtp_port">SMTP Port:</label></th>
		<td><input type="number" name="smtp_port" id="smtp_port" value="{{ index .Config "smtp_port" }}"></td>
	</tr>
	<tr>
		<th><label for="smtp_user">SMTP Username:</label></th>
		<td><input type="text" name="smtp_user" id="smtp_user" value="{{ index .Config "smtp_user" }}"></td>
	</tr>
	<tr>
		<th><label for="smtp_pass">SMTP Password:</label></th>
		<td><input type="text" name="smtp_pass" id="smtp_pass" value="{{ index .Config "smtp_pass" }}"></td>
	</tr>
	<tr>
		<th><label for="read_only">Read-only mode:</label></th>
		<td><input type="checkbox" name="read_only" id="read_only" value="1"{{ if index .Config "read_only" }} checked{{ end }}></td>
	</tr>
	<tr>
		<th><label for="signup_disabled">Signup disabled:</label></th>
		<td><input type="checkbox" name="signup_disabled" id="signup_disabled" value="1"{{ if index .Config "signup_disabled" }} checked{{ end }}></td>
	</tr>
	<tr>
		<th><label for="group_creation_disabled">Group creation disabled:</label></th>
		<td><input type="checkbox" name="group_creation_disabled" id="group_creation_disabled" value="1"{{ if index .Config "group_creation_disabled" }} checked{{ end }}></td>
	</tr>
	<tr>
		<th><label for="image_upload_enabled">Allow image upload:</label></th>
		<td><input type="checkbox" name="image_upload_enabled" id="image_upload_enabled" value="1"{{ if index .Config "image_upload_enabled" }} checked{{ end }}></td>
	</tr>
	<tr>
		<th><label for="allow_group_subscription">Allow e-mail subscriptions to groups:</label></th>
		<td><input type="checkbox" name="allow_group_subscription" id="allow_group_subscription" value="1"{{ if index .Config "allow_group_subscription" }} checked{{ end }}></td>
	</tr>
	<tr>
		<th><label for="allow_topic_subscription">Allow e-mail subscriptions to topics:</label></th>
		<td><input type="checkbox" name="allow_topic_subscription" id="allow_topic_subscription" value="1"{{ if index .Config "allow_topic_subscription" }} checked{{ end }}></td>
	</tr>
	{{ if .Common.Msg }}
	<tr>
		<th></th>
		<td><span class="alert">{{ .Common.Msg }}</span></td>
	</tr>
	{{ end }}
	<tr>
		<th></th>
		<td><input type="submit" value="Update"></td>
	</tr>
</table>
</form>

<h1>Footer links</h1>

{{ range .ExtraNotes }}
<form action="/admin" method="POST">
<input type="hidden" name="csrf" value="{{ $.Common.CSRF }}">
<input type="hidden" name="linkid" value="{{ .ID }}">
<table class="form">
	<tr>
		<th>Link name:</th>
		<td><input type="text" name="name" value="{{ .Name }}"></td>
	</tr>
	<tr>
		<th>External URL / Content:</th>
		<td><input type="text" name="url" value="{{ .URL }}"></td>
	</tr>
	<tr>
		<th></th>
		<td><textarea name="content" rows="6">{{ .Content }}</textarea></td>
	</tr>
	<tr>
		<th></th>
		<td>
			<input type="submit" name="submit" value="Update">
			<input type="submit" name="submit" value="Delete">
		</td>
	</tr>
</table>
</form>
{{ end }}

<form action="/admin" method="POST">
<input type="hidden" name="csrf" value="{{ .Common.CSRF }}">
<input type="hidden" name="linkid" value="new">
<table class="form">
	<tr>
		<th>New Link:</th>
		<td><input type="text" name="name" placeholder="Privacy Policy"></td>
	</tr>
	<tr>
		<th>External URL / Content:</th>
		<td><input type="text" name="url" placeholder="https://..."></td>
	</tr>
	<tr>
		<th></th>
		<td><textarea name="content" rows="6" placeholder="Lorem ipsum dolor sit amet, consectetur adipiscing elit. Pellentesque mollis elit hendrerit mattis vulputate. Vivamus congue convallis urna. Ut mollis ligula velit, vitae feugiat nulla laoreet sit amet. Mauris hendrerit arcu ut quam vestibulum tincidunt ac id erat. Integer vehicula congue orci a sagittis. Cras arcu nibh, scelerisque et ultricies luctus, mollis sit amet ex. Donec a volutpat nibh, ac venenatis nisi. Aliquam erat volutpat. Vivamus ut ex rutrum, tristique neque et, pulvinar velit. Duis facilisis tincidunt arcu nec imperdiet. Cras semper metus nec quam ornare, eget tempor magna tincidunt. Integer fringilla nisl ligula, vel iaculis orci commodo vitae. Duis a volutpat mauris. Aenean ac felis at metus tristique hendrerit cursus a velit. Curabitur a est tellus. Morbi cursus nisi porta nisi congue iaculis."></textarea></td>
	</tr>
	<tr>
		<th></th>
		<td><input type="submit" value="Create new link"></td>
	</tr>
</table>
</form>

<h1>Stats</h1>
<table>
	<tr>
		<th>Number of users:</th>
		<td>{{ .NumUsers }}</td>
	</tr>
	<tr>
		<th>Number of groups:</th>
		<td>{{ .NumGroups }}</td>
	</tr>
	<tr>
		<th>Number of topics:</th>
		<td>{{ .NumTopics }}</td>
	</tr>
	<tr>
		<th>Number of comments:</th>
		<td>{{ .NumComments }}</td>
	</tr>
</table>

{{ end }}`
