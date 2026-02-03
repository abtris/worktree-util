#!/bin/bash
# Shell wrapper for worktree-util to enable directory changing
# Add this function to your ~/.bashrc or ~/.zshrc

wt() {
    local output
    output=$(worktree-util "$@")
    local exit_code=$?
    
    # If output is a directory path, cd to it
    if [ $exit_code -eq 0 ] && [ -n "$output" ] && [ -d "$output" ]; then
        cd "$output" || return 1
    elif [ -n "$output" ]; then
        # If there's output but it's not a directory, just print it
        echo "$output"
    fi
    
    return $exit_code
}

