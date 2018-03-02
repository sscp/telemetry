# telemetry
This repo contains the car data processing server that collects [ProtocolBuffer](https://developers.google.com/protocol-buffers/) data over UDP and logs it to CSVs, binary log files and influxdb for monitoring and data exploration. Currently built for Sundae and expected to be updated for sunrise.

## Requirements
### Running
- [golang](https://golang.org/) 1.9
- [Docker CE](https://docs.docker.com/install/#cloud) (for service use only, ie with influx)

### Developing
- [golang](https://golang.org/) 1.9
- [go dep](https://github.com/golang/dep) (for adding new dependencies/libraries)

## Installing Go

Note: if you are on Mac, I *highly* recommend installing with [brew](https://brew.sh/). If you install brew, you can install go with `brew install golang`.

Otherwise, download and install from [here](https://golang.org/dl/).

After you install, you need to setup a folder as your GOPATH. There is a good guide [here](https://github.com/golang/go/wiki/SettingGOPATH). I set my folder to `~/go` (a folder in my home directory called "go"). While you are editing your shell config, it is also really helpful to add `~/go/bin` to your PATH variable to whenever you install a go program, telemetry included, it is available to run anywhere on the system.

## How to install and just collect data
1. Clone the repo: `git clone --recurse-submodules git@github.com:sscp/telemetry.git` for ssh or `https://github.com/sscp/telemetry.git` for https
2. Run `make install` to install the `telemetry` command in `$GOPATH/bin`
3. Run `telemetry collect <run_name>` to collect data (works if you have `$GOPATH/bin` in your path, otherwise run `$GOPATH/bin/telemetry collect <run_name>` directly)
4. Data will appear in the folders `blogs` and `csvs` in your current directory

## How to run the service
1. Clone the repo (Make sure to clone the `sandbox` submodule as well with `--recurse-submodules` as shown above)
2. Run `docker-compose up` in the repo, this will compile telemetry, and download and run all other services (influx, jaeger, etc)
3. Run `make install` to install the `telemetry` command in `$GOPATH/bin`
4. Run `telemetry call start <run_name>` to collect data (works if you have `$GOPATH/bin` in your path, otherwise run `$GOPATH/bin/telemetry call start <run_name>` directly)

## How to configure telemetry
Telemetry looks for a config file in `~/.telemetry.yml` in your home folder. You should create one if you want to modify settings like ports, but it should also work without any config. You can also select a config like this: `telemetry -c <configfile> collect <run_name>`.

## How to develop
If you are new go Go, I highly recommend the [tour](https://tour.golang.org). It is a great introduction that assumes basic programming skills and quickly gets you up to speed with the way Go works.

There are unit tests that can be run with `make test`.

We use [protocol buffers](https://developers.google.com/protocol-buffers/) and [GRPC](https://grpc.io/), which means we need to compile `.proto` files. This can be done with `make generate` (wrapper for `go generate`). First run `make install-tools` to install the required libraries/compilers.

### Github Workflow
We essentially follow this [workflow](https://guides.github.com/introduction/flow/). This means whenever you want to start working on something, pull the latest changes from master `git pull`, then create a new branch to work on `git checkout -b <firstname>/<feature>`. (For me this would look something like this: `git checkout -b jack/influx`, for working on adding influx). You can now work on this branch, committing as often as you like with `git commit -m "fixed this bug"`. You should also push periodically to GitHub. The first time you push a branch you run `git push -u origin <yourbranchname>` and every time after that you just run `git push`. Once you have your branch on GitHub, you should create a "Pull Request" by going to the [repo page](https://github.com/sscp/telemetry/) and clicking the "Pull Requests" tab and then clicking "New pull request". You then select your branch and create the PR, which lets us easily review code.

