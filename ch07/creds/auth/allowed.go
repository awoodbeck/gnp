// +build !linux

package auth

import "net"

func Allowed(conn *net.UnixConn, groups map[string]struct{}) bool {
	return true
}
