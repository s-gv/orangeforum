// Copyright (c) 2021 Orange Forum authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package models

import "github.com/jmoiron/sqlx"

var DB *sqlx.DB

var BannedIpsGroupedByDomain map[int][]string

/*
On startup :
	1.Get all banned ip addresses per domain
	2.Populate models DS , map<domainId, [] bannedIps>
	3.for each of element in map
		construct trie of banned ip addresses

On ip filter http handler
	1.Find domain id from the request
	2.Find ip address from the trie
	3.If found send forbidden status
	else continue
*/
