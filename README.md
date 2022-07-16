# TibiaData API in Golang

[![GitHub CI](https://github.com/TibiaData/tibiadata-api-go/workflows/build/badge.svg?branch=main)](https://github.com/TibiaData/tibiadata-api-go/actions?query=workflow%3Abuild)
[![Codecov](https://codecov.io/gh/TibiaData/tibiadata-api-go/branch/main/graph/badge.svg?token=PSBNLBI10C)](https://codecov.io/gh/TibiaData/tibiadata-api-go)
[![GitHub go.mod version](https://img.shields.io/github/go-mod/go-version/tibiadata/tibiadata-api-go)](https://github.com/tibiadata/tibiadata-api-go/blob/main/go.mod)
[![Docker version](https://img.shields.io/docker/v/tibiadata/tibiadata-api-go/latest)](https://hub.docker.com/r/tibiadata/tibiadata-api-go)
[![Docker size](https://img.shields.io/docker/image-size/tibiadata/tibiadata-api-go/latest)](https://hub.docker.com/r/tibiadata/tibiadata-api-go)
[![GitHub license](https://img.shields.io/github/license/tibiadata/tibiadata-api-go)](https://github.com/tibiadata/tibiadata-api-go/blob/main/LICENSE)

TibiaData API written in Golang and deployed in container (which contains v3).

Current status of v3 is released and information like documentation can be found on [docs.tibiadata.com](https://docs.tibiadata.com).

### Table of Contents

- [How to use](#how-to-use)
  - [Docker](#docker)
  - [Docker-compose](#docker-compose)
  - [Local development](#local-development)
  - [Environment variables](#environment-variables)
  - [Deployment note](#deployment-note)
- [API documentation](#api-documentation)
  - [Available endpoints](#available-endpoints)
- [General information](#general-information)
- [Credits](#credits)

## How to use

You can either use it in a Docker container or go download the code and deploy it yourself on any server.

Keep in mind that there are restrictions on tibia.com that might impact the usage of the application being hosted yourself.

### Docker

Our images are available on both [GitHub Container Registry](https://github.com/TibiaData/tibiadata-api-go/pkgs/container/tibiadata-api-go) and [Docker Hub](https://hub.docker.com/r/tibiadata/tibiadata-api-go).

This is how to pull and run the _latest_ release of TibiaData from [GHCR](https://github.com/TibiaData/tibiadata-api-go/pkgs/container/tibiadata-api-go):

```console
docker pull ghcr.io/tibiadata/tibiadata-api-go:latest
docker run -p 127.0.0.1:80:8080/tcp --rm -it ghcr.io/tibiadata/tibiadata-api-go:latest
```
You can also use [Docker Hub](https://hub.docker.com/r/tibiadata/tibiadata-api-go) to pull your images from.

If you want to run the latest code you can switch from _latest_ to _edge_.

### Docker-compose

This is a simple example on how you can get up and running with TibiaData in docker-compose, which will be running on port 8080 and be exposed locally.

```yaml
version: "3"

services:
  tibiadata:
    image: ghcr.io/tibiadata/tibiadata-api-go:latest
    restart: always
    environment:
      - TIBIADATA_HOST=tibiadata.example.com
    ports:
      - 8080:8080
```

### Local development

Build the code on your computer

```console
docker build -t tibiadata .
```

Run your build locally

```console
docker run -p 127.0.0.1:80:8080/tcp --rm -it tibiadata
```

### Environment variables

_Information will be added at a later stage._

### Deployment note

You should consider to add a layer in front of this application, so you can do caching of endpoints, access controll or what ever your needs are.

We do so at least by using [Kong](https://github.com/Kong/kong) API Gateway, which solves features like caching, rate-limiting, authentication and more.

## API documentation

The hosted API documentation for our [api.tibiadata.com](https://api.tibiadata.com) service can be viewd at [docs.tibiadata.com](https://docs.tibiadata.com).

There is a swagger-generated documentation available for download on the [GitHub Release](https://github.com/TibiaData/tibiadata-api-go/releases) of the version you are looking for.

### Available endpoints

Those are the current existing endpoints.

- GET `/ping`
- ~GET `/health`~
- GET `/healthz`
- GET `/readyz`
- GET `/v3/boostablebosses`
- GET `/v3/character/:name`
- GET `/v3/creature/:race`
- GET `/v3/creatures`
- GET `/v3/fansites`
- GET `/v3/guild/:name`
- GET `/v3/guilds/:world`
- GET `/v3/highscores/:world/:category/:vocation`
- GET `/v3/house/:world/:house_id`
- GET `/v3/houses/:world/:town`
- GET `/v3/killstatistics/:world`
- GET `/v3/news/archive`
- GET `/v3/news/archive/:days`
- GET `/v3/news/id/:news_id`
- GET `/v3/news/latest`
- GET `/v3/news/newsticker`
- GET `/v3/spell/:spell_id`
- GET `/v3/spells`
- GET `/v3/world/:name`
- GET `/v3/worlds`
- GET `/versions`

## General information

Tibia is a registered trademark of [CipSoft GmbH](https://www.cipsoft.com/en/). Tibia and all products related to Tibia are copyright by [CipSoft GmbH](https://www.cipsoft.com/en/).

## Credits

- Authors: [Tobias Lindberg](https://github.com/tobiasehlert) – [List of contributors](https://github.com/TibiaData/tibiadata-api-go/graphs/contributors)
- Distributed under [MIT License](LICENSE)
