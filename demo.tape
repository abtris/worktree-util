# Worktree Util Demo
# This demo showcases the main features of worktree-util:
# 1. Creating a new worktree
# 2. Creating a worktree from an existing branch
# 3. Removing a worktree

Output assets/demo.gif

Require worktree-util

Set Shell "bash"
Set FontSize 18
Set Width 1400
Set Height 800
Set Padding 20
Set TypingSpeed 100ms

# Start the worktree-util TUI
Type "worktree-util"
Sleep 500ms
Enter
Sleep 2s

# Show the initial list of worktrees
Sleep 1.5s

# Demo 1: Add a new worktree
# Press 'a' to add a new worktree
Type "a"
Sleep 1s

# Type the branch name
Type "feature/demo-feature"
Sleep 1s

# Press Enter to create the worktree
Enter
Sleep 2s

# Show the updated list with the new worktree
Sleep 1.5s

# Demo 2: Create worktree from existing branch
# Press 'c' to create from existing branch
Type "c"
Sleep 1.5s

# Navigate down to select a branch (if needed)
Down 2
Sleep 500ms

# Press Enter to create worktree from selected branch
Enter
Sleep 2s

# Show the updated list
Sleep 1.5s

# Demo 3: Remove a worktree
# Navigate to the worktree we want to remove
Down 1
Sleep 500ms

# Press 'd' to delete
Type "d"
Sleep 1s

# Confirm deletion with 'y'
Type "y"
Sleep 2s

# Show the final list
Sleep 1.5s

# Quit the application
Type "q"
Sleep 1s
