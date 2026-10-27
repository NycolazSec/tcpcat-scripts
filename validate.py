#!/usr/bin/env python3
"""Validate a tcpcat scripts hub: every index entry must have a safe, unique
name and, for a relative module, a file whose SHA-256 matches the index. Run
from the hub root (where index.json lives). Exits non-zero on any problem, so
it works as a PR gate."""
import hashlib
import json
import os
import re
import sys

NAME_RE = re.compile(r"^[A-Za-z0-9][A-Za-z0-9._-]*$")


def main() -> int:
    root = os.path.dirname(os.path.abspath(__file__))
    with open(os.path.join(root, "index.json"), "rb") as f:
        idx = json.load(f)

    errors = []
    if idx.get("schema_version") != 1:
        errors.append(f"unsupported schema_version: {idx.get('schema_version')!r}")

    seen = set()
    dissectors = idx.get("dissectors", [])
    for d in dissectors:
        name = d.get("name", "")
        if not NAME_RE.match(name):
            errors.append(f"unsafe or empty name: {name!r}")
        if name in seen:
            errors.append(f"duplicate name: {name}")
        seen.add(name)

        for field in ("description", "version", "url", "sha256"):
            if not d.get(field):
                errors.append(f"{name or '?'}: missing required field {field!r}")

        want = (d.get("sha256") or "").lower()
        if len(want) != 64 or any(c not in "0123456789abcdef" for c in want):
            errors.append(f"{name}: sha256 is not a 64-char hex digest")

        url = d.get("url", "")
        if "://" in url or url.startswith("/"):
            continue  # absolute URL: integrity is still enforced at install time
        path = os.path.join(root, url)
        if not os.path.isfile(path):
            errors.append(f"{name}: module file not found: {url}")
            continue
        with open(path, "rb") as f:
            got = hashlib.sha256(f.read()).hexdigest()
        if want and got != want:
            errors.append(f"{name}: sha256 mismatch (index {want}, file {got})")

    if errors:
        print("INVALID hub:")
        for e in errors:
            print("  -", e)
        return 1
    print(f"OK: {len(dissectors)} dissector(s) validated")
    return 0


if __name__ == "__main__":
    sys.exit(main())
