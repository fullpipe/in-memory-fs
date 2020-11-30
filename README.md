# ngserve

[![Tests Status](https://github.com/fullpipe/ngserve/workflows/Tests/badge.svg)](https://github.com/fullpipe/ngserve)
[![Docker Image](https://img.shields.io/docker/image-size/fullpipe/ngserve/latest)](https://cloud.docker.com/repository/docker/fullpipe/ngserve)

Simple and easy to use http server for angular 2+ apps.

## Usage

Add `Dockerfile` to the root of your angular app project.

```Dockerfile
# Build
FROM node:lts-alpine AS build

WORKDIR /app

COPY package.json package-lock.json ./
RUN npm install

COPY . .

RUN npm run build --prod

# App image
FROM fullpipe/ngserve:latest

COPY --from=build /app/dist/example/ /app/
```

First stage will build your app. Second will copy `dist` to `ngserve` web root
directory.

Or if you build on your own

```Dockerfile
FROM fullpipe/ngserve:latest

# Do not forget to add end-slashes to copy dir content
COPY dist/example/ /app/
```

Now you can build and run you app

```bash
docker build -t example .
docker run -p 8080:8080 example
```

## Example

See and try [example](https://github.com/fullpipe/ngserve/tree/main/example).

```bash
cd example
npm i
npm run build
docker-compose up --build
```

It will create three different endpoints:

- http://localhost:8080/
- http://localhost:8081/
- http://localhost:8082/

## TODO:

- benchmarks against nginx. Service uses in memory file cache. So it could be a
  little bit faster then default nginx
- logs?
