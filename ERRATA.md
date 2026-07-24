# Errata

Corrections to *Network Programming with Go* discovered since publication.
Each has been submitted to No Starch Press for future printings. Where a
correction affects code, this repository has been updated to match, so the
code here may intentionally differ from the printed listings as noted below.

## Introduction

### Page xx: minimum Go version requirements have moved

The "Installing Go" section says most programs in the book need at least
Go 1.12 and that some require Go 1.14 or newer. The floor has moved many
times since printing. This repository's [go.mod](go.mod) file records the
minimum Go version the code currently requires, and that is the number to
trust. (The `https://golang.org/` URL in the same section now redirects to
`https://go.dev/`.)

## Chapter 3: Reliable TCP Data Streams

### Page 48: TCP acknowledgment numbers name the next expected packet, not the last received one

The "Acknowledging Receipt of Packets by Using Their Sequence Numbers"
section describes an ACK as meaning "I've received all packets up to and
including the packet with this sequence number." TCP's acknowledgment
number actually names the *next* sequence number the receiver expects,
implicitly acknowledging everything before it. Figure 3-2 already reflects
this with its X + 1 and Y + 1 labels.

The retransmission example changes accordingly: a sender that transmits
packets up through 100 and receives an ACK with acknowledgment number 90
must retransmit packets 90 to 100 (the book says 91 to 100).

Two further notes for careful readers:

- The acknowledgment number is its own TCP header field, distinct from the
  sequence number, though the book's prose uses "sequence number" for both.
- TCP sequence and acknowledgment numbers count bytes in the stream, not
  whole packets. Each packet's sequence number identifies its first byte.
  The chapter numbers packets to keep the illustration simple.

Thanks to an attentive reader for catching this one.

## Chapter 9: Building HTTP Services

### Pages 208–209: `http.ServeMux` now redirects with 307 instead of 301

Go 1.26 changed the multiplexer's implicit trailing-slash redirect from
301 Moved Permanently to 307 Temporary Redirect so that redirected
requests preserve their HTTP method and body
([golang/go#50243](https://github.com/golang/go/issues/50243)).

This affects Listing 9-16 and the prose around it: the fourth test case's
expected status becomes `http.StatusTemporaryRedirect` and its expected
body becomes a "Temporary Redirect" link. On Go 1.25 and earlier, the
printed 301 behavior still holds.

The repository's [ch09/mux_test.go](ch09/mux_test.go) has been updated to
expect the Go 1.26 behavior, so it intentionally differs from the printed
listing.

## Chapter 10: Caddy: A Contemporary Web Server

### Page 219: Go version required to build Caddy from source

The "Building Caddy from Source Code" section says you need at least
Go 1.14. Caddy's minimum Go version moves with each of its releases.
Consult the `go.mod` file in Caddy's repository for the version it
currently requires.

## Chapter 12: Data Serialization

### Pages 281–293: the gRPC code-generation API changed after printing

The gRPC team removed the experimental service-struct registration API
(`RobotMaidService`, `RegisterRobotMaidService`) that the printed chapter
targets before `protoc-gen-go-grpc` reached v1.0, and the Go project has
since deprecated the `go get` tool-installation commands and the
`github.com/golang/protobuf` module the chapter references. The printed
listings no longer work with current tooling as a result (issues
[#3](https://github.com/awoodbeck/gnp/issues/3) and
[#5](https://github.com/awoodbeck/gnp/issues/5)).

This one outgrew an erratum entry: the [ch12-2026](ch12-2026) directory
contains a fully revised edition of the chapter with updated code and
text, including a table mapping each affected printed listing to its
replacement. The [ch12](ch12) directory continues to match the printed
listings.

## Chapter 13: Logging and Metrics

### Page 303: Listing 13-7's example output embeds a Go version

The zap logger in this example includes a `version` field populated from
`runtime.Version()`, so the listing's `// Output:` comments (printed with
`go1.15.5`) match only that exact Go release, and the example test fails
as printed on any other. Expect your installed Go version in the output.
The repository's [ch13/zap_test.go](ch13/zap_test.go) tracks the Go
version used to test this repository, so its `Output` comments
intentionally differ from the printed listing and change with toolchain
upgrades.

## Chapter 14: Moving to the Cloud

### Page 345: Cloud Functions no longer caps you at Go 1.13

The Google Cloud Functions discussion says deployed code must avoid
language features added after Go 1.13 because of the platform's supported
runtime. The constraint is real, but the version was a snapshot of the
platform in 2020. Consult the Go runtimes Cloud Functions currently
supports. They have long since moved past Go 1.13.
