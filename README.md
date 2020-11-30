# Go-REST
[![CircleCI](https://circleci.com/gh/rnov/Go-REST/tree/master.svg?style=svg)](https://circleci.com/gh/rnov/Go-REST/tree/master)
[![Coverage Status](https://coveralls.io/repos/github/rnov/Go-REST/badge.svg?branch=master)](https://coveralls.io/github/rnov/Go-REST?branch=master)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
### Description:

Gorest is a portfolio like REST API project that also serves as a reminder and template since most of the projects i'm involved are private.
Note that the business logic is rather simplistic and is not the focus of this project.

API REST endpoints:

| Name   | Method       | URL                   | Protected |
| :---:    | :---:      | :---:                 | :---:       |
| List   | `GET`        | `/recipes`             | ✘         |
| Create | `POST`       | `/recipes`             | ✓         |
| Get    | `GET`        | `/recipes/{ID}`        | ✘         |
| Update | `PUT/PATCH`  | `/recipes/{ID}`        | ✓         |
| Delete | `DELETE`     | `/recipes/{ID}`        | ✓         |
| Rate   | `POST`       | `/recipes/{ID}/rate`   | ✘         |


I tried to keep the code as vanilla as possible - avoiding third party packages some of them are :

| Packages |Description
|:--------:|:-----------------------------------:|
|github.com/gorilla/mux| Most popular router handler for Go, great performance easy to use |
|github.com/go-redis/redis| A redis client package for Go |
|github.com/go-logging| Great package to manage logs, very useful in larger projects |
|gopkg.in/yaml.v2| Package used to proceed the config files |


### Deployment :

1- Docker :
```sh
# only dockefile, mind to start a DB to connect to
docker built -t gorest:test .
docker run gorest:test

# with compose :
$ docker-compose build
$ docker-compose run
```

2- Local:

```sh
$ export ENV_PATH="config/envs/local/config.yml"
```

Once running, in order to make protected call redis db needs to be populated, run following command :
```sh
$ cat populate-Redis.sh
```

### Architecture :

* Hexagonal(Onion) like design, build in mind to separate the different adapters, from the application and domain
  layers, achieved by relaying on dependency injection through interfaces.
* Logs only in handlers (adapter).
* Defined error types, since all the errors when produced change the program flow but each error type (user, server...)
    have different behaviour and consequences.
* Twelve-Factor App principles in mind.
* No third party lib for testing, thanks to the design everything can be mocked easily.
