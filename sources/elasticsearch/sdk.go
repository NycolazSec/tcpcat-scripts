// tcpcat dissector SDK for TinyGo -- the glue between Dissect (dissector.go)
// and tcpcat's WebAssembly ABI. You should not need to edit this file; see
// docs/WASM_SDK.md for the ABI if you want to write a dissector in another
// language (Rust, Zig, C...).
package main

import "unsafe"

// tcpcat.probe: send req to the port being scanned, receive up to cap bytes.
// Returns the byte count, or a negative error (-1 connect, -2 too many
// probes in one run, -3 bad buffer).
//
//go:wasmimport tcpcat probe
func hostProbe(reqPtr, reqLen, respPtr, respCap uint32) int32

// Probe sends req to the port being scanned and returns what it answers
// (nil on error). An empty req just reads the banner. At most 4 probes per
// run and 64 KiB per answer.
func Probe(req []byte) []byte {
	resp := make([]byte, 16*1024)
	var reqPtr uint32
	if len(req) > 0 {
		reqPtr = uint32(uintptr(unsafe.Pointer(&req[0])))
	}
	n := hostProbe(reqPtr, uint32(len(req)), uint32(uintptr(unsafe.Pointer(&resp[0]))), uint32(len(resp)))
	if n < 0 {
		return nil
	}
	return resp[:n]
}

// Buffers shared with the host stay referenced here so the garbage
// collector leaves them alone until the host calls deallocate.
var pinned = map[uint32][]byte{}

func pin(buf []byte) uint32 {
	ptr := uint32(uintptr(unsafe.Pointer(&buf[0])))
	pinned[ptr] = buf
	return ptr
}

//go:wasmexport allocate
func allocate(size uint32) uint32 {
	if size == 0 {
		size = 1
	}
	return pin(make([]byte, size))
}

//go:wasmexport deallocate
func deallocate(ptr, size uint32) { delete(pinned, ptr) }

//go:wasmexport run
func run(port, ipPtr, ipLen uint32) uint64 {
	ip := string(pinned[ipPtr][:ipLen])
	out := Dissect(ip, int(port))
	if len(out) == 0 {
		return 0
	}
	js := encode(out)
	return uint64(pin(js))<<32 | uint64(len(js))
}

// encode writes a map as a JSON object (a few lines instead of pulling in
// encoding/json, which would triple the module's size).
func encode(m map[string]string) []byte {
	b := []byte{'{'}
	first := true
	for k, v := range m {
		if !first {
			b = append(b, ',')
		}
		first = false
		b = quote(append(quote(b, k), ':'), v)
	}
	return append(b, '}')
}

func quote(b []byte, s string) []byte {
	const hex = "0123456789abcdef"
	b = append(b, '"')
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case c == '"' || c == '\\':
			b = append(b, '\\', c)
		case c < 0x20:
			b = append(b, '\\', 'u', '0', '0', hex[c>>4], hex[c&0xf])
		default:
			b = append(b, c)
		}
	}
	return append(b, '"')
}

func main() {}
