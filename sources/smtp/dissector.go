// A tcpcat dissector: recognises an SMTP server and reads its greeting.
//
// An SMTP server announces itself the moment you connect with a "220 ..."
// greeting that usually names the mail software -- all read before any
// command, no authentication. Requiring "smtp" in the greeting keeps this
// from matching other "220 ..." services (e.g. FTP).
//
// Build:  ./build.sh smtp     ->  smtp.wasm
// Run:    tcpcat -sT -p 25,587 --scripts . <target>
package main

import "strings"

func Dissect(ip string, port int) map[string]string {
	banner := strings.TrimSpace(string(Probe(nil)))
	low := strings.ToLower(banner)
	if !strings.HasPrefix(banner, "220") || !strings.Contains(low, "smtp") {
		return nil
	}
	out := map[string]string{"service": "smtp", "banner": banner}
	for _, p := range []string{"postfix", "exim", "sendmail", "exchange", "opensmtpd", "qmail", "haraka", "microsoft"} {
		if strings.Contains(low, p) {
			out["product"] = p
			break
		}
	}
	return out
}
