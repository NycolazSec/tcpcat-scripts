// A tcpcat dissector: recognises MongoDB and estimates its version.
//
// It sends the legacy, unauthenticated `isMaster` command as an OP_QUERY --
// the handshake every MongoDB answers before auth. The reply is an OP_REPLY
// (opcode 1) whose BSON carries `maxWireVersion`, which maps to a release.
// An exposed MongoDB is a well-known data-leak class, so confirming it is
// worthwhile; the exact version string needs an authenticated buildInfo.
//
// Build:  ./build.sh mongodb     ->  mongodb.wasm
// Run:    tcpcat -sT -p 27017 --scripts . <target>
package main

import (
	"strconv"
	"strings"
)

func Dissect(ip string, port int) map[string]string {
	// A fixed 58-byte OP_QUERY for {isMaster: 1} on admin.$cmd.
	req := []byte{
		0x3a, 0x00, 0x00, 0x00, // messageLength = 58
		0x01, 0x00, 0x00, 0x00, // requestID
		0x00, 0x00, 0x00, 0x00, // responseTo
		0xd4, 0x07, 0x00, 0x00, // opCode = 2004 (OP_QUERY)
		0x00, 0x00, 0x00, 0x00, // flags
		'a', 'd', 'm', 'i', 'n', '.', '$', 'c', 'm', 'd', 0x00, // "admin.$cmd"
		0x00, 0x00, 0x00, 0x00, // numberToSkip
		0x01, 0x00, 0x00, 0x00, // numberToReturn
		// BSON {isMaster: 1}
		0x13, 0x00, 0x00, 0x00,
		0x10, 'i', 's', 'M', 'a', 's', 't', 'e', 'r', 0x00,
		0x01, 0x00, 0x00, 0x00,
		0x00,
	}
	resp := Probe(req)
	if len(resp) < 16 {
		return nil
	}
	// Header opCode (little-endian) must be 1 (OP_REPLY).
	opcode := uint32(resp[12]) | uint32(resp[13])<<8 | uint32(resp[14])<<16 | uint32(resp[15])<<24
	if opcode != 1 {
		return nil
	}

	out := map[string]string{"service": "mongodb"}
	// Find the maxWireVersion int32 in the BSON reply (type 0x10, name, int32).
	const key = "maxWireVersion"
	if i := strings.Index(string(resp), key); i >= 0 {
		p := i + len(key)
		if p < len(resp) && resp[p] == 0x00 && p+5 <= len(resp) {
			mwv := int(resp[p+1]) | int(resp[p+2])<<8 | int(resp[p+3])<<16 | int(resp[p+4])<<24
			out["maxwireversion"] = strconv.Itoa(mwv)
			if v := versionFromWire(mwv); v != "" {
				out["version"] = v
			}
		}
	}
	return out
}

// versionFromWire maps maxWireVersion to the MongoDB release that introduced it.
func versionFromWire(w int) string {
	switch {
	case w >= 25:
		return "8.0+"
	case w == 21:
		return "7.0"
	case w == 17:
		return "6.0"
	case w == 13:
		return "5.0"
	case w == 9:
		return "4.4"
	case w == 8:
		return "4.2"
	case w == 7:
		return "4.0"
	}
	return ""
}
