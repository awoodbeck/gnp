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
ok  	github.com/awoodbeck/gnp/ch03	22.757s
PASS
ok  	github.com/awoodbeck/gnp/ch04	1.192s
PASS
ok  	github.com/awoodbeck/gnp/ch05/echo	12.048s
?   	github.com/awoodbeck/gnp/ch06/sha512-256sum	[no test files]
2025/08/13 19:37:23 [127.0.0.1:51739] requested file: test
2025/08/13 19:37:23 [127.0.0.1:51739] sent 148 blocks
PASS
ok  	github.com/awoodbeck/gnp/ch06/tftp	1.056s
?   	github.com/awoodbeck/gnp/ch06/tftp/tftp	[no test files]
?   	github.com/awoodbeck/gnp/ch07/creds	[no test files]
?   	github.com/awoodbeck/gnp/ch07/creds/auth	[no test files]
goos: linux
goarch: amd64
pkg: github.com/awoodbeck/gnp/ch07/echo
cpu: AMD Ryzen 5 PRO 2500U w/ Radeon Vega Mobile Gfx
BenchmarkEchoServerUnixPacket-8     	   34728	     34472 ns/op
BenchmarkEchoServerUnixDatagram-8   	   24410	     43301 ns/op
BenchmarkEchoServerUDP-8            	   29680	     42264 ns/op
BenchmarkEchoServerTCP-8            	   19546	     60073 ns/op
BenchmarkEchoServerUnix-8           	   37960	     32372 ns/op
PASS
ok  	github.com/awoodbeck/gnp/ch07/echo	10.261s
PASS
ok  	github.com/awoodbeck/gnp/ch08	7.029s
PASS
ok  	github.com/awoodbeck/gnp/ch09	1.045s
PASS
ok  	github.com/awoodbeck/gnp/ch09/handlers	1.034s
PASS
ok  	github.com/awoodbeck/gnp/ch09/middleware	2.057s
?   	github.com/awoodbeck/gnp/ch10	[no test files]
?   	github.com/awoodbeck/gnp/ch10/backend	[no test files]
2025/08/13 19:37:46 http: TLS handshake error from 127.0.0.1:33652: remote error: tls: bad certificate
PASS
ok  	github.com/awoodbeck/gnp/ch11	3.298s
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
ok  	github.com/awoodbeck/gnp/ch13	2.051s
?   	github.com/awoodbeck/gnp/ch13/instrumentation	[no test files]
?   	github.com/awoodbeck/gnp/ch13/instrumentation/metrics	[no test files]
?   	github.com/awoodbeck/gnp/ch14/aws	[no test files]
?   	github.com/awoodbeck/gnp/ch14/azure	[no test files]
?   	github.com/awoodbeck/gnp/ch14/gcp	[no test files]
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
root@e171155cb667:/usr/src/gnp# go test -race ./...
ok  	github.com/awoodbeck/gnp/ch03	22.765s
ok  	github.com/awoodbeck/gnp/ch04	1.217s
ok  	github.com/awoodbeck/gnp/ch05/echo	12.036s
?   	github.com/awoodbeck/gnp/ch06/sha512-256sum	[no test files]
ok  	github.com/awoodbeck/gnp/ch06/tftp	1.053s
?   	github.com/awoodbeck/gnp/ch06/tftp/tftp	[no test files]
?   	github.com/awoodbeck/gnp/ch07/creds	[no test files]
?   	github.com/awoodbeck/gnp/ch07/creds/auth	[no test files]
ok  	github.com/awoodbeck/gnp/ch07/echo	1.030s
ok  	github.com/awoodbeck/gnp/ch08	7.371s
ok  	github.com/awoodbeck/gnp/ch09	1.039s
ok  	github.com/awoodbeck/gnp/ch09/handlers	1.035s
ok  	github.com/awoodbeck/gnp/ch09/middleware	2.092s
?   	github.com/awoodbeck/gnp/ch10	[no test files]
?   	github.com/awoodbeck/gnp/ch10/backend	[no test files]
ok  	github.com/awoodbeck/gnp/ch11	3.375s
?   	github.com/awoodbeck/gnp/ch11/cert	[no test files]
?   	github.com/awoodbeck/gnp/ch12/client	[no test files]
?   	github.com/awoodbeck/gnp/ch12/cmd	[no test files]
?   	github.com/awoodbeck/gnp/ch12/gob	[no test files]
?   	github.com/awoodbeck/gnp/ch12/housework	[no test files]
?   	github.com/awoodbeck/gnp/ch12/housework/v1	[no test files]
?   	github.com/awoodbeck/gnp/ch12/json	[no test files]
?   	github.com/awoodbeck/gnp/ch12/protobuf	[no test files]
?   	github.com/awoodbeck/gnp/ch12/server	[no test files]
ok  	github.com/awoodbeck/gnp/ch13	2.050s
?   	github.com/awoodbeck/gnp/ch13/instrumentation	[no test files]
?   	github.com/awoodbeck/gnp/ch13/instrumentation/metrics	[no test files]
?   	github.com/awoodbeck/gnp/ch14/aws	[no test files]
?   	github.com/awoodbeck/gnp/ch14/azure	[no test files]
?   	github.com/awoodbeck/gnp/ch14/gcp	[no test files]
root@e171155cb667:/usr/src/gnp#
```

## What's `gnp` in the context of this repository?

`gnp` was the acronym for the book's working name, "Go Network Programming." The book's name evolved 
while in development, but the repository did not. Before publishing, my energy was entirely focused
on completing the book, and I couldn't justify renaming this repository and correcting all
references to it while the book's deadline loomed. Perhaps in the second edition.

### Updates

* _August 2025_ -- Update to use Go release 1.25.0. Test on Linux Mint 22.1. Docker version 28.3.3, build 980b856.

* _February 2025_ -- Updated to use Go release 1.24.0. Tested on Linux Mint 22.1. Docker version 27.5.1, build 9f9e405.
