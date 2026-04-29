#!/bin/sh
DIR="$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)"
chmod +x "$DIR/scripts/setup-mac.sh"
"$DIR/scripts/setup-mac.sh"
