[![Build Status](https://travis-ci.org/ello/streams.svg)](https://travis-ci.org/ello/streams)

# Ello Streams
This repository contains Ello's Streams wrapper for Roshi.

## Getting Set Up With Go

* Many of these steps assume you have a correctly installed and working homebrew setup. If not, please set it up.  See http://brew.sh for details.
* Make sure you have go installed/updated (currently, we're on 1.5.1):  `brew install go` or `brew upgrade;brew update go`
* Clone this repository to your gopath (see https://golang.org/doc/code.html for information on gopath)
* To get/update the rest of the tools we make use of, run `make setup`
   * Tools include:  
      * https://github.com/Masterminds/glide
      * https://github.com/alecthomas/gometalinter
      * https://github.com/emcrisostomo/fswatch
      * https://cnswww.cns.cwru.edu/php/chet/readline/rltop.html
* For some of our services, we also recommend the use of docker to ease development.  For specific details, see the individual wiki's, but we'd recommend you install docker, docker-machine and docker-compose.  Either use docker toolbox, or install via homebrew.  

### Streams - Getting Started
After following the above steps, to run/test/build the streams application:

First, you need to make sure you have glide, gometalinter(and the linters it uses), and fswatch for all make commands to work.  There is a make target in the parent directory `make get-tools` that will do this for you.

Next, you need to make sure that you have vendored all of the dependencies for the streams project.  You can either run `glide up; glide rebuild` or use the make target in this project, `make dependencies`.  

* The streams app depends on a running roshi (and redis) instance.  By far, the easiest way to use this is via docker.  
  * Make sure your docker-machine is running `docker-machine start default` and then `eval "$(docker-machine env default)"` to make sure this terminal session is set up to refer to it.
  * `docker-compose up -d roshi` will start a roshi in the background (omit the -d flag for foreground)
  * You then need to tell the streams app where to find roshi.  For both tests and normal operation, this is done via the ROSHI_URL environment variable. `ROSHI_URL="http://$(docker-machine ip default):6302" make <command>` is the general structure you can use for running commands.  You could optionally just set that environment variable (you may need to reset if the docker-machine ip changes) and just run the make commnds alone.  
    * Example of running tests:  `ROSHI_URL="http://$(docker-machine ip default):6302" make test`
    * Example of running tests + build + docker: `ROSHI_URL="http://$(docker-machine ip default):6302" make all`
* Once built (`make build`), you can run it from `bin/streams` (use the -h flag to see what env variables you can set (ROSHI_URL is mandatory if roshi is not running on localhost, fyi))
* If you build the docker image, you can use docker-compose to run that, as well:
  * `ROSHI_URL="http://$(docker-machine ip default):6302" make all`
  * `docker-compose up -d` You may want to `docker-compose stop` and `docker-compose rm` first, if you started roshi by hand earlier.  Also, again note that the -d can be omitted to foreground it.
    * Once running, you can access it at http://$(docker-machine ip default):8080 (try http://$(docker-machine ip default):808/health/check)

#### TLDR

* From _$GOPATH/src/github.com/ello/streams_, execute `make setup`
* Verify you have a working docker install with a valid docker-machine daemon connected
* `docker-compose start roshi`
* `ROSHI_URL="http://$(docker-machine ip default):6302" make test`
