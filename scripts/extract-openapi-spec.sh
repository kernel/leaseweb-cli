#!/usr/bin/env bash
#
# Extract the Leaseweb OpenAPI spec from the developer docs.
#
# The spec is embedded in the HTML as a JavaScript object literal inside the
# __redoc_state variable (Redocly). This script:
#   1. Fetches the HTML from https://developer.leaseweb.com/docs/
#   2. Uses Python to locate and extract the JS object via brace-matching
#   3. Preprocesses JS shorthand (!0 -> true, !1 -> false)
#   4. Uses Node.js to eval the JS object literal into valid JSON
#
# Requirements: curl, python3, node
#
set -euo pipefail

OUTDIR="${1:-docs}"
mkdir -p "$OUTDIR"

TMPDIR=$(mktemp -d)
trap 'rm -rf "$TMPDIR"' EXIT

echo "Fetching developer docs..."
curl -sL "https://developer.leaseweb.com/docs/" -o "$TMPDIR/page.html"

echo "Extracting OpenAPI spec from embedded JS..."
python3 - "$TMPDIR/page.html" "$TMPDIR/spec.js" <<'PYEOF'
import sys, re

html = open(sys.argv[1]).read()

# The HTML contains: __redoc_state={menu:{...},spec:{data:{...THE SPEC...},...},searchIndex:{...}}
# We want just {THE SPEC} -- the value of the `data:` key inside `spec:`.
marker = 'spec:{data:'
idx = html.find('__redoc_state=')
if idx == -1:
    print("ERROR: Could not find __redoc_state in HTML", file=sys.stderr)
    sys.exit(1)

spec_idx = html.find(marker, idx)
if spec_idx == -1:
    print("ERROR: Could not find spec:{data: in __redoc_state", file=sys.stderr)
    sys.exit(1)

start = spec_idx + len(marker)

# Brace-match to find the end of the OpenAPI spec object.
# Must handle JS strings (double-quoted, single-quoted, and backtick template literals).
depth = 0
i = start
n = len(html)
while i < n:
    ch = html[i]
    if ch == '{':
        depth += 1
    elif ch == '}':
        if depth == 0:
            break
        depth -= 1
    elif ch == '"':
        i += 1
        while i < n and html[i] != '"':
            if html[i] == '\\':
                i += 1
            i += 1
    elif ch == "'":
        i += 1
        while i < n and html[i] != "'":
            if html[i] == '\\':
                i += 1
            i += 1
    elif ch == '`':
        i += 1
        while i < n and html[i] != '`':
            if html[i] == '\\':
                i += 1
            i += 1
    i += 1

js = html[start:i]

# Replace JS boolean shorthand: !0 -> true, !1 -> false
js = re.sub(r'(?<![a-zA-Z0-9_])!0(?![a-zA-Z0-9_])', 'true', js)
js = re.sub(r'(?<![a-zA-Z0-9_])!1(?![a-zA-Z0-9_])', 'false', js)
js = js.replace('void 0', 'null')

open(sys.argv[2], 'w').write(js)
print(f"Extracted {len(js)} bytes of JS object literal")
PYEOF

echo "Converting JS object literal to JSON via Node.js..."
node -e "
const fs = require('fs');
const js = fs.readFileSync('$TMPDIR/spec.js', 'utf8');
const obj = new Function('return (' + js + ')')();
const json = JSON.stringify(obj, null, 2);
fs.writeFileSync('$OUTDIR/leaseweb-openapi.json', json);
const paths = Object.keys(obj.paths || {}).length;
const ops = Object.values(obj.paths || {}).reduce((n, p) => n + Object.keys(p).filter(m => ['get','post','put','patch','delete'].includes(m)).length, 0);
console.log('Wrote ' + paths + ' paths, ' + ops + ' operations to $OUTDIR/leaseweb-openapi.json');
"

echo "Generating API reference markdown..."
python3 - "$OUTDIR/leaseweb-openapi.json" "$OUTDIR/leaseweb-api-reference.md" <<'PYEOF'
import json, sys

spec = json.load(open(sys.argv[1]))
out = open(sys.argv[2], 'w')

out.write("# Leaseweb API Reference\n\n")
out.write(f"OpenAPI {spec.get('openapi', '?')} — {spec.get('info', {}).get('title', '?')}\n\n")

tag_groups = spec.get('x-tagGroups', [])
if tag_groups:
    for group in tag_groups:
        out.write(f"## {group['name']}\n\n")
        for tag_name in group.get('tags', []):
            out.write(f"### {tag_name}\n\n")
            for path, methods in sorted(spec.get('paths', {}).items()):
                for method, op in methods.items():
                    if method not in ('get', 'post', 'put', 'patch', 'delete'):
                        continue
                    tags = op.get('tags', [])
                    if tag_name in tags:
                        summary = op.get('summary', '')
                        out.write(f"- `{method.upper()} {path}` — {summary}\n")
            out.write("\n")
else:
    for path, methods in sorted(spec.get('paths', {}).items()):
        for method, op in methods.items():
            if method not in ('get', 'post', 'put', 'patch', 'delete'):
                continue
            summary = op.get('summary', '')
            tags = ', '.join(op.get('tags', []))
            out.write(f"- `{method.upper()} {path}` [{tags}] — {summary}\n")

print(f"Wrote API reference to {sys.argv[2]}")
PYEOF

echo "Done."
