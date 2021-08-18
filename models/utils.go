// Copyright (c) 2021 Orange Forum authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package models

import (
	"strconv"
	"time"
)

func RelTimeNowStr(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)
	if diff.Hours() > 24*30 {
		return t.Format("2 Jan 2006")
	} else if diff.Hours() > 24 {
		days := int(diff.Hours() / 24)
		if days == 1 {
			return strconv.Itoa(days) + " day ago"
		}
		return strconv.Itoa(days) + " days ago"
	} else if diff.Minutes() > 120 {
		return strconv.Itoa(int(diff.Hours())) + " hours ago"
	} else {
		mins := int(diff.Minutes())
		if mins == 1 {
			return strconv.Itoa(mins) + " minute ago"
		}
		return strconv.Itoa(mins) + " minutes ago"
	}
}

func ApproxNumStr(i int) string {
	if i >= (1000 * 1000) {
		return strconv.Itoa(i/(1000*1000)) + "M"
	} else if i >= 1000 {
		return strconv.Itoa(i/1000) + "k"
	} else {
		return strconv.Itoa(i)
	}
}
