// Copyright (c) 2021 Orange Forum authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package models

func Migrate() error {
	migrateSqlite0001(DB)
	return nil
}