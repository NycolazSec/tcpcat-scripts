// A tcpcat dissector: recognises Elasticsearch and reads its version.
//
// Elasticsearch answers an unauthenticated GET / with a JSON document that
// carries the famous tagline "You Know, for Search" and a version.number
// field. An Elasticsearch reachable without auth is a classic data-exposure
// finding, so confirming it (and its version) is worth a probe.
//
// Build:  ./build.sh elasticsearch     ->  elasticsearch.wasm
// Run:    tcpcat -sT -p 9200 --scripts . <target>
package main

import "strings"

func Dissect(ip string, port int) map[string]string {
	resp := string(Probe([]byte("GET / HTTP/1.0\r\nHost: localhost\r\nUser-Agent: tcpcat\r\n\r\n")))
	if !strings.Contains(resp, "You Know, for Search") {
		return nil // the Elasticsearch tagline
	}
	out := map[string]string{"service": "elasticsearch"}
	if v := jsonString(resp, "number"); v != "" {
		out["version"] = v
	}
	return out
}

// jsonString returns the string value of the first "key":"value" pair, or "".
// A deliberately tiny scanner so the module stays small (no encoding/json).
func jsonString(s, key string) string {
	i := strings.Index(s, "\""+key+"\"")
	if i < 0 {
		return ""
	}
	i += len(key) + 2
	// skip spaces and the colon
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
