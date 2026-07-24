# Chapter 12: Data Serialization (2026 Edition)

This directory contains a contemporary revision of Chapter 12 from *Network
Programming with Go* (No Starch Press, 2021). The code in the [ch12](../ch12)
directory continues to mirror the listings as they appear in print. The code
here, and the chapter text below, bring those listings up to date with
current Go and gRPC practices. If you're following along with the book in
hand, the listing numbers below match the printed listing numbers, so you can
compare any listing with its printed counterpart to see exactly what changed
and why.

## What Changed Since Publication

The gRPC code-generation API changed shortly after this book went to print,
which broke the printed listings for readers using current tooling (see
issues [#3](https://github.com/awoodbeck/gnp/issues/3) and
[#5](https://github.com/awoodbeck/gnp/issues/5), and
[grpc/grpc-go#3794](https://github.com/grpc/grpc-go/issues/3794)). The Go
project deprecated a few packages and commands in the interim as well. This
table summarizes the differences between the printed chapter and this
revision.

| Printed listing | What changed | Why |
| --- | --- | --- |
| Page 281 | `go get -u github.com/golang/protobuf/protoc-gen-go` is now `go get -tool` plus `go install tool` | `go get` stopped installing binaries in Go 1.17, and `github.com/golang/protobuf` is deprecated in favor of `google.golang.org/protobuf`. Go 1.24's `tool` directive pins the generator versions in go.mod. |
| 12-17 | The `go_package` option uses the `import/path;package_name` form | The printed form appended the package name to the import path, which named an import path that doesn't exist. The semicolon form states both accurately. |
| Pages 282, 285–286 | One `protoc` invocation generates both message and gRPC code | Running `--go_out` and `--go-grpc_out` together is the conventional invocation now that both plug-ins are stable. |
| 12-18 | `ioutil.ReadAll` is now `io.ReadAll` | The `io/ioutil` package was deprecated in Go 1.16. |
| 12-21 | `Rosie` embeds `housework.UnimplementedRobotMaidServer` and drops its `Service` method | The `RobotMaidService` struct and `RegisterRobotMaidService` function came from an experimental generator API that gRPC removed before `protoc-gen-go-grpc` v1.0. The generated code now expects you to implement the `RobotMaidServer` interface instead. |
| 12-21 | `Complete` returns `status.Errorf(codes.OutOfRange, ...)` instead of `fmt.Errorf` | A plain error crosses the wire as the generic `Unknown` status code. Returning a status error gives the client a meaningful code to act on. |
| 12-22 | `housework.RegisterRobotMaidServer(server, rosie)` replaces `RegisterRobotMaidService`, and the TLS configuration moves from `tls.NewListener` to `grpc.Creds` | Same generator change as Listing 12-21. Passing transport credentials to `grpc.NewServer` is the idiomatic way to serve TLS in gRPC, and `PreferServerCipherSuites` has been a no-op since Go 1.18. |
| 12-27 | `grpc.NewClient` replaces `grpc.Dial`, and `os.ReadFile` replaces `ioutil.ReadFile` | gRPC deprecated `Dial` in v1.63. `NewClient` returns immediately and connects lazily when you make your first call. |

The chapter follows, revised accordingly.

---

A sizable portion of our work as developers involves integrating our network
services with existing ones, some of which may be legacy services or services
written in languages other than Go. For these services to work together, they
need to exchange bytes in a format that both sides understand, regardless of
the language on either end of the connection. The sender converts structured
data into a series of bytes, sends those bytes over the network, and trusts
the receiver to convert them back into structured data. Transforming an
object into a series of bytes this way is called *data serialization*.
Serialized data is equally suitable for transmission over a network or
persistence to disk. Any code that understands the serialization format can
reconstruct the original object from the bytes alone.

You've serialized data in this book already, perhaps without labeling it as
such. The type-length-value encoding from Chapter 4, the JSON you posted
over HTTP in Chapter 8, and the PEM-encoded certificates you generated in
Chapter 11 are all data serialization formats at work.

This chapter takes a deeper look at data serialization for storage and for
exchange with other systems, particularly systems you didn't write. We'll
concentrate on the three formats you're most likely to reach for in Go
network programming: JSON, protocol buffers, and Gob. We'll then use protocol
buffers as the foundation for gRPC, a framework for calling code on remote
machines as if it were local. By the end of this chapter, you will know how
to serialize objects for storage or transmission, deserialize them back into
useful data structures, and stand up a gRPC service that other systems can
call over a network.

## Serializing Objects

Objects cannot traverse a network connection as is. `net.Conn`'s `Write`
method accepts a byte slice, not a struct, so you must first serialize the
object into bytes. Thankfully, Go makes this easy. The standard library's
`encoding` packages cover the popular formats: you've used
`encoding/binary` to serialize numbers, `encoding/json` to serialize objects
for HTTP, and `encoding/pem` to write certificates to disk. (When a function
or method name includes *encode* or *marshal*, it almost certainly
serializes data. Likewise, *decode* and *unmarshal* mean deserialization.)

In this section, you'll write an application that serializes its data in
three formats: JSON, Gob, and protocol buffers. I have a habit of promising
my wife I'll get to a household chore and then promptly losing track of it,
so the application manages a list of household chores, persisting the list
to disk between runs. Listing 12-1 starts the application with a type
representing a single chore.

```go
package housework

type Chore struct {
	Complete    bool
	Description string
}
```

*Listing 12-1: A type to represent a household chore
([housework/housework.go](housework/housework.go))*

The `Chore` struct has a field indicating whether the chore is complete and
a description of the work itself. Both fields are exported, which matters
here: the serialization packages you're about to use can access only
exported fields.

The application needs to do three things from the command line: list the
chores, add new chores, and mark chores complete. Listing 12-2 lays the
groundwork.

```go
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/awoodbeck/gnp/ch12-2026/housework"
	❶ storage "github.com/awoodbeck/gnp/ch12-2026/json"
	// storage "github.com/awoodbeck/gnp/ch12-2026/gob"
	// storage "github.com/awoodbeck/gnp/ch12-2026/protobuf"
)

var dataFile string

func init() {
	flag.StringVar(&dataFile, "file", "housework.db", "data file")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(),
			`Usage: %s [flags] [add chore, ...|complete #]
    add         add comma-separated chores
    complete    complete designated chore

Flags:
`, filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}
}
```

*Listing 12-2: Initial housework application code
([cmd/housework.go](cmd/housework.go))*

You'll write one storage package per serialization format, each exposing the
same two functions, so the application can switch formats by swapping a
single import ❶. The import is aliased to `storage`, which keeps the rest of
the code oblivious to the format in use. You start with JSON and leave the
Gob and protocol buffer imports commented out for later in this chapter.
The application accepts an optional flag designating the file that backs
the chore list, defaulting to *housework.db* in the current directory.

Now let's write the code to load chores from storage (Listing 12-3).

```go
func load() ([]*housework.Chore, error) {
	❶ if _, err := os.Stat(dataFile); os.IsNotExist(err) {
		return make([]*housework.Chore, 0), nil
	}

	df, err := os.Open(dataFile)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := df.Close(); err != nil {
			fmt.Printf("closing data file: %v", err)
		}
	}()

	❷ return storage.Load(df)
}
```

*Listing 12-3: Deserializing chores from a file
([cmd/housework.go](cmd/housework.go))*

If the data file doesn't exist yet ❶, you return an empty chore list rather
than an error, since a missing file simply means you haven't procrastinated
anything yet. Otherwise, you open the file and pass it to the storage
package's `Load` function ❷, which accepts an `io.Reader`. Accepting
interfaces rather than concrete file types keeps the storage packages
flexible and easy to test.

Listing 12-4 defines the complementary function, which flushes the in-memory
chores back to disk.

```go
func flush(chores []*housework.Chore) error {
	df, err := os.Create(dataFile)
	if err != nil {
		return err
	}
	defer func() {
		if err := df.Close(); err != nil {
			fmt.Printf("closing data file: %v", err)
		}
	}()

	return storage.Flush(df, chores)
}
```

*Listing 12-4: Flushing chores to storage
([cmd/housework.go](cmd/housework.go))*

Here you create or truncate the data file and hand it to the storage
package's `Flush` function, which serializes the chores and writes them to
its `io.Writer`. This approach rewrites the entire file on every change,
which is perfectly reasonable for a chore list. If your chores multiply to
the point that rewriting the file becomes a bottleneck, you have bigger
problems than serialization overhead.

You need a way to see the list on the command line. Listing 12-5 writes the
chores to standard output.

```go
func list() error {
	chores, err := load()
	if err != nil {
		return err
	}

	if len(chores) == 0 {
		fmt.Println("You're all caught up!")
		return nil
	}

	fmt.Println("#\t[X]\tDescription")
	for i, chore := range chores {
		c := " "
		if chore.Complete {
			c = "X"
		}
		fmt.Printf("%d\t[%s]\t%s\n", i+1, c, chore.Description)
	}

	return nil
}
```

*Listing 12-5: Printing the list of chores to standard output
([cmd/housework.go](cmd/housework.go))*

The function loads the chores and, if the list is empty, congratulates you.
Otherwise, it prints a header followed by one line per chore: the chore's
number, an `X` if the chore is complete, and its description. Since Go
slices are zero-indexed and humans generally aren't, you number the chores
starting at 1. Listing 12-6 shows what the output looks like with a few
chores on the list.

```
#       [X]     Description
1       [ ]     Mop floors
2       [X]     Wash dishes
3       [ ]     Vacuum
```

*Listing 12-6: Example output of the `list` function with three chores in
the list*

Listing 12-7 implements the `add` function, which accepts new chores as a
comma-separated string.

```go
func add(s string) error {
	chores, err := load()
	if err != nil {
		return err
	}

	❶ for _, chore := range strings.Split(s, ",") {
		if desc := strings.TrimSpace(chore); desc != "" {
			chores = append(chores, &housework.Chore{
				Description: desc,
			})
		}
	}

	return flush(chores)
}
```

*Listing 12-7: Adding chores to the list
([cmd/housework.go](cmd/housework.go))*

You split the input on commas ❶ so a single command can add several chores
at once, trimming whitespace from each description and skipping any empty
entries. This does mean chore descriptions cannot contain commas.
I'm confident you can devise a delimiter-free approach once
you've finished this chapter, and mowing the lawn will still be there when
you're done.

Since I've yet to convince my wife that deleting a chore is the same as
completing it, Listing 12-8 settles for marking chores complete.

```go
func complete(s string) error {
	i, err := strconv.Atoi(s)
	if err != nil {
		return err
	}

	chores, err := load()
	if err != nil {
		return err
	}

	if i < 1 || i > len(chores) {
		return fmt.Errorf("chore %d not found", i)
	}

	❶ chores[i-1].Complete = true

	return flush(chores)
}
```

*Listing 12-8: Marking a chore as complete
([cmd/housework.go](cmd/housework.go))*

The function converts its argument to an integer, makes sure the number
falls within the bounds of the list, and marks the corresponding chore
complete, accounting for the 1-based numbering the `list` function
presented ❶. It then flushes the modified list to disk.

Listing 12-9 ties the functions to the command line in `main`.

```go
func main() {
	flag.Parse()

	var err error

	switch strings.ToLower(flag.Arg(0)) {
	case "add":
		err = add(strings.Join(flag.Args()[1:], " "))
	case "complete":
		err = complete(flag.Arg(1))
	}

	if err != nil {
		log.Fatal(err)
	}

	err = list()
	if err != nil {
		log.Fatal(err)
	}
}
```

*Listing 12-9: The main logic of the housework application
([cmd/housework.go](cmd/housework.go))*

If the first argument is `add`, you join the remaining arguments into a
single string for the `add` function to split on commas. If it's
`complete`, you pass the second argument along to mark that chore done.
In every case, the application finishes by printing the current list, so
you always see the effect of your command.

All that remains is a storage package. Let's write the JSON one first.

### JSON

JSON, or JavaScript Object Notation, may be the most widely supported
text-based serialization format in use today. It represents data using
key-value pairs and arrays built from strings, numbers, Booleans, and null,
and nearly every contemporary language offers first-class support for it —
one reason it's the lingua franca of RESTful APIs. You worked with Go's
`encoding/json` package in Chapter 8. Here, you'll use it to back the
housework application. Listing 12-10 shows the contents of the
*housework.db* file that the JSON storage implementation produces, formatted
for readability. (The encoder writes it as a single line. You could use the
encoder's `SetIndent` method to make the file itself look like this.)

```json
[
  {
    "Complete": false,
    "Description": "Mop floors"
  },
  {
    "Complete": true,
    "Description": "Wash dishes"
  },
  {
    "Complete": false,
    "Description": "Vacuum"
  }
]
```

*Listing 12-10: Formatted contents of the housework.db file after
serializing the chores to JSON*

Notice that the field names match the `Chore` struct exactly. If you wanted
different names in the serialized output — lowercase keys, say — you'd add
struct tags to Listing 12-1. Since this file is yours alone, the exported
field names work fine as is. Listing 12-11 details the entire JSON storage
implementation.

```go
package json

import (
	"encoding/json"
	"io"

	"github.com/awoodbeck/gnp/ch12-2026/housework"
)

func Load(r io.Reader) ([]*housework.Chore, error) {
	var chores []*housework.Chore

	return chores, json.NewDecoder(r).Decode(&chores)
}

func Flush(w io.Writer, chores []*housework.Chore) error {
	return json.NewEncoder(w).Encode(chores)
}
```

*Listing 12-11: JSON storage implementation
([json/housework.go](json/housework.go))*

`Load` wraps its `io.Reader` in a JSON decoder and decodes the contents into
a slice of chores, and `Flush` does the reverse with an encoder. That's the
entire package. Let's put it through its paces on the command line
(Listing 12-12).

```
$ go run cmd/housework.go
You're all caught up!
$ go run cmd/housework.go add Mop floors, Wash dishes, Vacuum
#       [X]     Description
1       [ ]     Mop floors
2       [ ]     Wash dishes
3       [ ]     Vacuum
$ go run cmd/housework.go complete 2
#       [X]     Description
1       [ ]     Mop floors
2       [X]     Wash dishes
3       [ ]     Vacuum
```

*Listing 12-12: Testing the housework application with JSON storage on the
command line*

You add three chores in one command, then mark the second one complete.
So far, the housework application is working better than my own follow-through.

### Gob

Gob is Go's native binary serialization format — the name is short for
"gobs of binary data." It's efficient, it's self-describing, and it was
designed by the Go team to take full advantage of Go's type system, which
makes it the path of least resistance when Go programs need to exchange
data with each other. The trade-off is exactly that Go-centric design: if a
Python service needs to read your data, JSON or protocol buffers will serve
you better. Listing 12-13 implements the Gob storage backend.

```go
package gob

import (
	"encoding/gob"
	"io"

	"github.com/awoodbeck/gnp/ch12-2026/housework"
)

func Load(r io.Reader) ([]*housework.Chore, error) {
	var chores []*housework.Chore

	return chores, gob.NewDecoder(r).Decode(&chores)
}

func Flush(w io.Writer, chores []*housework.Chore) error {
	return gob.NewEncoder(w).Encode(chores)
}
```

*Listing 12-13: Gob storage implementation
([gob/housework.go](gob/housework.go))*

If your first thought is that I copied Listing 12-11 and replaced `json`
with `gob`, you've caught me. In Go, Gob really is a drop-in replacement
for JSON, which is precisely why the housework application can treat
storage as a swappable import. Listing 12-14 makes the swap.

```go
	"github.com/awoodbeck/gnp/ch12-2026/housework"
	❶ // storage "github.com/awoodbeck/gnp/ch12-2026/json"
	❷ storage "github.com/awoodbeck/gnp/ch12-2026/gob"
	// storage "github.com/awoodbeck/gnp/ch12-2026/protobuf"
```

*Listing 12-14: Swapping the JSON storage package for the Gob storage
package ([cmd/housework.go](cmd/housework.go))*

Comment out the JSON storage import ❶ and uncomment the Gob one ❷. Since
the previous *housework.db* is JSON, delete it before testing the
application with its new storage format (Listing 12-15).

```
$ rm housework.db
$ go run cmd/housework.go add Mop floors, Wash dishes
#       [X]     Description
1       [ ]     Mop floors
2       [ ]     Wash dishes
$ go run cmd/housework.go complete 1
#       [X]     Description
1       [X]     Mop floors
2       [ ]     Wash dishes
$ xxd housework.db | head -n 3
00000000: 0dff 8102 0102 ff82 0001 ff80 0000 287f  ..............(.
00000010: 0301 02ff 8000 0102 0108 436f 6d70 6c65  ..........Comple
00000020: 7465 0102 0001 0b44 6573 6372 6970 7469  te.....Descripti
```

*Listing 12-15: Testing the housework application with Gob storage on the
command line*

The application works exactly as before, but the data file is no longer
something you'd want to edit by hand. You can make out the field names in
the hex dump — Gob streams describe their own types — but the rest is
binary. Go doesn't mind one bit, and Gob's encoding is more compact and
quicker to process than the equivalent JSON.

### Protocol Buffers

Protocol buffers are a binary serialization format designed by Google for
exchanging data across platforms and languages. Like JSON, protocol buffers
are language neutral. Like Gob, they're compact and fast. Unlike either,
they're *schema first*: you define your messages in a *.proto* file using
the protocol buffer language, then compile that file into code for each
language that needs to read or write the data. We'll use the proto3 syntax
in this chapter. A protocol buffer definition for our `Chore` type looks
like Listing 12-16.

```proto
message Chore {
  bool complete = 1;
  string description = 2;
}
```

*Listing 12-16: Protocol buffer message definition representing a Chore
type*

You define a *message* with the `message` keyword and a Pascal-cased name.
Each field has a type, a snake-cased name, and a field number that is unique
within the message. The name is for you. The number is what identifies the
field in the binary payload, which means field numbers must never change
once in use, or you'll break backward compatibility with every serialized
payload that came before. Adding new fields and new messages is always safe.

Because third parties may generate code from your *.proto* file to talk to
your service, treat the file like the API it is. One useful convention is
to version the definitions themselves, so you can introduce breaking
changes in a *v2* package while existing clients continue using *v1*. You'll
see this convention in the file path shortly.

You'll need two things to generate Go code from a *.proto* file: the
protocol buffer compiler, `protoc`, and its Go generator plug-ins. Your
operating system's package manager likely carries the compiler. On Debian
and Ubuntu, run the following:

```
$ sudo apt install protobuf-compiler
```

On macOS or Linux with Homebrew, run this:

```
$ brew install protobuf
```

On Windows, download the latest protocol buffer compiler ZIP file from
[https://github.com/protocolbuffers/protobuf/releases/](https://github.com/protocolbuffers/protobuf/releases/),
extract it, and add its *bin* subdirectory to your `PATH`. You should now be
able to run `protoc` on your command line.

The Go plug-ins are Go programs, and as of Go 1.24 you can track them in
your *go.mod* file with the `tool` directive, the same way you track any
other dependency. This pins the exact generator versions for everyone who
clones your module — a lesson this chapter earned the hard way, as you'll
see in the next section. Add the tools to your module like so:

```
$ go get -tool google.golang.org/protobuf/cmd/protoc-gen-go
$ go get -tool google.golang.org/grpc/cmd/protoc-gen-go-grpc
```

This module's *go.mod* already includes both directives, so you only need
to install them:

```
$ go install tool
```

The `go install tool` command builds every tool listed in *go.mod* and
installs the binaries in your *$GOPATH/bin* directory. Make sure that
directory is in your `PATH` environment variable, since `protoc` discovers
its plug-ins there.

> **NOTE** *You may see older documentation instruct you to run
> `go get -u github.com/golang/protobuf/protoc-gen-go`, as the first
> edition of this book did. Neither half of that command survives: `go get`
> stopped installing binaries in Go 1.17, and the
> `github.com/golang/protobuf` module is deprecated in favor of
> `google.golang.org/protobuf`.*

Now that you've installed the compiler and its plug-ins, let's write the
full protocol buffer definition for the housework application
(Listing 12-17). This file lives at *housework/v1/housework.proto*, where
*v1* is the versioning convention mentioned earlier.

```proto
syntax = "proto3";
❶ package housework;

❷ option go_package = "github.com/awoodbeck/gnp/ch12-2026/housework/v1;housework";

message Chore {
  bool complete = 1;
  string description = 2;
}

message Chores {
  ❸ repeated Chore chores = 1;
}
```

*Listing 12-17: Protocol buffer definition for your housework application
([housework/v1/housework.proto](housework/v1/housework.proto))*

You declare the proto3 syntax and a package name ❶ for the generated code.
The `go_package` option ❷ tells the Go generator two things, separated by a
semicolon: the full import path of the generated package and the package
name to use in the generated code. You then define the `Chore` message and
a `Chores` message that holds any number of `Chore` messages, courtesy of
the `repeated` keyword ❸.

Now compile the *.proto* file to Go code:

```
$ protoc --go_out=. --go_opt=paths=source_relative housework/v1/housework.proto
```

The `--go_out` flag tells `protoc` to invoke the `protoc-gen-go` plug-in
and where to write its output, while `paths=source_relative` places the
generated file next to its source file. If `protoc` complains that it
cannot find `protoc-gen-go`, revisit your `PATH`. When the compilation
succeeds, you'll find a new *housework/v1/housework.pb.go* file whose first
few lines look like this, allowing for version drift:

```
$ head -n 7 housework/v1/housework.pb.go
// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
//      protoc-gen-go v1.36.11
//      protoc        v7.35.1
// source: housework/v1/housework.proto

package housework
```

As the comment insists, don't edit this file. When you need a change, make
it in the *.proto* file and recompile. Listing 12-18 puts the generated
code to work in the protocol buffer storage backend.

```go
package protobuf

import (
	"io"

	❶ "google.golang.org/protobuf/proto"

	❷ "github.com/awoodbeck/gnp/ch12-2026/housework/v1"
)

func Load(r io.Reader) ([]*housework.Chore, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	var chores housework.Chores

	return chores.Chores, ❸ proto.Unmarshal(b, &chores)
}

func Flush(w io.Writer, chores []*housework.Chore) error {
	b, err := ❹ proto.Marshal(&housework.Chores{Chores: chores})
	if err != nil {
		return err
	}

	_, err = w.Write(b)

	return err
}
```

*Listing 12-18: Protocol buffers storage implementation
([protobuf/housework.go](protobuf/housework.go))*

Two imports distinguish this backend from its siblings. The `proto`
package ❶ supplies the `Marshal` and `Unmarshal` functions, and the
generated *housework/v1* package ❷ replaces the hand-written `housework`
package from Listing 12-1, since the storage functions now traffic in
generated types. Unlike `encoding/json` and `encoding/gob`, the `proto`
package doesn't wrap an `io.Reader` or `io.Writer`, so `Load` reads the
full payload before unmarshaling it ❸ and `Flush` marshals to a byte
slice ❹ before writing.

Revisit Listing 12-2 and plug in the protocol buffer implementation with
the two changes shown in Listing 12-19.

```go
	❶ "github.com/awoodbeck/gnp/ch12-2026/housework/v1"
	// storage "github.com/awoodbeck/gnp/ch12-2026/json"
	// storage "github.com/awoodbeck/gnp/ch12-2026/gob"
	❷ storage "github.com/awoodbeck/gnp/ch12-2026/protobuf"
```

*Listing 12-19: Swapping the JSON storage package for the protobuf storage
package ([cmd/housework.go](cmd/housework.go))*

This swap involves one more change than the last one: in addition to
switching the storage import ❷, you replace the `housework` package with
the generated *housework/v1* package ❶, since the rest of the application
must now pass around the generated `Chore` type. The generated package
declares itself `housework`, so nothing else in the file changes. Delete
the Gob-formatted *housework.db* and the application works as it did in
Listings 12-12 and 12-15, now speaking a format that services in dozens of
other languages can consume with nothing more than your *.proto* file.

## Transmitting Serialized Objects

Serializing objects to disk is useful, but the more common reason to
serialize data is to send it to another service. Picture an online store:
the storefront needs to consult inventory, billing, shipping, and
notification services to do its job. If every service ran on one server,
scaling the store would mean buying an ever-larger server. Running each
service on its own hardware scales better, but now the services can't share
memory. They need to exchange data over the network — and ideally in a way
that doesn't make every function call feel like a networking exercise.

*Remote procedure calls* (RPCs) are that abstraction. With RPC, a client
calls what appears to be a local function, and machinery beneath it relays
the request to a server, executes it there, and returns the result as the
function's return value. Your code reads the same whether the work happens
in-process or on another continent, which lets you distribute services
across servers without distributing that complexity throughout your code.

The dominant RPC framework today is *gRPC*, a cross-platform,
language-agnostic framework built on HTTP/2 and protocol buffers. Let's use
it to build something more ambitious than a to-do list: a service that
delegates your chores to Rosie, the robotic maid from *The Jetsons*. She
isn't due for a few more decades, but the code will be ready for her.

### Connecting Services with gRPC

The gRPC framework hides most of its wire-level details behind generated
code, and you already have the toolchain for generating it: the
`protoc-gen-go-grpc` plug-in you installed with `go install tool` earlier
in this chapter. Everything starts with the *.proto* file. Listing 12-20
adds a service definition and the supporting messages to Listing 12-17.

```proto
--snip--

service RobotMaid {
  ❶ rpc Add (Chores) returns (Response);
  rpc Complete (CompleteRequest) returns (Response);
  rpc List (Empty) returns (Chores);
}

message CompleteRequest {
  ❷ int32 chore_number = 1;
}

❸ message Empty {}

message Response {
  string message = 1;
}
```

*Listing 12-20: Additional definitions to support a gRPC RobotMaid service
([housework/v1/housework.proto](housework/v1/housework.proto))*

The service mirrors the three commands from the command line application.
You declare a service named `RobotMaid` containing three RPC methods, each
introduced by the `rpc` keyword with its request message type in
parentheses and its response message type after the `returns` keyword ❶.
Rather than call the `add`, `complete`, and `list` functions locally,
you'll call these methods on the server.

Until Rosie ships with proper artificial intelligence, you'll need to tell
her when a chore is done, so the `CompleteRequest` message carries the
number of the completed chore ❷, and a `Response` message carries her
acknowledgment. The `List` method needs no input at all, but every RPC
method requires a request message. gRPC's equivalent of `nil` is a message
with no fields ❸.

> **NOTE** *The `google.protobuf.Empty` well-known type serves this purpose
> in most contemporary .proto files, saving you the three keystrokes of
> defining your own. We define ours here to keep the file self-contained.*

Recompile the *.proto* file, this time asking for the gRPC code as well.
The two plug-ins happily run in a single invocation:

```
$ protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    housework/v1/housework.proto
```

The `--go-grpc_out` flag invokes the `protoc-gen-go-grpc` plug-in, which
writes the service code to *housework/v1/housework_grpc.pb.go* beside the
message code. You now have everything you need to build a gRPC server and
client.

### Creating a TLS-Enabled gRPC Server

gRPC assumes a secure connection by default, so let's give it one. You'll
reuse the *cert.pem* and *key.pem* files you created in Chapter 11 for the
server, and the client will pin the server's certificate, as in the
"Generating Certificates for Authentication" section of that chapter.
Copies of both files live in this directory's *server* subdirectory, valid
into 2030. Regenerate them per Chapter 11 if you're reading this from the
far future, and say hello to Rosie for me.

On the server side, your job is to implement the `RobotMaidServer`
interface that `protoc-gen-go-grpc` generated from the service definition:

```go
type RobotMaidServer interface {
	Add(context.Context, *Chores) (*Response, error)
	Complete(context.Context, *CompleteRequest) (*Response, error)
	List(context.Context, *Empty) (*Chores, error)
	mustEmbedUnimplementedRobotMaidServer()
}
```

That unexported method at the bottom deserves a moment of attention,
because you cannot write it yourself. The generated code includes an
`UnimplementedRobotMaidServer` struct that satisfies it, and by embedding
that struct in your implementation, you satisfy it too. The embedded struct
also provides a default implementation of every RPC method, each returning
an `Unimplemented` status code. That's forward compatibility by
construction: when you add a new RPC method to the *.proto* file next year,
existing implementations continue to compile, responding to the new method
with a sensible error instead of a build failure.

> **NOTE** *The first printing of this book targeted a pre-release version
> of `protoc-gen-go-grpc` that generated a `RobotMaidService` struct and a
> `RegisterRobotMaidService` function. gRPC removed that API before the
> plug-in's v1.0 release, and the plug-in has generated the interface you
> see above ever since. This is why the generator versions now live in
> go.mod, where an upgrade is a code review instead of a surprise.*

Listing 12-21 implements the interface with a new type named `Rosie`.

```go
package main

import (
	"context"
	"sync"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/awoodbeck/gnp/ch12-2026/housework/v1"
)

type Rosie struct {
	❶ housework.UnimplementedRobotMaidServer
	mu     sync.Mutex
	❷ chores []*housework.Chore
}

func (r *Rosie) Add(_ context.Context, chores *housework.Chores) (
	*housework.Response, error) {
	r.mu.Lock()
	r.chores = append(r.chores, chores.Chores...)
	r.mu.Unlock()

	return ❸ &housework.Response{Message: "ok"}, nil
}

func (r *Rosie) Complete(_ context.Context,
	req *housework.CompleteRequest) (*housework.Response, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.chores == nil || req.ChoreNumber < 1 ||
		int(req.ChoreNumber) > len(r.chores) {
		return nil, ❹ status.Errorf(codes.OutOfRange, "chore %d not found",
			req.ChoreNumber)
	}

	r.chores[req.ChoreNumber-1].Complete = true

	return &housework.Response{Message: "ok"}, nil
}

func (r *Rosie) List(_ context.Context, _ *housework.Empty) (
	*housework.Chores, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.chores == nil {
		r.chores = make([]*housework.Chore, 0)
	}

	return &housework.Chores{Chores: r.chores}, nil
}
```

*Listing 12-21: Building a RobotMaid-compatible type named Rosie
([server/rosie.go](server/rosie.go))*

`Rosie` embeds `UnimplementedRobotMaidServer` ❶ to satisfy the generated
interface, then overrides all three RPC methods with real implementations.
She keeps her chore list in memory ❷, guarded by a mutex, since multiple
clients may call her concurrently. The `Add` and `List` methods return a
response message or an error ❸, either of which makes its way back to the
client.

The `Complete` method returns its error using the `status` package ❹
rather than `fmt.Errorf`. Any plain error you return crosses the wire as
the generic `Unknown` status code. Wrapping the message in a status error
assigns it a meaningful code — `OutOfRange` here — that clients can inspect
with `status.Code` and act on. Reserve plain errors for problems you
haven't classified yet.

Now, let's stand up a gRPC server around `Rosie` (Listing 12-22).

```go
package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/awoodbeck/gnp/ch12-2026/housework/v1"
)

var addr, certFn, keyFn string

func init() {
	flag.StringVar(&addr, "address", "localhost:34443",
		"listen address")
	flag.StringVar(&certFn, "cert", "cert.pem", "certificate file")
	flag.StringVar(&keyFn, "key", "key.pem", "private key file")
}

func main() {
	flag.Parse()

	cert, err := tls.LoadX509KeyPair(certFn, keyFn)
	if err != nil {
		log.Fatal(err)
	}

	server := ❶ grpc.NewServer(
		❷ grpc.Creds(credentials.NewTLS(
			&tls.Config{
				Certificates:     []tls.Certificate{cert},
				CurvePreferences: []tls.CurveID{tls.CurveP256},
				MinVersion:       tls.VersionTLS12,
			},
		)),
	)
	rosie := new(Rosie)
	❸ housework.RegisterRobotMaidServer(server, rosie)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Listening for TLS connections on %s ...", addr)
	❹ log.Fatal(server.Serve(listener))
}
```

*Listing 12-22: Creating a new gRPC server using Rosie
([server/server.go](server/server.go))*

You load the server's key pair, then create the gRPC server ❶, handing it
the TLS configuration through the `grpc.Creds` server option ❷. Letting
gRPC manage the TLS handshake this way, rather than wrapping the listener
with `tls.NewListener` yourself, keeps the transport details where gRPC
expects them and works with the framework's connection management. You then
register a new `Rosie` instance with the server ❸ using the generated
`RegisterRobotMaidServer` function, bind a TCP listener on port 34443, and
put the server to work ❹.

### Creating a gRPC Client to Test the Server

The client-side code should feel familiar: it's the command line
application from Listing 12-2, with the storage-backed functions replaced
by calls to a gRPC client whose methods correspond to the RPC methods you
defined in Listing 12-20. Listing 12-23 shows the client application's
scaffolding.

```go
package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/awoodbeck/gnp/ch12-2026/housework/v1"
)

var addr, caCertFn string

func init() {
	flag.StringVar(&addr, "address", "localhost:34443",
		"server address")
	flag.StringVar(&caCertFn, "ca-cert", "cert.pem", "CA certificate")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(),
			`Usage: %s [flags] [add chore, ...|complete #]
    add         add comma-separated chores
    complete    complete designated chore

Flags:
`, filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}
}
```

*Listing 12-23: Initial gRPC client code for our housework application
([client/client.go](client/client.go))*

In place of a data file flag, the client accepts the server's address and
the path to the certificate it should trust. Listing 12-24 uses the gRPC
client to list the current chores.

```go
func list(ctx context.Context,
	❶ client housework.RobotMaidClient) error {
	chores, err := ❷ client.List(ctx, new(housework.Empty))
	if err != nil {
		return err
	}

	if len(chores.Chores) == 0 {
		fmt.Println("You have nothing to do!")
		return nil
	}

	fmt.Println("#\t[X]\tDescription")
	for i, chore := range chores.Chores {
		c := " "
		if chore.Complete {
			c = "X"
		}
		fmt.Printf("%d\t[%s]\t%s\n", i+1, c, chore.Description)
	}

	return nil
}
```

*Listing 12-24: Using the gRPC client to list the current chores
([client/client.go](client/client.go))*

This is Listing 12-5 with the storage swapped for a `RobotMaidClient` ❶,
the client-side counterpart that `protoc-gen-go-grpc` generated alongside
the server interface. Asking the server for the chore list is a method
call ❷ that accepts a context and an `Empty` message. The fact that the
answer arrives over a TLS connection is gRPC's business, not yours.

Listing 12-25 sends new chores to the server.

```go
func add(ctx context.Context, client housework.RobotMaidClient,
	s string) error {
	chores := new(housework.Chores)

	for _, chore := range strings.Split(s, ",") {
		if desc := strings.TrimSpace(chore); desc != "" {
			chores.Chores = append(chores.Chores, &housework.Chore{
				Description: desc,
			})
		}
	}

	_, err := ❶ client.Add(ctx, chores)

	return err
}
```

*Listing 12-25: Adding new chores using the gRPC client
([client/client.go](client/client.go))*

As in Listing 12-7, you split the comma-separated input into individual
chores, but rather than flushing them to disk, you gather them into a
`Chores` message and hand them to Rosie ❶. In Listing 12-26, you write the
code to mark chores complete, which takes a bit less code than Listing 12-8
now that the server does the bookkeeping.

```go
func complete(ctx context.Context, client housework.RobotMaidClient,
	s string) error {
	i, err := strconv.Atoi(s)
	if err == nil {
		_, err = client.Complete(ctx,
			&housework.CompleteRequest{❶ ChoreNumber: int32(i)})
	}

	return err
}
```

*Listing 12-26: Marking chores complete by using the gRPC client
([client/client.go](client/client.go))*

Notice that `protoc-gen-go` converted the snake-cased `chore_number` field
from Listing 12-20 into the Pascal-cased `ChoreNumber` ❶ in the generated
Go code. You also convert the `int` from `strconv.Atoi` to an `int32` to
match the field's type.

Listing 12-27 creates the client's connection to the server, pinning the
server's certificate to its TLS configuration.

```go
--snip--

func main() {
	flag.Parse()

	caCert, err := os.ReadFile(caCertFn)
	if err != nil {
		log.Fatal(err)
	}
	certPool := x509.NewCertPool()
	if ok := certPool.AppendCertsFromPEM(caCert); !ok {
		log.Fatal("failed to add certificate to pool")
	}

	conn, err := ❶ grpc.NewClient(addr, ❷ grpc.WithTransportCredentials(
		❸ credentials.NewTLS(
			&tls.Config{
				CurvePreferences: []tls.CurveID{tls.CurveP256},
				MinVersion:       tls.VersionTLS12,
				RootCAs:          certPool,
			},
		),
	))
	if err != nil {
		log.Fatal(err)
	}
```

*Listing 12-27: Creating a new gRPC connection using TLS and certificate
pinning ([client/client.go](client/client.go))*

You read the server's certificate into a new certificate pool and build a
TLS configuration that trusts only that pool, exactly as you pinned
certificates in Chapter 11. The `grpc.NewClient` function ❶ accepts the
server's address and any options, and you pass your TLS configuration
through the `grpc.WithTransportCredentials` option ❷ after wrapping it
with `credentials.NewTLS` ❸.

Despite appearances, `grpc.NewClient` doesn't touch the network. It returns
immediately, and the connection is established lazily when you make your
first call, with gRPC transparently reconnecting as needed thereafter.
(Its predecessor, `grpc.Dial`, behaved the same way by default, which
surprised enough people expecting a function named *Dial* to dial that its
replacement got a more honest name. gRPC deprecated `Dial` in v1.63.)

You use this connection to instantiate a new gRPC client in Listing 12-28.

```go
--snip--

	rosie := ❶ housework.NewRobotMaidClient(conn)
	ctx := context.Background()

	switch strings.ToLower(flag.Arg(0)) {
	case "add":
		err = add(ctx, rosie, strings.Join(flag.Args()[1:], " "))
	case "complete":
		err = complete(ctx, rosie, flag.Arg(1))
	}

	if err != nil {
		log.Fatal(err)
	}

	err = list(ctx, rosie)
	if err != nil {
		log.Fatal(err)
	}
}
```

*Listing 12-28: Instantiating a new gRPC client and making calls
([client/client.go](client/client.go))*

Aside from creating the client from the connection ❶, this code differs
little from Listing 12-9. The difference lives below it: every interaction
with the chore list now transparently crosses a TLS connection to the gRPC
server.

Give it a try. In one terminal, start the server:

```
$ go run ./server -cert server/cert.pem -key server/key.pem
Listening for TLS connections on localhost:34443 ...
```

And in another terminal, run the client:

```
$ go run ./client -ca-cert server/cert.pem
You have nothing to do!
$ go run ./client -ca-cert server/cert.pem add Mop floors, Wash dishes
#       [X]     Description
1       [ ]     Mop floors
2       [ ]     Wash dishes
$ go run ./client -ca-cert server/cert.pem complete 2
#       [X]     Description
1       [ ]     Mop floors
2       [X]     Wash dishes
$ go run ./client -ca-cert server/cert.pem complete 9
2026/07/23 17:36:13 rpc error: code = OutOfRange desc = chore 9 not found
exit status 1
```

*Testing the RobotMaid service from the command line*

That last command shows the `status` package earning its keep: the
`OutOfRange` code you set on the server in Listing 12-21 arrives intact on
the client, where code could branch on it rather than parse an error
string. Of course, restarting the server wipes out the chore list, since
Rosie keeps it in memory. I'll leave persistence as an exercise for you —
the server could load and flush chores using any storage backend from the
first half of this chapter.

## What You've Learned

Data serialization lets you exchange structured data with systems that
share nothing with yours but a network connection and a format. It also
lets you persist that data and pick up where your application left off.

JSON is the most widely used text-based serialization format and the
default choice for RESTful APIs, with excellent support in Go's standard
library and virtually everywhere else. Gob is Go's native binary format:
faster and more compact than JSON, and effectively a drop-in replacement
for it, so long as Go programs sit on both ends of the exchange. Protocol
buffers ask more of you up front — a *.proto* file, a compiler, and
generated code — and repay the effort with compact, well-specified messages
that any supported language can produce or consume from the same
definition.

gRPC builds on protocol buffers to provide remote procedure calls: a
client invokes what looks like a local method, and gRPC carries the request
over HTTP/2 to a server that executes it and returns the result. You
define the service alongside your messages in the *.proto* file, generate
the client and server code, and implement the generated server interface —
embedding the generated `Unimplemented` struct so that tomorrow's additions
to the service don't break today's implementation. A tool as central to
your build as a code generator deserves version pinning, and Go's `tool`
directive now records it in go.mod beside your other dependencies. This
chapter's revision history is the argument for doing so.
