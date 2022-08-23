# Network Programming with Go
Code repository for [the book](https://nostarch.com/networkprogrammingwithgo)


[![Build Status](https://travis-ci.org/awoodbeck/gnp.svg?branch=master)](https://travis-ci.org/awoodbeck/gnp)

## How to run the examples:
Unlike a book to learn programming, this book discusses network concepts and provides code snippets that utilize the testing functionality of Go. If you have a good knowledge of Go and are following along, you can figure out the code or write your own.

If you are still learning Go, here are the steps you can take to run the examples from this repo:

Checkout the code:

    git clone git@github.com:awoodbeck/gnp.git

Run all the tests:

    cd gnp
    go test -timeout 300s -race -bench=. ./...
    
Alternatively, run the tests from one chapter:

    go test -v -timeout 300s -race -bench=. ./ch03/dial_fanout_test.go
