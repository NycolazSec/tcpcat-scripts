#!/usr/bin/env python3
"""Generate the dissector-author badges and pages from index.json.

For each distinct `author` in the index this writes:
  contributors/<slug>/badge.svg   -- a flat "tcpcat | dissector author" badge
  contributors/<slug>/README.md   -- the author's page: their dissectors, a
                                      copy-paste README badge snippet (a backlink
                                      to this hub), and pre-written share text
and a top-level CONTRIBUTORS.md listing everyone.

Run from the repo root:  python3 tools/build_contributors.py
CI runs it and fails if the committed output is stale (see validate.yml)."""
import html
import json
import os
import re
import sys

REPO = "NycolazSec/tcpcat-scripts"
RAW = f"https://raw.githubusercontent.com/{REPO}/main"
BLOB = f"https://github.com/{REPO}/blob/main"
TCPCAT = "https://github.com/NycolazSec/tcpcat"

GREEN = "#3ddc84"
DARK = "#1b2733"


def slug(name: str) -> str:
    s = re.sub(r"[^a-z0-9]+", "-", name.lower()).strip("-")
    return s or "anonymous"


def text_width(s: str) -> int:
    # Rough Verdana-11 advance; good enough for a stable, legible badge.
    return int(len(s) * 7) + 12


def badge_svg(left: str, right: str) -> str:
    lw, rw = text_width(left), text_width(right)
    total = lw + rw
    le, re_ = html.escape(left), html.escape(right)
    return f"""<svg xmlns="http://www.w3.org/2000/svg" width="{total}" height="20" role="img" aria-label="{le}: {re_}">
  <title>{le}: {re_}</title>
  <linearGradient id="s" x2="0" y2="100%">
    <stop offset="0" stop-color="#bbb" stop-opacity=".1"/>
    <stop offset="1" stop-opacity=".1"/>
  </linearGradient>
  <clipPath id="r"><rect width="{total}" height="20" rx="3" fill="#fff"/></clipPath>
  <g clip-path="url(#r)">
    <rect width="{lw}" height="20" fill="{DARK}"/>
    <rect x="{lw}" width="{rw}" height="20" fill="{GREEN}"/>
    <rect width="{total}" height="20" fill="url(#s)"/>
  </g>
  <g text-anchor="middle" font-family="Verdana,Geneva,DejaVu Sans,sans-serif" font-size="11">
    <text x="{lw/2:.0f}" y="14" fill="#fff">{le}</text>
    <text x="{lw + rw/2:.0f}" y="14" fill="#0b0f14">{re_}</text>
  </g>
</svg>
"""


def badge_snippet(sl: str) -> str:
    return (
        f"[![tcpcat dissector author]({RAW}/contributors/{sl}/badge.svg)]"
        f"({BLOB}/contributors/{sl}/README.md)"
    )


def main() -> int:
    root = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
    with open(os.path.join(root, "index.json"), "rb") as f:
        idx = json.load(f)

    authors: dict[str, dict] = {}
    for d in idx.get("dissectors", []):
        name = (d.get("author") or "anonymous").strip()
        a = authors.setdefault(name, {"slug": slug(name), "dissectors": []})
        a["dissectors"].append(d)

    # Clean out stale author dirs so a renamed/removed author doesn't linger.
    contrib_root = os.path.join(root, "contributors")
    keep = {a["slug"] for a in authors.values()}
    if os.path.isdir(contrib_root):
        for entry in os.listdir(contrib_root):
            if entry not in keep and os.path.isdir(os.path.join(contrib_root, entry)):
                for f in os.listdir(os.path.join(contrib_root, entry)):
                    os.remove(os.path.join(contrib_root, entry, f))
                os.rmdir(os.path.join(contrib_root, entry))

    rows = []
    for name in sorted(authors):
        a = authors[name]
        sl = a["slug"]
        ds = sorted(a["dissectors"], key=lambda d: d["name"])
        d_dir = os.path.join(contrib_root, sl)
        os.makedirs(d_dir, exist_ok=True)

        with open(os.path.join(d_dir, "badge.svg"), "w") as f:
            f.write(badge_svg("tcpcat", "dissector author"))

        names = ", ".join(f"`{d['name']}`" for d in ds)
        share = (
            f"My {'dissector' if len(ds) == 1 else 'dissectors'} "
            f"({names.replace('`', '')}) {'is' if len(ds) == 1 else 'are'} published in the "
            f"tcpcat scripts hub — sandboxed WASM detection for the tcpcat network scanner. {TCPCAT}"
        )
        lines = [
            f"# {name} — tcpcat dissector author",
            "",
            badge_snippet(sl),
            "",
            f"{len(ds)} dissector(s) published in the [tcpcat scripts hub](../../README.md):",
            "",
            "| dissector | version | description |",
            "|---|---|---|",
        ]
        for d in ds:
            lines.append(f"| `{d['name']}` | {d.get('version','')} | {d.get('description','')} |")
        lines += [
            "",
            "## Show it off",
            "",
            "Add the badge to your own README (it links back here):",
            "",
            "```markdown",
            badge_snippet(sl),
            "```",
            "",
            "Share text:",
            "",
            "```",
            share,
            "```",
            "",
            f"Want your own? [Publish a dissector]({BLOB}/README.md#contributing-a-dissector).",
            "",
        ]
        with open(os.path.join(d_dir, "README.md"), "w") as f:
            f.write("\n".join(lines))

        rows.append((name, sl, len(ds), names))

    # Top-level CONTRIBUTORS.md
    out = [
        "# Dissector authors",
        "",
        "Everyone with a dissector in this hub. Publish one and you get a badge "
        "to show off — see [Contributing](README.md#contributing-a-dissector).",
        "",
        "| author | dissectors | badge |",
        "|---|---|---|",
    ]
    for name, sl, n, names in rows:
        out.append(f"| [{name}](contributors/{sl}/README.md) | {n} ({names}) | ![badge](contributors/{sl}/badge.svg) |")
    out.append("")
    with open(os.path.join(root, "CONTRIBUTORS.md"), "w") as f:
        f.write("\n".join(out))

    print(f"generated {len(rows)} contributor(s)")
    return 0


if __name__ == "__main__":
    sys.exit(main())
