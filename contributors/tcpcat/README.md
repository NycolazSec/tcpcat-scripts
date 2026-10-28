# tcpcat — tcpcat dissector author

[![tcpcat dissector author](https://raw.githubusercontent.com/NycolazSec/tcpcat-scripts/main/contributors/tcpcat/badge.svg)](https://github.com/NycolazSec/tcpcat-scripts/blob/main/contributors/tcpcat/README.md)

4 dissector(s) published in the [tcpcat scripts hub](../../README.md):

| dissector | version | description |
|---|---|---|
| `memcached` | 1.0.0 | Detects memcached and reads its version; flags unauthenticated instances. |
| `mysql` | 1.0.0 | Detects MySQL / MariaDB and reads the server version from its pre-auth handshake. |
| `postgres` | 1.0.0 | Detects PostgreSQL and reports its TLS posture via the pre-auth SSLRequest. |
| `redis` | 1.0.0 | Detects Redis and reads its version; flags unauthenticated instances. |

## Show it off

Add the badge to your own README (it links back here):

```markdown
[![tcpcat dissector author](https://raw.githubusercontent.com/NycolazSec/tcpcat-scripts/main/contributors/tcpcat/badge.svg)](https://github.com/NycolazSec/tcpcat-scripts/blob/main/contributors/tcpcat/README.md)
```

Share text:

```
My dissectors (memcached, mysql, postgres, redis) are published in the tcpcat scripts hub — sandboxed WASM detection for the tcpcat network scanner. https://github.com/NycolazSec/tcpcat
```

Want your own? [Publish a dissector](https://github.com/NycolazSec/tcpcat-scripts/blob/main/README.md#contributing-a-dissector).
