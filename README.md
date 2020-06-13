[![Build Status](https://travis-ci.com/erickoliv/finances-api.svg?branch=master)](https://travis-ci.com/erickoliv/finances-api)
[![codecov](https://codecov.io/gh/erickoliv/finances-api/branch/master/graph/badge.svg)](https://codecov.io/gh/erickoliv/finances-api)
[![Go Report Card](https://goreportcard.com/badge/github.com/erickoliv/finances-api)](https://goreportcard.com/report/github.com/erickoliv/finances-api)
[![HitCount](http://hits.dwyl.io/erickoliv/erickoliv/finances-api.svg)](http://hits.dwyl.io/erickoliv/erickoliv/finances-api)
> This is just a personal project I'm using to learn the Go programming language and its libraries.

*Using different package approaches to validate the better structure for each use case*

#### The road so far
- [x] configs from environment
- [x] go modules
- [x] GORM atabase models, using UUID as identifiers
- [x] Dockerfile with build layer
- [x] docker-compose for single machine deployment 
- [x] JWT Authentication

#### In progress
- CRUD for entry tags

#### Pending 
- [ ] improve database migration
- [ ] Remove common and utils packages
- [ ] Remove DB connection from gin.Context !!!!
- [ ] Validate Request Payloads
- [ ] better error handling, with "stack trace" build, using errors.Wrap 
- [ ] metrics using influx and grafana 
- [ ] > 90% test coverage