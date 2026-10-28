// A tcpcat dissector: recognises PostgreSQL and reports its SSL posture.
//
// The PostgreSQL SSLRequest is answered before authentication with a single
// byte: 'S' (TLS available) or 'N' (TLS not available). That one byte both
// confirms the service and tells you whether the wire is encryptable -- a
// plaintext-only Postgres on the network is worth flagging. The server
// version needs an authenticated session, so it is not reported here.
//
// Build:  ./build.sh postgres     ->  postgres.wasm
// Run:    tcpcat -sT -p 5432 --scripts . <target>
package main

func Dissect(ip string, port int) map[string]string {
	// SSLRequest message: int32 length = 8, int32 code = 80877103 (0x04d2162f).
	reply := Probe([]byte{0x00, 0x00, 0x00, 0x08, 0x04, 0xd2, 0x16, 0x2f})
	if len(reply) != 1 {
		return nil // PostgreSQL answers SSLRequest with exactly one byte
	}
	switch reply[0] {
	case 'S':
		return map[string]string{"service": "postgresql", "tls": "supported"}
	case 'N':
		return map[string]string{"service": "postgresql", "tls": "unavailable"}
	}
	return nil
}
