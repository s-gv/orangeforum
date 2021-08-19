Orange Forum
============

[Orange Forum](http://www.goodoldweb.com/orangeforum/) is an easy to deploy forum that has minimal dependencies (only Postgres) and uses almost no javascript. It is written is golang and a [compiled binary](https://github.com/s-gv/orangeforum/releases) is available for linux. Try the latest version hosted [here](https://community.goodoldweb.com/). Please contact [info@goodoldweb.com](mailto:info@goodoldweb.com) if you have any questions or want support.

How to use
----------

Install PostgreSQL and create a database and user using the `psql` command:

```
CREATE DATABASE testdb;
CREATE USER dbuser WITH PASSWORD 'dbpass';
ALTER ROLE dbuser SET client_encoding TO 'utf8';
ALTER ROLE dbuser SET default_transaction_isolation TO 'read committed';
ALTER ROLE dbuser SET timezone TO 'UTC';
GRANT ALL PRIVILEGES ON DATABASE testdb TO dbuser;
```

[Download](https://github.com/s-gv/orangeforum/releases) the Orange Forum binary and migrate the database with:

```
./orangeforum -migrate
```

Create a domain (example: community.goodoldweb.com):

```
./orangeforum -createdomain
```

Create an admin:

```
./orangeforum -createsuperuser
```

Enter the SMTP server details:

```
./orangeforum -setsmtp
```

Finally, start the server (don't forget to change the secret key!):

```
SECRET_KEY=s6JM1e8JTAphtKNR2y27XA8kkAaXOSYB ORANGEFORUM_DSN="postgres://dbuser:dbpass@localhost:5432/testdb" ./orangeforum -alsologtostderr
```

Goto [http://localhost:9123/forums/community.goodoldweb.com](http://localhost:9123/forums/community.goodoldweb.com) in your browser.

Notes
-----

There are two types of privileged users in orangeforum: `admin`, and `mod`. Both can edit/delete/close posts and ban users. In addition to these privileges, `admin`s control which users are designated as `mod`s.

Dependencies
------------

- Go 1.16 (only for compiling)
- Postgres 9.5 or newer