#!/bin/bash

sudo -u postgres psql postgres -c "DROP DATABASE IF EXISTS testdb;"
sudo -u postgres psql postgres -c "CREATE DATABASE testdb;"
sudo -u postgres psql postgres -c "GRANT ALL PRIVILEGES ON DATABASE testdb TO dbuser;"

