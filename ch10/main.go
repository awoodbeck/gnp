package main

import (
	cmd "github.com/caddyserver/caddy/v2/cmd"
	_ "github.com/caddyserver/caddy/v2/modules/standard"

	// Injecting custom modules into Caddy
	_ "github.com/awoodbeck/caddy-restrict-prefix"
	_ "github.com/awoodbeck/caddy-toml-adapter"
)

func main() {
	cmd.Main()
}
