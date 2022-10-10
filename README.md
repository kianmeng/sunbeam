# Running the project

## Prerequisites

- Requires [Go](https://golang.org/dl/) 1.19 or later
- Optional
  - [direnv](https://direnv.net/): loads environment variables from .envrc

## Running the project

```console
direnv allow # load the SUNBEAM_SCRIPT_DIR environment variable (you can also use `source .envrc`)
go run main.go
```

The logs are redirected to the `debug.log` file, use `tail -f debug.log` to follow them.

The scripts available in the ex

## Installing the `sunbeam` command

```console
go install
```
