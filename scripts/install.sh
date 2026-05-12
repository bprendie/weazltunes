#!/usr/bin/env bash
set -euo pipefail

root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
prefix="${WEAZLTUNES_PREFIX:-$HOME/.weazltunes}"
bin_dir="$prefix/bin"

mkdir -p "$bin_dir"
go build -o "$bin_dir/weazltunes" "$root/cmd/weazltunes"

case ":$PATH:" in
  *":$bin_dir:"*) ;;
  *)
    shell_rc="$HOME/.bashrc"
    if [ -n "${ZSH_VERSION:-}" ]; then
      shell_rc="$HOME/.zshrc"
    fi
    printf '\nexport PATH="%s:$PATH"\n' "$bin_dir" >> "$shell_rc"
    ;;
esac

printf 'installed %s\n' "$bin_dir/weazltunes"
