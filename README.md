# Record

[![Build Status](https://travis-ci.org/ksoichiro/record.svg?branch=master)](https://travis-ci.org/ksoichiro/record)

TODO

## MySQL

Local:

```sh
docker-compose up -d
docker-compose exec db bash -c "mysql -uroot -p test"
# or
mysql -h 127.0.0.1 -uroot -p test
```

## API

```sh
cd api
go run main.go
```

## RSA Keys

```sh
$ ssh-keygen -t rsa -b 4096 -m PEM -f jwtRS256.key -C "user@localhost"
Generating public/private rsa key pair.
Enter passphrase (empty for no passphrase):
Enter same passphrase again:
:
$ ssh-keygen -f jwtRS256.key.pub -e -m pkcs8 > jwtRS256.key.pub.pkcs8
$ chmod 0600 jwtRS256.key*
```
