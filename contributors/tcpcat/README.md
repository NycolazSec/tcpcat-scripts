# tcpcat — tcpcat dissector author

[![tcpcat dissector author](https://raw.githubusercontent.com/NycolazSec/tcpcat-scripts/main/contributors/tcpcat/badge.svg)](https://github.com/NycolazSec/tcpcat-scripts/blob/main/contributors/tcpcat/README.md)

10 dissector(s) published in the [tcpcat scripts hub](../../README.md):

| dissector | version | description |
|---|---|---|
| `couchdb` | 1.0.0 | Detects Apache CouchDB and reads its version from the unauthenticated GET /. |
| `elasticsearch` | 1.0.0 | Detects Elasticsearch and reads its version from the unauthenticated GET /. |
| `ftp` | 1.0.0 | Recognises an FTP server and reads its greeting (names the daemon). |
| `memcached` | 1.0.0 | Detects memcached and reads its version; flags unauthenticated instances. |
| `mongodb` | 1.0.0 | Detects MongoDB via the unauthenticated isMaster handshake and estimates its version. |
| `mqtt` | 1.0.0 | Detects an MQTT broker via CONNECT/CONNACK and whether it accepts anonymous clients. |
| `mysql` | 1.0.0 | Detects MySQL / MariaDB and reads the server version from its pre-auth handshake. |
| `postgres` | 1.0.0 | Detects PostgreSQL and reports its TLS posture via the pre-auth SSLRequest. |
| `redis` | 1.0.0 | Detects Redis and reads its version; flags unauthenticated instances. |
| `smtp` | 1.0.0 | Recognises an SMTP server and reads its greeting (names the mail software). |

## Show it off

Add the badge to your own README (it links back here):

```markdown
[![tcpcat dissector author](https://raw.githubusercontent.com/NycolazSec/tcpcat-scripts/main/contributors/tcpcat/badge.svg)](https://github.com/NycolazSec/tcpcat-scripts/blob/main/contributors/tcpcat/README.md)
```

Share text:

```
My dissectors (couchdb, elasticsearch, ftp, memcached, mongodb, mqtt, mysql, postgres, redis, smtp) are published in the tcpcat scripts hub — sandboxed WASM detection for the tcpcat network scanner. https://github.com/NycolazSec/tcpcat
```

Want your own? [Publish a dissector](https://github.com/NycolazSec/tcpcat-scripts/blob/main/README.md#contributing-a-dissector).
