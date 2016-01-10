<img src="http://d324imu86q1bqn.cloudfront.net/uploads/user/avatar/641/large_Ello.1000x1000.png" width="200px" height="200px" />

# Streams - Roshi-based activity feeds
Streams is a RESTful Go wrapper for [Soundcloud's Roshi](https://github.com/soundcloud/roshi), an awesome tool for building activity feeds. Streams improves upon the built-in `roshi-server` REST API by mapping some of its low-level concepts into higher-level ones and using more conventional REST semantics.

[![Build Status](https://travis-ci.org/ello/streams.svg)](https://travis-ci.org/ello/streams)

### Quickstart

* Install Go 1.5
* Clone this repo to `$GOPATH/src/github.com/ello/streams`
* From `$GOPATH/src/github.com/ello/streams`, execute `make setup`
* Verify you have a working docker install with a valid docker-machine daemon connected
* Fire up a Roshi instance by executing `docker-compose start roshi`
* Run the tests with `ROSHI_URL="http://$(docker-machine ip default):6302" make test`

## Overview
The Streams service acts as an intermediate layer between Roshi and our Rails application. It essentially acts as a replacement for making requests against the activities table to load stream data for a user.  You query the Streams service with the ids of the users you wish to have a stream of and it will return IDs that you can then query directly from Postgres.  

### What is Roshi?
[Roshi](https://github.com/soundcloud/roshi) is an open source software product originally written by the engineers at Soundcloud to [power their activity feeds](https://developers.soundcloud.com/blog/roshi-a-crdt-system-for-timestamped-events). Rather than using a Fan-Out-On-Write model, it uses a Fan-In-On-Read approach, modeled on [CRDTs](https://en.wikipedia.org/wiki/Conflict-free_replicated_data_type).

Fan-Out-On-Write is one model for managing a 'Twitter-like' content system. Every time a user posts, each follower of that user gets a record written to a collection associated with them. At read time, you can simply read a users collection to load the content they should see. This trades *O(N)* writes, where *N* is the follower count, for near-constant-time reads. These are desirable characteristics when writes are not in the critical path of a single web request, but reads are.

Fan-In-On-Read is an alternative means for accomplishing the same thing. In this approach, every time a user posts, a single record is added to a collection associated to the posting user. When a user loads their stream of content, at read time you request content from each of the followed users' collections and combine them into a single view. This trades potentially variable *O(N)* performance at read time for reduced storage requirements and constant-time write performance. 

Roshi persists data using a number of [Redis](http://redis.io/) instances, laid out in multiple clusters for durability, and makes use of bounded-size [sorted sets](http://redis.io/commands#sorted_set).

The core type that Roshi uses is:

```
type KeyScoreMember struct {
	Key string,
	Score float64,
	Member string,
}
```

`Key` is used as a stream identifier. In Ello's case, it is a user's ID hashed with a non-cryptographic hash ([xxhash](https://github.com/OneOfOne/xxhash)).
`Score` is used to order the items. In Ello's case, it is the timestamp when the post was created, converted to a float value (based on nanoseconds from epoch).
`Member` is used to store information that identifies the exact item inserted into the stream. In Ello's case, it is a JSON object which includes the post ID, posting user ID and the type (Post/Repost/etc.). This allows us to load a time ordered stream of posts authored by an individual user.

Critically, Roshi supports an efficient coalesce function, which allows us to load multiple user streams into a single consolidated, time-ordered stream. This is the primary access avenue for Ello. When loading content for a user, we will load the user ids for the accounts that user follows (adjusting as needed for any blocked/blocking users) and request a coalesced stream of those ids from the Streams service (and thus, from Roshi).

#### Pagination
Pagination in Roshi is a little complicated and thus deserves a bit of discussion.

Natively, Roshi supports two methods of pagination:

##### Limit/Offset
In many data systems, Limit refers to how many records to return. Offset refers to how far down the list of records to move before you start returning them.

Given an ordered list `[A,B,C,D,E,F,G,H]`, a request with a limit of 2 and an offset of 0 would return `[A,B]`. A request with a limit of 2, and an offset of 2 would return `[C,D]`.

The combination of limit and offset can be used to partition a set of data into pages. This works well with static data, but in a system with frequent inserts, it has limitations that are difficult to overcome. Specifically, as the head continues to move, the offset is often incorrect, resulting in duplicated entries.  

##### Cursor-based
Cursor-based pagination works similarly to a limit/offset-based system, but rather than using an offset count to describe where to begin returning records, it uses an actual record. Given the ordered nature of Roshi collections, this effectively eliminates the inaccuracies that can arise with pure limit/offset where the head is frquently changing.

This approach has two limitations, however. The first is that it requires the calling client to keep track of the cursor and use/update it on subsequent requests. The second is that cursor-based performance degrades the deeper into the collection you retrieve. For systems like Ello, users are typically looking at or near the head of the content collection, so this effect is somewhat mitigated.  

The format of a Roshi cursor is a bit odd and worth calling out here:

```
<IEEE 754 binary representation of nanoseconds since epoch in base 10>A<base64 representation of the Member>
 ```

For example:

```
4894443175316128785AMWNmMjYyM2QtYmExNi00N2VmLWE2ZTktNmU1NTE1MzNiOdNk
```

#### References
- Roshi Server Documentation: https://github.com/soundcloud/roshi/blob/master/roshi-server/README.md
- Roshi Overview: https://github.com/soundcloud/roshi/blob/master/README.md

### What is Streams?
The Streams service is an intermediary service, written in Go. It is structured as a fairly standard Go REST API.  It is using [httprouter](https://github.com/julienschmidt/httprouter), [negroni](https://github.com/codegangsta/negroni), and [render](https://github.com/unrolled/render) but otherwise is built on the stdlib for the actual REST interface. The entry point for the service is in `streams/main.go`. This reads environment variables for configuration and sets up the application. The bulk of the API code is in the `streams/api` package. `streams/model` contains representations of both the objects we use for REST communication to clients, as well as those for communicating to Roshi. It also handles the translation between those two worlds. `streams/service` contains the necessary code for interfacing with the underlying Roshi instance. `streams/util` has other random bits of useful common code (validation, environment variable helpers, etc).

The project has been set up to vendor its dependencies, using the Go 1.5 experimental feature. It can be a little tricky to get this setup correctly, but there is a `Makefile` that handles most of this for you.

### Motivations
For Ello, adopting a fan-in model has a few distinct advantages. Given that our user base visits the site at widely varying frequencies, a fan-out model incurs a large cost to store content on behalf of users who may not visit frequently enough to see all of it. Additionally, this approach lowers the cost of adding/removing relationships (whether through blocking, onboarding, etc.) with full history (e.g., you can immediately can see a new follower's entire history in your stream).

#### Cost Breakdown
Ello current utilizes a sharded array of dedicated Heroku Postgres instances to handle activity feeds (friend/noise streams, notifications, etc). While effective, this array constitutes our largest single fixed engineering cost, and the manual effort required to scale shards is a bottleneck for future growth.

Streams separates the handling of this to a new service and adopt a new approach for storage, providing to us a substantial cost benefit and increased scaling capabilities.

##### Current Costs
5 Heroku Postgres Standard-6 @ $2K/month/instance + 33 2x worker instances @ $50/month/instance = $11,650 per month (average)

##### Estimated Costs
- $1619.61 per month on average for Redis instances (includes upfront cost amortized) - 9 `cache.r3.xlarge` Elasticache instances (using reserved pricing)
- ~ $500 per month for EC2 instances for Streams API - 2 m4.2xlarge (reserved, no upfront)
- = $2200 per month total cost

##### Deployment, Operations, and Gotchas
To be written

## Development

### Getting Set Up With Go

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

### Development
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

## License
Streams is released under the [MIT License](blob/master/LICENSE.txt)

## Code of Conduct
Ello was created by idealists who believe that the essential nature of all human beings is to be kind, considerate, helpful, intelligent, responsible, and respectful of others. To that end, we will be enforcing [the Ello rules](https://ello.co/wtf/policies/rules/) within all of our open source projects. If you don’t follow the rules, you risk being ignored, banned, or reported for abuse.

