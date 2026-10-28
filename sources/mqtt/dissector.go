// A tcpcat dissector: recognises an MQTT broker.
//
// It sends a minimal MQTT CONNECT and looks for a CONNACK (packet type 0x20)
// in reply -- the broker's handshake. The CONNACK's return code tells whether
// the broker accepts anonymous clients, which is the interesting posture
// question for an exposed broker. MQTT carries no version string, only a
// protocol level, so none is reported.
//
// Build:  ./build.sh mqtt     ->  mqtt.wasm
// Run:    tcpcat -sT -p 1883 --scripts . <target>
package main

func Dissect(ip string, port int) map[string]string {
	// CONNECT: protocol "MQTT", level 4 (3.1.1), clean session, empty client id.
	connect := []byte{
		0x10, 0x0c, // CONNECT, remaining length 12
		0x00, 0x04, 'M', 'Q', 'T', 'T', // protocol name
		0x04,       // protocol level 4 (MQTT 3.1.1)
		0x02,       // connect flags: clean session
		0x00, 0x3c, // keep-alive 60s
		0x00, 0x00, // client id length 0
	}
	resp := Probe(connect)
	if len(resp) < 4 || resp[0] != 0x20 { // CONNACK fixed header byte
		return nil
	}
	out := map[string]string{"service": "mqtt"}
	switch resp[3] { // connect return / reason code
	case 0x00:
		out["anonymous"] = "accepted"
	case 0x04, 0x05, 0x87, 0x86:
		out["anonymous"] = "rejected (auth required)"
	default:
		out["anonymous"] = "refused"
	}
	return out
}
