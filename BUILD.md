# Building

Normally the CI/CD pipeline would build it all, but if you want to do it manually, below are the steps.

## The app

```
GOPATH=$PWD GO111MODULE=off go build -v app
```

## The image

```
docker build -t app -f docker/Dockerfile .
```

## Running

Now you can run it:

```
docker run \
  --rm \
  -e VERSION="$(git describe --tags --abbrev=0)" \
  -e GH_USER=arno01 \
  -e GH_REPO=self-update \
  -e SLEEP=62 \
  -ti \
  app
```
