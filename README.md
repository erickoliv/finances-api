[![Build Status](https://travis-ci.com/erickoliv/finances-api.svg?branch=master)](https://travis-ci.com/erickoliv/finances-api)
[![codecov](https://codecov.io/gh/erickoliv/finances-api/branch/master/graph/badge.svg)](https://codecov.io/gh/erickoliv/finances-api)
[![Go Report Card](https://goreportcard.com/badge/github.com/erickoliv/finances-api)](https://goreportcard.com/report/github.com/erickoliv/finances-api)
[![HitCount](http://hits.dwyl.io/erickoliv/erickoliv/finances-api.svg)](http://hits.dwyl.io/erickoliv/erickoliv/finances-api)
> This is just a personal project I'm using to learn the Go programming language and its libraries.

**using different package approaches to validate the better structure for each use case**

## The road so far
- [x] configs from environment
- [x] go modules
- [x] GORM atabase models, using UUID as identifiers
- [x] Dockerfile with build layer
- [x] docker-compose for single machine deployment 
- [x] JWT Authentication


## Pending 
- [ ] Replace GORM for a better solution, without so many abstrations 
- [ ] improve database migration
- [ ] Validate Request Payloads
- [ ] better error handling, with "stack trace" build, using errors.Wrap 
- [ ] metrics using influx and grafana 
- [ ] > 90% test coverage

## Running:

### create a .env file containing application environment variables:
```sh
APP_TOKEN=aRandomGeneratedString
DB_NAME=databaseName
DB_USER=databaseUsername
DB_PASSWORD=databasePassword
DB_HOST=0.0.0.0
DB_PORT=5432
IMAGE_NAME=dockerRepositoryImage
IMAGE_VERSION=dockerTagVersion
GIN_MODE=release 
```

### vscode configuration file for debugging
```json
{
    // Use IntelliSense to learn about possible attributes.
    // Hover to view descriptions of existing attributes.
    // For more information, visit: https://go.microsoft.com/fwlink/?linkid=830387
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Launch",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "${workspaceFolder}/main.go",
            "envFile": "${workspaceFolder}/.env",
        }
    ]
}
```

### deploy a local database server for development 
```sh
$ make database
```

### show database logs 
```sh
$ docker logs financedb
```

### build docker image and execute development environment with a local database 
```sh
$ make dev
```

