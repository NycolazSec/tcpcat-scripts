# tcpcat-scripts

The community hub of **WebAssembly detection modules (dissectors)** for
[tcpcat](https://github.com/NycolazSec/tcpcat). It's tcpcat's equivalent of
Nmap's NSE — but every module runs in a sandbox and talks only to the port
being scanned.

`tcpcat scripts` installs modules from this repo and verifies each one's
**SHA-256** before it ever reaches the scanner.

```bash
tcpcat scripts search                 # browse this hub
tcpcat scripts install redis          # download + verify into the scripts dir
tcpcat -sV -p 6379 --scripts "$(tcpcat scripts dir)" <target>
```

(This repo is the default hub; point `tcpcat` elsewhere with `--index` or
`$TCPCAT_SCRIPTS_INDEX`.)

## Layout

```
index.json                       # the catalog
dissectors/
  redis.wasm                     # a module, referenced from index.json by its sha256
validate.py                      # the PR gate (also runnable locally)
.github/workflows/validate.yml   # runs validate.py on every PR/push
```

## The index format

A hub is a static `index.json` served over HTTPS. Each entry pins the module's
**SHA-256**, so `tcpcat scripts install` rejects a tampered or corrupted
`.wasm`. `url` may be absolute or **relative to `index.json`** (as below), which
keeps the hub self-contained.

```jsonc
{
  "schema_version": 1,
  "dissectors": [
    {
      "name": "redis",                 // also the installed filename: redis.wasm
      "description": "…",
      "version": "1.0.0",
      "author": "…",
      "license": "AGPL-3.0-only",
      "homepage": "https://…",
      "ports": [6379],                 // informational
      "url": "dissectors/redis.wasm",  // absolute, or relative to index.json
      "sha256": "fb1d58a1…41c7"         // required; install is pinned to it
    }
  ]
}
```

## Contributing a dissector

1. Write and build it from the template in
   [tcpcat/examples/wasm-dissector](https://github.com/NycolazSec/tcpcat/tree/main/examples/wasm-dissector)
   (`./build.sh myname` → `myname.wasm`).
2. Add `dissectors/myname.wasm` and an entry in `index.json` with its digest:
   ```bash
   sha256sum dissectors/myname.wasm
   ```
3. Regenerate the author badges and pages, then commit them:
   ```bash
   python3 tools/build_contributors.py
   ```
4. Open a PR. The [validation workflow](.github/workflows/validate.yml) checks
   that every entry's `sha256` matches its file, names are unique and safe, the
   JSON is well-formed, and the badges are up to date — a mismatched or
   malformed module can't be merged.

Run the same checks locally:

```bash
python3 validate.py
python3 tools/build_contributors.py && git diff --exit-code
```

## Dissector authors

Publish a dissector and you get a **badge** to show off — see everyone on the
[dissector authors page](CONTRIBUTORS.md). Each author has a page with a
copy-paste snippet:

[![tcpcat dissector author](https://raw.githubusercontent.com/NycolazSec/tcpcat-scripts/main/contributors/tcpcat/badge.svg)](https://github.com/NycolazSec/tcpcat-scripts/blob/main/contributors/tcpcat/README.md)

Drop it in your own README and it links back here — a small nod to the people
who grow the hub.

## Sandbox

Modules run in tcpcat's sandbox (wazero, no WASI): a dissector has no files, no
clock, no environment, and can reach only the port being scanned, through one
host function, within a strict probe budget. Keep probes read-only. See the
[WASM SDK](https://github.com/NycolazSec/tcpcat/blob/main/docs/WASM_SDK.md) for
the ABI and the sandbox rules.

## License

Each module declares its own `license` in `index.json`. The reference `redis`
module is AGPL-3.0-only, like tcpcat itself.
