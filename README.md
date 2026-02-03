# Worktree Util

<p align="center">
  <img src="assets/logo.png" alt="Worktree Util Logo" width="200"/>
</p>

[![Release](https://img.shields.io/github/v/release/abtris/worktree-util)](https://github.com/abtris/worktree-util/releases)
[![Test](https://github.com/abtris/worktree-util/workflows/Test/badge.svg)](https://github.com/abtris/worktree-util/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/abtris/worktree-util)](https://goreportcard.com/report/github.com/abtris/worktree-util)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A simple TUI (Terminal User Interface) for managing Git worktrees, built with Go and Bubble Tea.

## Features

- 📋 **List** all git worktrees in your repository
- 🚀 **Quick navigation** - press Enter to instantly change to a worktree directory
- ➕ **Add** new worktrees with custom branches (auto-organized in `.worktrees/` folder)
- 🗑️ **Remove** worktrees safely
- 🎨 Beautiful terminal interface with keyboard navigation
- ⚠️ **Smart error handling** with helpful messages
- 🔧 **Simple workflow** - just enter a branch name, path is auto-generated

## Installation

### Homebrew (macOS/Linux)

```bash
brew install abtris/tap/worktree-util
```

### Download Binary

Download the latest release for your platform from the [releases page](https://github.com/abtris/worktree-util/releases).

### Build from Source

```bash
git clone https://github.com/abtris/worktree-util.git
cd worktree-util
go build -o worktree-util
```

## Usage

### TUI Mode

Run the tool from within a git repository to start the interactive TUI:

```bash
./worktree-util
```

### Quick Directory Change (Recommended Setup)

To enable changing to a worktree directory by pressing Enter, add this function to your shell configuration:

**For Bash (~/.bashrc):**
```bash
wt() {
    local output
    output=$(worktree-util "$@")
    local exit_code=$?

    if [ $exit_code -eq 0 ] && [ -n "$output" ] && [ -d "$output" ]; then
        cd "$output" || return 1
    elif [ -n "$output" ]; then
        echo "$output"
    fi

    return $exit_code
}
```

**For Zsh (~/.zshrc):**
```zsh
wt() {
    local output
    output=$(worktree-util "$@")
    local exit_code=$?

    if [ $exit_code -eq 0 ] && [ -n "$output" ] && [ -d "$output" ]; then
        cd "$output" || return 1
    elif [ -n "$output" ]; then
        echo "$output"
    fi

    return $exit_code
}
```

After adding the function, reload your shell config:
```bash
source ~/.bashrc  # or ~/.zshrc
```

Now you can use `wt` to launch the tool and press Enter on any worktree to instantly change to that directory!

### CLI Commands

Manage configuration from the command line:

```bash
# Show current configuration
worktree-util config

# Initialize config file with defaults
worktree-util config init

# Set configuration values
worktree-util config set worktree_dir my-worktrees

# Get configuration values
worktree-util config get worktree_dir

# Add files to copy to new worktrees
worktree-util config add-copy-file .env
worktree-util config add-copy-file .env.local

# Remove files from copy list
worktree-util config remove-copy-file .env

# Show help
worktree-util --help

# Show version
worktree-util --version
```

### Keyboard Shortcuts

#### List View
- `Enter` - Change to selected worktree directory (requires shell wrapper - see above)
- `a` - Add a new worktree
- `c` - Create worktree from existing branch (shows searchable list of local and remote branches)
- `d` - Delete selected worktree
- `r` - Refresh the list
- `↑/↓` - Navigate through worktrees
- `q` - Quit

#### Add Worktree View
- `Enter` - Create the worktree
- `Esc` - Cancel and return to list

#### Branch Selection View
- `↑/↓` or `j/k` - Navigate through branches
- `/` - Filter/search branches
- `Enter` - Create worktree from selected branch
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

### Auto-Generated Paths

When you create a new worktree, you only need to provide the branch name. The tool automatically:
1. Creates a `.worktrees/` directory in your repository root (if it doesn't exist)
2. Sanitizes the branch name (e.g., `feature/new-feature` → `feature-new-feature`)
3. Creates the worktree at `.worktrees/<sanitized-branch-name>`

This keeps all your worktrees organized in one place!

## Configuration

Worktree Util supports optional configuration via `~/.config/worktree-util/config.yml`.

### Creating a Configuration File

1. Create the config directory:
   ```bash
   mkdir -p ~/.config/worktree-util
   ```

2. Create `~/.config/worktree-util/config.yml`:
   ```yaml
   # Directory where worktrees will be created (relative to repo root)
   worktree_dir: .worktrees

   # Files to copy from repo root to new worktrees
   copy_files:
     - .env
     - .env.local
   ```

### Configuration Options

- **`worktree_dir`**: Directory where worktrees will be created (relative to repository root)
  - Default: `.worktrees`
  - Examples:
    - `.worktrees` - Creates worktrees in `.worktrees/` folder
    - `worktrees` - Creates worktrees in `worktrees/` folder
    - `../my-worktrees` - Creates worktrees outside the repository

- **`copy_files`**: List of files to copy from repository root to new worktrees
  - Default: `[]` (no files copied)
  - Useful for files in `.gitignore` that are needed for development (e.g., `.env` files)
  - Files that don't exist will be silently skipped
  - Supports nested paths (e.g., `config/local.yml`)
  - Examples:
    ```yaml
    copy_files:
      - .env
      - .env.local
      - config/local.yml
    ```

See [`config.example.yml`](config.example.yml) for a complete example.

**Note:** Configuration is completely optional. If no config file exists, the tool uses sensible defaults.

## Example Workflow

1. Run `./worktree-util` in your git repository
2. Press `a` to add a new worktree
3. Enter the branch name (e.g., `feature/new-feature`)
4. Watch the path auto-generate (e.g., `.worktrees/feature-new-feature`)
5. Press `Enter` to create
6. The new worktree will appear in the list

Your worktrees will be organized like this:
```
my-repo/
├── .git/
├── .worktrees/
│   ├── feature-new-feature/
│   ├── bugfix-123/
│   └── experiment-api/
├── src/
└── README.md
```

## License

MIT

