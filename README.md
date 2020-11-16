# Go-REST
[![Coverage Status](https://coveralls.io/repos/github/rnov/Go-REST/badge.svg?branch=master)](https://coveralls.io/github/rnov/Go-REST?branch=fix/code-refactor)
[![Go Report Card](https://goreportcard.com/badge/github.com/rnov/Go-REST)](https://goreportcard.com/report/github.com/rnov/Go-REST)
### Description:

An api REST in go - a toy project - implementing SOLID principles and applying best practices whenever possible,
it aims to be used as a reference and a portfolio since everything i do is private.

This simple Api REST is used to CRUD restaurant's recipes, as today has only 6 calls :


| Name   | Method       | URL                   | Protected |
| :---:    | :---:      | :---:                 | :---:       |
| List   | `GET`        | `/recipes`             | ✘         |
| Create | `POST`       | `/recipes`             | ✓         |
| Get    | `GET`        | `/recipes/{ID}`        | ✘         |
| Update | `PUT/PATCH`  | `/recipes/{ID}`        | ✓         |
| Delete | `DELETE`     | `/recipes/{ID}`        | ✓         |
| Rate   | `POST`       | `/recipes/{ID}/rate`   | ✘         |


I tried to keep the code as vanilla as possible - avoiding installing many third party packages, some of them are:

| Packages |Description
|:--------:|:-----------------------------------:|
|github.com/julienschmidt/httprouter| Most popular router handler for Go, great performance easy to use |
|github.com/go-redis/redis| A redis client package for Go |
|github.com/go-logging| Great package to manage logs, very useful in larger projects |
|gopkg.in/yaml.v2| Package used to proceed the config files |


### Deployment :

1- Easy way, using Docker :
```sh
$ docker-compose build
$ docker-compose run
```

2- The long way is configuring the environment in the host, need to install some packages and set a env variable to read the configuration file.

```sh
$ export ENV_PATH="config/environments/production/config.yml"
```

Once running, in order to make protected call redis db needs to be populated, run following command :
```sh
$ cat populateRedis.txt | redis-cli -p 6379
```
Note: There is a json file with the calls that can be imported by Postman.
### Architecture :

It is designed to separate low level from high level abstraction, making it easier to add new functionalities and scaling.
It would be quite easy to add or migrate to another DB, lets say psql/mongo or adding multilevel db architecture without changing any high level abstraction code and vice versa.