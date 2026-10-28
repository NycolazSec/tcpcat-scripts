// A tcpcat dissector: recognises memcached and reads its version.
//
// `version` is an unauthenticated, read-only memcached command: the server
// answers "VERSION <x.y.z>". An exposed memcached (no auth, often reachable
// from the internet) is a classic misconfiguration, so this also flags it.
//
// Build:  ./build.sh memcached     ->  memcached.wasm
// Run:    tcpcat -sT -p 11211 --scripts . <target>
package main

import "strings"

func Dissect(ip string, port int) map[string]string {
	reply := string(Probe([]byte("version\r\n")))
	if !strings.HasPrefix(reply, "VERSION ") {
		return nil // not memcached (or it needs SASL auth)
	}
	version := strings.TrimSpace(strings.TrimPrefix(reply, "VERSION "))
	return map[string]string{
		"service": "memcached",
		"version": version,
		"auth":    "none", // answered an unauthenticated command
	}
}
