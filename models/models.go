// Copyright (c) 2021 Orange Forum authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package models

import "github.com/jmoiron/sqlx"

var DB *sqlx.DB

var BannedIpsGroupedByDomain map[int][]string

type ipv4AddressTrieNode struct {
	addresOctet        byte
	children           map[byte]*ipv4AddressTrieNode
	octectIndex        int
	isLastAddressOctet bool
}

type ipv4AddressTrie struct {
	root *ipv4AddressTrieNode
}

var BannedIpv4AddressTriesPerDomain map[int]*ipv4AddressTrie
