# Forking this project

Once you fork this project, make sure to generate the Personal Access Token (PAT),
and then create the environment variable as described below.

This is needed by the GitHub Actions CI/CD pipeline so it can upload new docker
images to the GitHub Packages registry. (`ghcr.io` registry).

## Creating the PAT

Head to https://github.com/settings/tokens/new?scopes=write:packages page,
where you can create your PAT with the limited scope to upload packages to
GitHub Package Registry - `write:packages` only.

## Adding the environment variable

Go to your forked project's `Settings` tab, and then head to
`Secrets` -> `Actions` -> `New repository secret` where you can create the
following environment variable with your PAT as the value:

- `GH_PAT_WRITE_PACKAGES` - used by [release.yml](.github/workflows/release.yml)

## How about Docker Hub?

If you want to store your images at the [Docker Hub](https://hub.docker.com)
instead of GHCR, you absolutely can use it and it is simple to switch to it:

1. Create a new Access Token (read-write permissions) at https://hub.docker.com/settings/security
2. Set `GH_PAT_WRITE_PACKAGES` environemnt variable to your new Docker Hub access token.
3. Remove `registry:` line from the `Login to Docker Registry` (`docker/login-action`) step in [release.yml](.github/workflows/release.yml) file so the pipeline will default to the Docker Hub registry.

