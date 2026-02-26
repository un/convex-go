#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
ALLOWLIST_FILE="${ROOT_DIR}/scripts/no_scaffold_guard.allowlist"

RUNTIME_PATHS=(
  "${ROOT_DIR}/convex"
  "${ROOT_DIR}/internal/protocol"
  "${ROOT_DIR}/internal/sync"
  "${ROOT_DIR}/internal/baseclient"
)

PATTERNS=(
  "TODO\\(runtime-parity\\)"
  "RUNTIME_SCAFFOLD"
  "scaffold-only"
  "placeholder runtime"
  "panic\\(\"not implemented\"\\)"
  "transition chunk unsupported"
)

tmp_matches="$(mktemp)"
trap 'rm -f "$tmp_matches"' EXIT

for pattern in "${PATTERNS[@]}"; do
  rg -n --glob '*.go' --fixed-strings "$pattern" "${RUNTIME_PATHS[@]}" >>"$tmp_matches" || true
done

if [[ -s "$tmp_matches" && -f "$ALLOWLIST_FILE" ]]; then
  filtered="$(mktemp)"
  trap 'rm -f "$tmp_matches" "$filtered"' EXIT
  cp "$tmp_matches" "$filtered"
  while IFS= read -r allowed; do
    if [[ -z "$allowed" || "$allowed" == \#* ]]; then
      continue
    fi
    grep -F -v "$allowed" "$filtered" >"${filtered}.next" || true
    mv "${filtered}.next" "$filtered"
  done <"$ALLOWLIST_FILE"
  mv "$filtered" "$tmp_matches"
fi

if [[ -s "$tmp_matches" ]]; then
  echo "No-scaffold guard failed. Forbidden runtime scaffold markers detected:"
  cat "$tmp_matches"
  exit 1
fi

echo "No-scaffold guard passed."
