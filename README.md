Orange Forum
============

Orange Forum is an easy to deploy forum that has few dependencies. It is written is golang and a (compiled binary)[https://github.com/s-gv/orangeforum/releases] is available for Linux. Try the latest version hosted [here](https://groups.goodoldweb.com/). By default, sqlite is used, so it's easy to get started. Download the binary and create an admin user with:

```
./orangeforum --migrate && ./orangeforum --createsuperuser
```

Now, start the server:

```
./orangeforum
```

Dependencies
------------

- Go 1.8 (only for compiling)
- Postgres 9.5 (or use embedded sqlite3)
