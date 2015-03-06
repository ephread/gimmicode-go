# Gimmicode Go #

## Description ##

Gimmicode is the contraction of "Give me Unicode".

This is a web-accessible, go port of the old
[gimmicode](https://github.com/ephread/gimmicode) command-line utility.

## Building Gimmicode ##

Just run `docker build` inside the repository.

```
$ docker build -t gimmicode .
```

## Running Gimmicode ##

You will also need a linked container aliased as `redis`, which
will run a redis database. You can use the official
[docker image](https://registry.hub.docker.com/_/redis/).

```
$ docker run --name gimmicode-redis \
         -v /docker/gimmicode/data:/data \
         -d redis

$ docker run --name gimmicode-go \
         -p 3000:3000 \
         --link gimmicode-redis:redis \
         -d gimmicode-go 
```

## Author ##

Gimmicode was made by [Frédéric Maquin](ephread.com).

## Unlicense ##

This is public domain. Do what you want with it.

Please see the UNLICENSE included for details.