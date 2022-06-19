# self-updating container

This container keeps updating itself whenever it finds there is a new release.

## Demo

[![asciicast](https://asciinema.org/a/502745.svg)](https://asciinema.org/a/502745)

## How does it work?

CI/CD pipeline used: Github Actions

- pipeline gets triggered when a version tag `v*` (e.g. `v0.0.1`) is pushed to the repo;
- pipeline builds the golang application and the container image with `/version` file in it containing the current image tag;
- container app runs the `update` script along which monitors the new version by checking the releases page via github API;
- if it finds the new version, it then:
  - stops the app (waits for the full stop);
  - downloads and extracts the new version over;
  - starts the new app;

## Motivation

While seeking for an alternative solution to `presearch/auto-updater` (a watchtower which wants CRI socket to be exposed) image which automates updates automatically (image pull & container restart), I've developed a self-updating container as a PoC, which properly handles signals and reaps child processes (so no defunct aka zombie processes in the container), and it can be easily extended with more services should one want, by placing them into `docker/service/<app-name>/run` (`run` file can be a shell script / or a binary directly).

I know that having multiple services in a single docker container is considered evil by some (including myself), but until Akash supports all use-cases such as this or automated periodic backup service, ... it's probably best to use this image since it runs multiple services properly (handling signals, reaping zombies), plus the image is below 16 MiB, alpine-based, but can be easily ported to be ubuntu-based.

If looking forward to supporting self-updating images similar to how [watchtower](https://github.com/containrrr/watchtower) does it, but in K8s rather than inside the container itself, then there is [Keel](https://keel.sh) we may consider.

### Pros

- Users automatically get the latest version either while the container is running or when they start it for the first time;
- It does not require CRI socket API exposure (i.e. `/var/run/docker.sock`, `run/containerd/containerd.sock`, ...) as found in solutions such as watchtower. Thus reduces the attack surface;
- This container uses a custom init script and [runit](http://smarden.org/runit/) to manage services better:
  - it properly reaps child processes, i.e. no "defunct" aka "zombie" processes; You can read more on zombie reaping problem [here](https://blog.phusion.nl/2015/01/20/docker-and-the-pid-1-zombie-reaping-problem/);
  - it properly handles signals, i.e. it will relay SIGTERM to its services and wait for them to finish (7 seconds by default);

### Cons

- The image tag won't correspond to the actual version running inside the container which might not be desirable;
- It relies on github API;

## Using this image

There is Akash deployment manifst file [deploy.yaml](./deploy.yaml) which you can use to deploy it on Akash decentralized cloud!

To learn how to deploy on Akash click [here](https://docs.akash.network/guides)

### Environment variables

Set these variables:

- `GH_USER` - your github username;
- `GH_REPO` - your github repository;

Optional:

- `VERBOSE` - set to any non-empty value to increase verbosity;
- `SLEEP` - set to any value higher than `300` seconds (default), this tells how frequently to check for the new release. This is due to Github API Rate limits. For unauthenticated requests, the rate limit allows for up to 60 requests per hour. So make sure to not set sleep to a lower than `60` seconds value;
- `SVWAIT` - override the default 7 seconds to wait for the runit (`sv stop/start <service>`) command to take effect;
