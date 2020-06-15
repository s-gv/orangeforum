Orange Forum
============

Edit: This branch (V1.x.x) is no longer maintained. Please see the master branch for the latest version.

[Orange Forum](http://www.goodoldweb.com/orangeforum/) is an easy to deploy forum that has minimal dependencies and uses very little javascript.
It is written is golang and a [compiled binary](https://github.com/s-gv/orangeforum/releases) is available for linux.
Try the latest version hosted [here](https://groups.goodoldweb.com/).
Please contact [info@goodoldweb.com](mailto:info@goodoldweb.com) if you have any questions or want support.

How to use
----------

By default, sqlite is used, so it's easy to get started.
[Download](https://github.com/s-gv/orangeforum/releases) the binary and migrate the database with:

```
./orangeforum -migrate
```

Create a superadmin:

```
./orangeforum -createsuperuser
```

Finally, start the server:

```
./orangeforum
```

Notes
-----

There are three types of privileged users in orangeforum: `superadmin`, `admin`, and `mod`. A `superadmin` has rights to
everything across all groups. This includes editing/deleting/closing posts across all groups and editing the list
of mods/admins for all groups in the forum. `mod`s of a group can edit/delete/close posts in that group. `admin`s of
a group have all the privileges of mods for that group and can also edit the name/description of the group and assign
mods for that group.

Orangeforum allows all users to create groups. The user that creates a group becomes an admin of that group.
This can be disabled and group creation can be restricted to the superadmin.

Dependencies
------------

- Go 1.8 (only for compiling)
- Postgres 9.5 (or use embedded sqlite3)

Options
-------

- `-addr <port>`: Use `./orangeforum -addr :8086` to listen on port 8086.
- `-dbdriver <db>` and `-dsn <data_source_name>`: PostgreSQL and SQLite are supported. SQLite is the default driver.

To use postgres, run `./orangeforum -dbdriver postgres -dsn postgres://pguser:pgpasswd@localhost/orangeforum`

To save an sqlite db at a different location, run `./orangeforum -dsn path/to/myforum.db`.

- `-usei2p=<bool>`: Use `./orangeforum -usei2p=true` to forward the service to i2p.
- `-i2pini file`: Use `./orangeforum -i2pini contrib/tunnels.orangeforum.conf` to configure an i2p service with an ini-like file.

When using i2p, the listening port will be set by the i2p configuration, and
arguments passed to -addr will be canceled out.

### Docker

To build the an image for running orangeforum within docker from source, first
clone this git repository and change to it's directory. Then run:

        docker build -t $(whoami)/orangeforum .

To use within docker, you can run by creating a data container interactively to
set up the superuser, then automatically close:

        docker run -i -t --name orangeforum-volume \
            --volumes orangeforum:/opt/orangeforum \
            $(whoami)/orangeforum orangeforum -createsuperuser

Then, run orangeforum as a docker daemon, using the volume from the other docker
container:

        docker run -i -t -d -e args="" \
            --name orangeforum \
            --volumes-from orangeforum-volume \
            -p 127.0.0.0:9123:9123 \
            $(whoami)/orangeforum

If you want, you can check if it's running, for instance

        docker logs orangeforum

Commands
--------

- `-help`: Show a list of all commands and options.
- `-migrate`: Migrate the database. Run this once after updating the orangeforum binary (or when starting afresh).
- `-createsuperuser`: Create a super admin.
- `-createuser`: Create a new user with no special privileges.
- `-changepasswd`: Change password of a user.
- `-deletesessions`: Drop all sessions and log out all users.

optionally, you can pass commands to the docker container by setting the
environment variable args when running the container, for instance

        docker run -i -t -d \
            -e args="-deletesessions" \
            --volumes-from orangeforum-volume \
            --name orangeforum \
            $(whoami)/orangeforum
