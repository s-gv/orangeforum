#!/bin/bash

# Copyright (c) 2021 Orange Forum authors. All rights reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.

sudo -u postgres psql postgres -c "DROP DATABASE IF EXISTS oftestdb;"
sudo -u postgres psql postgres -c "DROP USER IF EXISTS ofdbuser;"

sudo -u postgres psql postgres -c "CREATE USER ofdbuser WITH PASSWORD 'ofdbpass';"
sudo -u postgres psql postgres -c "ALTER ROLE ofdbuser SET client_encoding TO 'utf8';"
sudo -u postgres psql postgres -c "ALTER ROLE ofdbuser SET default_transaction_isolation TO 'read committed';"
sudo -u postgres psql postgres -c "ALTER ROLE ofdbuser SET timezone TO 'UTC';"

sudo -u postgres psql postgres -c "CREATE DATABASE oftestdb;"
sudo -u postgres psql postgres -c "GRANT ALL PRIVILEGES ON DATABASE oftestdb TO ofdbuser;"

