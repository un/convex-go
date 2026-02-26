#!/usr/bin/env python3

import argparse
import base64
import json
import re
import struct
import subprocess
from pathlib import Path


def rust_commit_hash(repo: Path) -> str:
    result = subprocess.run(
        ["git", "rev-parse", "HEAD"],
        cwd=repo,
        check=True,
        capture_output=True,
        text=True,
    )
    return result.stdout.strip()


def extract_auth_compat_vectors(json_rs: Path):
    text = json_rs.read_text(encoding="utf-8")
    pattern = re.compile(
        r"let\s+(old_[a-z_]+_auth_message)\s*=\s*json!\((\{.*?\})\);"
    )
    vectors = []
    for match in pattern.finditer(text):
        name = match.group(1)
        payload = json.loads(match.group(2))
        vectors.append(
            {
                "name": name,
                "kind": "client_decode",
                "payload": payload,
                "rust_source": "sync_types/src/types/json.rs:1088",
            }
        )
    return vectors


def timestamp_vectors():
    values = [0, 1, 42, (1 << 53) + 7, (1 << 64) - 1]
    vectors = []
    for value in values:
        encoded = base64.b64encode(struct.pack("<Q", value)).decode("ascii")
        vectors.append(
            {
                "name": f"timestamp_{value}",
                "kind": "timestamp",
                "value": value,
                "encoded": encoded,
                "rust_source": "sync_types/src/types/json.rs:u64_to_string",
            }
        )
    return vectors


def main():
    parser = argparse.ArgumentParser(description="Import protocol fixtures from convex-rs")
    parser.add_argument(
        "--rust-repo",
        default="/Users/omar/code/convex-rs",
        help="Path to convex-rs repository",
    )
    parser.add_argument(
        "--output",
        default="internal/protocol/testdata/rust_fixture_vectors.json",
        help="Output fixture file path",
    )
    args = parser.parse_args()

    rust_repo = Path(args.rust_repo).resolve()
    json_rs = rust_repo / "sync_types/src/types/json.rs"
    if not json_rs.exists():
        raise SystemExit(f"missing source file: {json_rs}")

    fixtures = {
        "source_repo": str(rust_repo),
        "source_commit": rust_commit_hash(rust_repo),
        "source_file": "sync_types/src/types/json.rs",
        "vectors": [],
    }
    fixtures["vectors"].extend(extract_auth_compat_vectors(json_rs))
    fixtures["vectors"].extend(timestamp_vectors())

    output_path = Path(args.output).resolve()
    output_path.parent.mkdir(parents=True, exist_ok=True)
    output_path.write_text(json.dumps(fixtures, indent=2, sort_keys=True) + "\n", encoding="utf-8")


if __name__ == "__main__":
    main()
