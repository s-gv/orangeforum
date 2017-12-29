// Copyright (c) 2017 Sagar Gubbi. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package static

const ScriptSrc = `
var inputs = document.getElementsByClassName("no-double-post"); //document.getElementsByTagName('input');
for (var i = 0; i < inputs.length; i++) {
	if (inputs[i].type == "submit") {
		var btn = inputs[i];
		var form = btn.form;
		if (!!form && form.method == "post") {
			btn.onclick = function() {
				var form = this.form;
				if(!!form && form.checkValidity()) {
					var submitBtn = this;
				 	setTimeout(function() {
				 		submitBtn.disabled = true;
				 	}, 1);
				}
			};
		}
	}
}
`
