## Requirements
* GoLang 1.13+
* PostgreSQL 13+(with pg_trgm extension)
* Make
* Docker
* docker-compose


# The structure of project
* assets - store migrations, dumps and etc
* configs - store config files
* deployments - store docker, docker-compose files
* di - contain different implementation DI libraries


# The structure of the DI folders
* cmd - store different commands
* container - DI helper folder
* definition - main implementation common DI containers
* internal - code of project


# DI folder
Folder contain different implementation DI libraries. I will add another implementation in future.
* sarulabsdi - implementation with use `github.com/sarulabs/di`
* sarulabsdingo - implementation with use `github.com/sarulabs/dingo`


## Examples
# Run http server
`./bin/theboatscom http-server -c /go/src/github.com/tumani1/diexample/configs/config.json`

`-c` - path to config


# Run migrations
Roll up migration

`./bin/theboatscom migrate up -c /go/src/github.com/tumani1/diexample/configs/config.json -m /go/src/github.com/tumani1/diexample/asstets/migrations`

Roll back migration

`./bin/theboatscom migrate down -c /go/src/github.com/tumani1/diexample/configs/config.json -m /go/src/github.com/tumani1/diexample/asstets/migrations`

`-c` - path to config
`-m` - path to migration folder


#Example requests to API
Curl request for look up data:

`curl -iv -X GET "http://127.0.0.1:9090/api/search?query=test"`

Curl request for autocomplete:

`curl -iv -X GET "http://127.0.0.1:9090/api/autocomplete?query=test"`
