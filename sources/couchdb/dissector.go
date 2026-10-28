// A tcpcat dissector: recognises Apache CouchDB and reads its version.
//
// CouchDB answers an unauthenticated GET / with a JSON welcome document
// ({"couchdb":"Welcome","version":"3.x",...}). An exposed CouchDB is a known
// data-leak class, so confirming it (and its version) is worth a probe.
//
// Build:  ./build.sh couchdb     ->  couchdb.wasm
// Run:    tcpcat -sT -p 5984 --scripts . <target>
package main

import "strings"

func Dissect(ip string, port int) map[string]string {
	resp := string(Probe([]byte("GET / HTTP/1.0\r\nHost: localhost\r\n\r\n")))
	if !strings.Contains(resp, "couchdb") || !strings.Contains(resp, "Welcome") {
		return nil
	}
	out := map[string]string{"service": "couchdb"}
	if v := jsonString(resp, "version"); v != "" {
		out["version"] = v
	}
	return out
}

// jsonString returns the string value of the first "key":"value" pair, or "".
func jsonString(s, key string) string {
	i := strings.Index(s, "\""+key+"\"")
	if i < 0 {
		return ""
	}
	i += len(key) + 2
	for i < len(s) && (s[i] == ' ' || s[i] == ':') {
		i++
	}
	if i >= len(s) || s[i] != '"' {
		return ""
	}
	i++
	j := strings.IndexByte(s[i:], '"')
	if j < 0 {
		return ""
	}
	return s[i : i+j]
}
