// A tcpcat dissector: recognises an FTP server and reads its greeting.
//
// An FTP server sends a "220 ..." greeting on connect, before login, that
// usually names the daemon (vsFTPd, ProFTPD, Pure-FTPd...). Requiring "ftp"
// in the greeting avoids matching other "220 ..." services such as SMTP.
//
// Build:  ./build.sh ftp     ->  ftp.wasm
// Run:    tcpcat -sT -p 21 --scripts . <target>
package main

import "strings"

func Dissect(ip string, port int) map[string]string {
	banner := strings.TrimSpace(string(Probe(nil)))
	low := strings.ToLower(banner)
	if !strings.HasPrefix(banner, "220") || !strings.Contains(low, "ftp") {
		return nil
	}
	out := map[string]string{"service": "ftp", "banner": banner}
	for _, p := range []string{"vsftpd", "proftpd", "pure-ftpd", "filezilla", "wu-ftpd", "pyftpd", "microsoft ftp"} {
		if strings.Contains(low, p) {
			out["product"] = p
			break
		}
	}
	return out
}
