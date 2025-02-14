# Network Programming with Go
Code repository for [Network Programming with Go](https://nostarch.com/networkprogrammingwithgo) from
No Starch Press.

Although this book was targeted for developers familiar with the Go programming language, it's
reasonable to assume that you may have picked it up early in your journey of mastering Go. If you
aren't comfortable running the tests and examples presented in the book you can either clone this
repository and run it on your operating system's command line, or run them in a Docker container.

## Running the examples on your command line:

Make sure you have `git` installed on your command line. If not, these [instructions](https://docs.github.com/en/github-cli)
should get you started.

First, clone this repository by clicking on the green `Code` button near the top of this page and
selecting the appropriate command. To clone the repository over HTTPS, run:

    git clone https://github.com/awoodbeck/gnp.git

To clone the repository over SSH, run this command:

    git clone git@github.com:awoodbeck/gnp.git

Once cloned, you can change into the `gnp` directory and run all of the tests, like so:

    cd gnp
    go test -timeout 300s -race -bench=. ./...

with the following typical results:
```fish
PASS
ok  	github.com/awoodbeck/gnp/ch03	22.760s
PASS
ok  	github.com/awoodbeck/gnp/ch04	1.212s
PASS
ok  	github.com/awoodbeck/gnp/ch05/echo	12.046s
?   	github.com/awoodbeck/gnp/ch06/sha512-256sum	[no test files]
2024/12/16 17:38:05 [127.0.0.1:44793] requested file: test
2024/12/16 17:38:05 [127.0.0.1:44793] sent 148 blocks
PASS
ok  	github.com/awoodbeck/gnp/ch06/tftp	1.050s
?   	github.com/awoodbeck/gnp/ch06/tftp/tftp	[no test files]
?   	github.com/awoodbeck/gnp/ch07/creds	[no test files]
?   	github.com/awoodbeck/gnp/ch07/creds/auth	[no test files]
goos: linux
goarch: amd64
pkg: github.com/awoodbeck/gnp/ch07/echo
cpu: AMD Ryzen 5 PRO 2500U w/ Radeon Vega Mobile Gfx
BenchmarkEchoServerUnixPacket-8     	   32089	     33895 ns/op
BenchmarkEchoServerUnixDatagram-8   	   28700	     36013 ns/op
BenchmarkEchoServerUDP-8            	   31693	     39171 ns/op
BenchmarkEchoServerTCP-8            	   24196	     58420 ns/op
BenchmarkEchoServerUnix-8           	   43057	     35109 ns/op
PASS
ok  	github.com/awoodbeck/gnp/ch07/echo	9.341s
PASS
ok  	github.com/awoodbeck/gnp/ch08	6.897s
PASS
ok  	github.com/awoodbeck/gnp/ch09	1.045s
PASS
ok  	github.com/awoodbeck/gnp/ch09/handlers	1.025s
PASS
ok  	github.com/awoodbeck/gnp/ch09/middleware	2.060s
?   	github.com/awoodbeck/gnp/ch10	[no test files]
?   	github.com/awoodbeck/gnp/ch10/backend	[no test files]
2024/12/16 17:38:27 http: TLS handshake error from 127.0.0.1:51674: remote error: tls: bad certificate
PASS
ok  	github.com/awoodbeck/gnp/ch11	3.393s
?   	github.com/awoodbeck/gnp/ch11/cert	[no test files]
?   	github.com/awoodbeck/gnp/ch12/client	[no test files]
?   	github.com/awoodbeck/gnp/ch12/cmd	[no test files]
?   	github.com/awoodbeck/gnp/ch12/gob	[no test files]
?   	github.com/awoodbeck/gnp/ch12/housework	[no test files]
?   	github.com/awoodbeck/gnp/ch12/housework/v1	[no test files]
?   	github.com/awoodbeck/gnp/ch12/json	[no test files]
?   	github.com/awoodbeck/gnp/ch12/protobuf	[no test files]
?   	github.com/awoodbeck/gnp/ch12/server	[no test files]
PASS
ok  	github.com/awoodbeck/gnp/ch13	2.055s
?   	github.com/awoodbeck/gnp/ch13/instrumentation	[no test files]
?   	github.com/awoodbeck/gnp/ch13/instrumentation/metrics	[no test files]
?   	github.com/awoodbeck/gnp/ch14/aws	[no test files]
?   	github.com/awoodbeck/gnp/ch14/azure	[no test files]
```

Alternatively, run the tests from a single chapter:

    go test -v -timeout 300s -race -bench=. ./ch03/dial_fanout_test.go

## Run examples in a Docker container:

First, ensure Docker is installed by following [instructions](https://docs.docker.com/engine/install/) for your operating system.

Next, install the [latest docker-buildx release](https://github.com/docker/buildx/releases/latest) for your operating system.

Then, clone this repository and build the `gnp` docker container by running the following commands 
in the `gnp` directory:

    git clone git@github.com:awoodbeck/gnp.git
    cd gnp
    docker buildx build -t gnp .

Once finished, you should see a `gnp` image in the output of the `docker image ls`
command, like this:

```bash
$ docker image ls                                                                                                                                                                                             ✗ master
REPOSITORY   TAG       IMAGE ID       CREATED          SIZE
gnp          latest    ce9980c7834f   8 minutes ago    1.07GB
```

Finally, you can run the container using the `docker run --rm -it gnp bash`
command. You should find yourself at a bash prompt where you can run the tests
by issuing the `go test -race ./...` command, for example, as seen below:

```bash
root@7f00d5d8ad21:/usr/src/gnp# go test -race ./...
?       github.com/awoodbeck/gnp/ch06/sha512-256sum     [no test files]
?       github.com/awoodbeck/gnp/ch06/tftp/tftp [no test files]
?       github.com/awoodbeck/gnp/ch07/creds     [no test files]
?       github.com/awoodbeck/gnp/ch07/creds/auth        [no test files]
ok      github.com/awoodbeck/gnp/ch03   21.764s
ok      github.com/awoodbeck/gnp/ch04   0.067s
ok      github.com/awoodbeck/gnp/ch05/echo      11.039s
ok      github.com/awoodbeck/gnp/ch06/tftp      0.030s
ok      github.com/awoodbeck/gnp/ch07/echo      0.021s
ok      github.com/awoodbeck/gnp/ch08   5.733s
ok      github.com/awoodbeck/gnp/ch09   0.019s
ok      github.com/awoodbeck/gnp/ch09/handlers  0.019s
ok      github.com/awoodbeck/gnp/ch09/middleware        1.038s
?       github.com/awoodbeck/gnp/ch10   [no test files]
?       github.com/awoodbeck/gnp/ch10/backend   [no test files]
?       github.com/awoodbeck/gnp/ch11/cert      [no test files]
?       github.com/awoodbeck/gnp/ch12/client    [no test files]
?       github.com/awoodbeck/gnp/ch12/cmd       [no test files]
?       github.com/awoodbeck/gnp/ch12/gob       [no test files]
?       github.com/awoodbeck/gnp/ch12/housework [no test files]
?       github.com/awoodbeck/gnp/ch12/housework/v1      [no test files]
?       github.com/awoodbeck/gnp/ch12/json      [no test files]
?       github.com/awoodbeck/gnp/ch12/protobuf  [no test files]
?       github.com/awoodbeck/gnp/ch12/server    [no test files]
?       github.com/awoodbeck/gnp/ch13/instrumentation   [no test files]
?       github.com/awoodbeck/gnp/ch13/instrumentation/metrics   [no test files]
?       github.com/awoodbeck/gnp/ch14/aws       [no test files]
?       github.com/awoodbeck/gnp/ch14/azure     [no test files]
?       github.com/awoodbeck/gnp/ch14/gcp       [no test files]
ok      github.com/awoodbeck/gnp/ch11   2.305s
ok      github.com/awoodbeck/gnp/ch13   1.076s
```

## What's `gnp` in the context of this repository?

`gnp` was the acronym for the book's working name, "Go Network Programming." The book's name evolved 
while in development, but the repository did not. Before publishing, my energy was entirely focused
on completing the book, and I couldn't justify renaming this repository and correcting all
references to it while the book's deadline loomed. Perhaps in the second edition.

### Updates

* _February 2025_ -- Updated to use Go release 1.24.0. Tested on Linux Mint 22.1. Docker version 27.5.1, build 9f9e405.
