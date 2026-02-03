# Worktree Util

A simple TUI (Terminal User Interface) for managing Git worktrees, built with Go and Bubble Tea.

## Features

- 📋 **List** all git worktrees in your repository
- ➕ **Add** new worktrees with custom branches
- 🗑️ **Remove** worktrees safely
- 🎨 Beautiful terminal interface with keyboard navigation
- ⚠️ **Smart error handling** with helpful messages

## Installation

```bash
go build -o worktree-util
```

## Usage

Run the tool from within a git repository:

```bash
./worktree-util
```

### Keyboard Shortcuts

#### List View
- `a` - Add a new worktree
- `d` - Delete selected worktree
- `r` - Refresh the list
- `↑/↓` - Navigate through worktrees
- `q` - Quit

#### Add Worktree View
- `Tab` - Switch between path and branch input fields
- `Enter` - Create the worktree
- `Esc` - Cancel and return to list

#### Delete Confirmation
- `y` - Confirm deletion
- `n` or `Esc` - Cancel deletion

## Requirements

- Go 1.21 or higher
- Git installed and available in PATH
- Must be run from within a git repository

## How It Works

The tool uses `git worktree` commands under the hood:
- `git worktree list --porcelain` - to list worktrees
- `git worktree add` - to create new worktrees
- `git worktree remove` - to delete worktrees

## Example Workflow

1. Run `./worktree-util` in your git repository
2. Press `a` to add a new worktree
3. Enter the path (e.g., `../my-feature`)
4. Enter the branch name (e.g., `feature/new-feature`)
5. Press `Enter` to create
6. The new worktree will appear in the list

## License

MIT

# worktree-util
