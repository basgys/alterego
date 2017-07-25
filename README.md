# Alter Ego [![CircleCI](https://circleci.com/gh/basgys/alterego.svg?style=svg)](https://circleci.com/gh/basgys/alterego) [![Go Report Card](https://goreportcard.com/badge/github.com/basgys/alterego)](https://goreportcard.com/report/github.com/basgys/alterego)

alter ego; noun; a person's secondary or alternative personality.

## Introduction

Alter Ego is a tiny HTTP router that can redirect requests based on a set of rules.

It has been built primarily to solve the well-known problem of handling similar domain names, such as misspelled domains or domains with other TLDs.

Example:

```
http://gooogle.com -> https://google.com
http://google.co -> https://google.com
http://google.us -> https://google.com
```

## Installing

### Docker

It has been designed to run on the cloud with Docker, so you can run it with Docker like that:

```
docker pull basgys/alterego
docker run -t basgys/alterego
```

### Go binary

Or you can download the source code and run it like that:

```
go get github.com/basgys/alterego
cd $workspace/src/github.com/basgys/alterego && go run main.go
```

## Configuration

The configuration is immutable, you have to redeploy `alterego` service to make the changes take effects

|Environment Variable|Default|Description|
|:-----:|:-----:|:----------|
|IP                   |0.0.0.0|IP address on which the HTTP server will listen to|
|PORT                 |8080|Port on which the HTTP server will listen to|
|REDIRECTS            |REDIRECT1|List of redirection names separated by comma, e.g. REDIRECT1,REDIRECT2,REDIRECT3. You also need to specify each redirection as separate env variables like so: REDIRECT1="http://gooogle.com,https://google.com", REDIRECT2="http://gogle.com,http://google.com", REDIRECT3="http://google.us,https://google.com"|
|REDIRECT1            |http://127.0.0.1:8080,http://localhost:8080|Default redirection, just to give an example and allow the server to start|
|REQUEST_LOGGING      |true|Log requests to stdout|
|REDIRECT_STATUS_CODE |308|Redirection status code to use. By default it is a permanent redirection, but for development purpose it is more convenient to use a temporary redirection to avoid caching.|

## Matching rules

### Glossary

* Request URL (URL requested by the client)
* Source URL (The request URL has to match this URL)
* Destination URL (Where to redirect the request)
* Redirection URL (URL returned to the client)

### Order
The rules are compared in the order defined on `$REDIRECTS`, so define the more restrictive rules first.

### Examples

|Request URL|Rules|Redirection URL|
|:-----:|:-----:|:----------|
|http://gooogle.com|http://gooogle.com - > http://google.com|http://google.com|
|http://google.us|http://google.us - > http://google.com|http://google.com|
|http://google.us/search?q=hello|http://google.us - > http://google.com|http://google.com/search?q=hello|
|http://google.us/search?q=hello|http://google.us - > http://google.com/|http://google.com|
|http://google.us/search?q=hello|http://google.us/ - > http://google.com/|404|

## Versioning

We use [SemVer](http://semver.org/) for versioning.

## Authors

* **Bastien Gysler** - *Author* - [Twitter](https://twitter.com/basgys)

See also the list of [contributors](https://github.com/basgys/alterego/contributors) who participated in this project.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details
