[![Build Status](https://travis-ci.com/ello/ello-go.svg?token=GLeHVitCiVCzVGXgUezV&branch=master)](https://travis-ci.com/ello/ello-go)

# Ello Go
This repository contains Ello's go projects.  For now, at least, we're keeping them in the same repository for developer ease. This may change down the road.

## Getting Started

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

* `cd streams`
* The streams app depends on a running roshi (and redis) instance.  By far, the easiest way to use this is via docker.  
  * Make sure your docker-machine is running `docker-machine start default` and then `eval "$(docker-machine env default)"` to make sure this terminal session is set up to refer to it.
  * `docker-compose up -d roshi` will start a roshi in the background (omit the -d flag for foreground)
  * You then need to tell the streams app where to find it.  For both tests and normal operation, this is done via the ELLO_ROSHI_HOST environment variable. `ELLO_ROSHI_HOST="http://$(docker-machine ip default):6302" make <command>` is the general structure you can use for running commands.  You could optionally just set that environment variable (you may need to reset if the docker-machine ip changes) and just run the make commnds alone.  
    * Example of running tests:  `ELLO_ROSHI_HOST="http://$(docker-machine ip default):6302" make test`
    * Example of running tests + build + docker: `ELLO_ROSHI_HOST="http://$(docker-machine ip default):6302" make all`
* Once built (`make build`), you can run it from `bin/streams` (use the -h flag to see what env variables you can set (ELLO_ROSHI_HOST is mandatory if roshi is not running on localhost, fyi))
* If you build the docker image, you can use docker-compose to run that, as well:
  * `ELLO_ROSHI_HOST="http://$(docker-machine ip default):6302" make all`
  * `docker-compose up -d` You may want to `docker-compose stop` and `docker-compose rm` first, if you started roshi by hand earlier.  Also, again note that the -d can be omitted to foreground it.
    * Once running, you can access it at http://$(docker-machine ip default):8080 (try http://$(docker-machine ip default):808/health/check)
