#!/usr/bin/env bash
# build-all.sh — export every manual under docs/manuals/src/ in all formats.
# Needs: pandoc, weasyprint (for PDF), zip. Calibre `ebook-convert` is
# optional (adds .mobi).
set -euo pipefail

REPO_ROOT=$(git rev-parse --show-toplevel)
SRC="$REPO_ROOT/docs/manuals/src"
DIST="$REPO_ROOT/docs/manuals/dist"
mkdir -p "$DIST"

have() { command -v "$1" >/dev/null 2>&1; }

if ! have pandoc; then
    echo "pandoc not installed — aborting." >&2
    exit 2
fi

manuals=()
for d in "$SRC"/*/; do
    [ -d "$d" ] || continue
    manuals+=("$(basename "$d")")
done

if [ ${#manuals[@]} -eq 0 ]; then
    echo "No manuals under $SRC yet. Add content and re-run." >&2
    exit 0
fi

for m in "${manuals[@]}"; do
    src_dir="$SRC/$m"
    combined="$DIST/$m.md"
    echo "→ $m"

    # Concatenate chapters in numeric order.
    : >"$combined"
    for f in $(ls "$src_dir"/*.md 2>/dev/null | sort); do
        cat "$f" >>"$combined"
        printf '\n\n' >>"$combined"
    done

    pandoc "$combined" -o "$DIST/$m.pdf" --pdf-engine=weasyprint 2>/dev/null \
        || pandoc "$combined" -o "$DIST/$m.pdf" 2>/dev/null \
        || echo "  warn: PDF skipped"

    pandoc "$combined" -o "$DIST/$m.epub" 2>/dev/null || echo "  warn: epub skipped"
    pandoc "$combined" -o "$DIST/$m.docx" 2>/dev/null || echo "  warn: docx skipped"
    pandoc "$combined" --to plain -o "$DIST/$m.txt" 2>/dev/null || echo "  warn: txt skipped"
    if have ebook-convert && [ -f "$DIST/$m.epub" ]; then
        ebook-convert "$DIST/$m.epub" "$DIST/$m.mobi" >/dev/null 2>&1 || true
    fi
    (cd "$DIST" && zip -q "$m.zip" "$m."*) || true
    rm -f "$combined"
done

echo ""
echo "Exports in $DIST:"
ls -1 "$DIST" | sort
