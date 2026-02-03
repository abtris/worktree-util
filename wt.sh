#!/bin/bash
# Shell wrapper for worktree-util to enable directory changing
# Add this function to your ~/.bashrc or ~/.zshrc
#
# All error handling and logic is in the binary.
# This wrapper just reads the temp file and cd's if it exists.

wt() {
    worktree-util "$@"
    local tmpfile="/tmp/worktree-util-cd"
    if [ -f "$tmpfile" ]; then
        local target=$(cat "$tmpfile")
        rm -f "$tmpfile"
        if [ -d "$target" ]; then
            cd "$target"
        fi
    fi
}

