#!/bin/bash

extract() {
  set -e

  FILE="$1"
  DEST_DIR="$2"

  if [ -z "$FILE" ]; then
    >&2  echo "Usage: extract <file>"
    return 1
  fi

  if [ ! -f "$FILE" ]; then
    >&2 echo "Error: File not found: $FILE"
    return 1
  fi

  if [ -z "$DEST_DIR" ]; then
    # If no destination directory provided, use the current directory
    DEST_DIR="."
  fi

    # Create the destination directory if it doesn't exist
  if [ ! -d "$DEST_DIR" ]; then
    >&2 echo "Destination directory does not exist. Creating: $DEST_DIR"
    mkdir -p "$DEST_DIR"
  fi

  >&2 echo "Extracting: $FILE"

  case "$FILE" in
    *.tar.gz|*.tgz)
      tar -xzf "$FILE" -C "$DEST_DIR"
      ;;
    *.tar)
      tar -xf "$FILE" -C "$DEST_DIR"
      ;;
    *.gz)
      gunzip -c "$FILE" > "$DEST_DIR/$(basename "$FILE" .gz)"
      ;;
    *)
      >&2 echo "Unsupported file type: $FILE"
      >&2 echo "Only: .tar.gz, .tgz, .tar, .gz"
      return 1
      ;;
  esac

  >&2 echo "Done."
}
