# Setup

First, you need to make sure you have glide, gometalinter(and the linters it uses), and fswatch for all make commands to work.  There is a make target in the parent directory `make get-tools` that will do this for you.

Next, you need to make sure that you have vendored all of the dependencies for the streams project.  You can either run `glide up; glide rebuild` or use the make target in this project, `make dependencies`.  

As a convienance two steps can be accomplished in one command from the parent project, ello-go, of `make setup`

Finally, for tests to pass, you need to have an instance of roshi running.  You can install it yourself locally, or you can use a docker image of it (recommended).  To do this, make sure you have docker properly installed and a docker-machine instance running (confirm via `docker-machine ip default` returning an ip address).  Then, you bring up roshi with:

`docker-compose start roshi` or `docker-compose up roshi` if you prefer it in the foreground.  

You then need to tell the streams app where to find it.  For both tests and normal operation, this is done via the ELLO_ROSHI_HOST environment variable.  

An example of a one off run of the tests would then look like:

`ELLO_ROSHI_HOST="http://$(docker-machine ip default):6302" make test`

You can, of course, set that environment variable so you don't need to prefix it for every run.

## TLDR

* From _$GOPATH/src/github.com/ello/ello-go_, execute `make setup`
* Verify you have a working docker install with a valid docker-machine daemon connected
* `docker-compose start roshi`
* `ELLO_ROSHI_HOST="http://$(docker-machine ip default):6302" make test`
