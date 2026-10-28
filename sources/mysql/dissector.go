// A tcpcat dissector: recognises MySQL / MariaDB and reads its version.
//
// MySQL-family servers send their handshake packet the moment you connect,
// before any request -- so Probe(nil) (a banner read) is all it takes, with
// no authentication. The server version sits in cleartext in that packet.
//
// Build:  ./build.sh mysql     ->  mysql.wasm
// Run:    tcpcat -sT -p 3306 --scripts . <target>
package main

import "strings"

func Dissect(ip string, port int) map[string]string {
	// Handshake v10: [3]len [1]seq [1]protocol(0x0a) <NUL-terminated version> ...
	b := Probe(nil)
	if len(b) < 6 || b[4] != 0x0a {
		return nil // not a MySQL-family protocol-10 handshake
	}
	end := 5
	for end < len(b) && b[end] != 0 {
		end++
	}
	version := string(b[5:end])
	if version == "" {
		return nil
	}
	// MariaDB advertises itself inside the version string (often prefixed with
	// a "5.5.5-" compatibility marker, e.g. "5.5.5-10.11.5-MariaDB-...").
	service := "mysql"
	if strings.Contains(version, "MariaDB") {
		service = "mariadb"
	}
	return map[string]string{"service": service, "version": version}
}
