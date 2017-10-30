Orange Forum
============

Orange Forum is an easy to deploy forum that has minimal dependencies and uses very little javascript. It is written is golang and a [compiled binary](https://github.com/s-gv/orangeforum/releases) is available for linux. Try the latest version hosted [here](https://groups.goodoldweb.com/). If you have any questions, send an e-mail to: [info@goodoldweb.com](mailto:info@goodoldweb.com)


How to use
----------

By default, sqlite is used, so it's easy to get started. [Download](https://github.com/s-gv/orangeforum/releases) the binary, and run the migrate command to setup the database tables:

```
./orangeforum --migrate
```

Create an admin user with:

```
./orangeforum --createsuperuser
```

Now, start the server:

```
./orangeforum
```

Dependencies
------------

- Go 1.8 (only for compiling)
- Postgres 9.5 (or use embedded sqlite3)
