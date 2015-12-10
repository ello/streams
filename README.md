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
