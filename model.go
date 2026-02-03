package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type mode int

const (
	modeList mode = iota
	modeAdd
	modeCheckout
	modeConfirmDelete
)

type model struct {
	list         list.Model
	branchList   list.Model
	mode         mode
	pathInput    textinput.Model
	branchInput  textinput.Model
	inputFocus   int
	err          error
	message      string
	selectedItem Worktree
	width        int
	height       int
	cdPath       string // Path to cd to when exiting
}

type worktreesLoadedMsg []Worktree
type branchesLoadedMsg []Branch
type errMsg error

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("111")).
			MarginLeft(2)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("246")).
			MarginLeft(2).
			MarginTop(1)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("203")).
			Bold(true).
			MarginLeft(2)

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("108")).
			Bold(true).
			MarginLeft(2)
)

func initialModel() model {
	// Create text input for branch name
	branchInput := textinput.New()
	branchInput.Placeholder = "feature/my-feature"
	branchInput.Focus()
	branchInput.CharLimit = 256
	branchInput.Width = 50

	// Create read-only path display (will be auto-generated)
	pathInput := textinput.New()
	pathInput.Placeholder = "(auto-generated)"
	pathInput.CharLimit = 256
	pathInput.Width = 50
	pathInput.Blur() // Always blurred since it's read-only

	// Create list
	delegate := list.NewDefaultDelegate()
	l := list.New([]list.Item{}, delegate, 0, 0)
	l.Title = "Git Worktrees"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle

	// Create branch list
	branchDelegate := list.NewDefaultDelegate()
	bl := list.New([]list.Item{}, branchDelegate, 0, 0)
	bl.Title = "Select Branch"
	bl.SetShowStatusBar(false)
	bl.SetFilteringEnabled(true)
	bl.Styles.Title = titleStyle

	return model{
		list:        l,
		branchList:  bl,
		mode:        modeList,
		pathInput:   pathInput,
		branchInput: branchInput,
		inputFocus:  0,
	}
}

func (m model) Init() tea.Cmd {
	return loadWorktrees
}

func loadWorktrees() tea.Msg {
	worktrees, err := ListWorktrees()
	if err != nil {
		return errMsg(err)
	}
	return worktreesLoadedMsg(worktrees)
}

func loadBranches() tea.Msg {
	branches, err := GetAllBranches()
	if err != nil {
		return errMsg(err)
	}
	return branchesLoadedMsg(branches)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.list.SetSize(msg.Width, msg.Height-6)
		m.branchList.SetSize(msg.Width, msg.Height-6)
		return m, nil

	case worktreesLoadedMsg:
		items := make([]list.Item, len(msg))
		for i, wt := range msg {
			items[i] = wt
		}
		m.list.SetItems(items)
		m.err = nil // Clear any previous errors on successful load
		return m, nil

	case branchesLoadedMsg:
		items := make([]list.Item, len(msg))
		for i, br := range msg {
			items[i] = br
		}
		m.branchList.SetItems(items)
		m.err = nil
		return m, nil

	case errMsg:
		m.err = msg
		return m, nil

	case tea.KeyMsg:
		// Global keys
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

		switch m.mode {
		case modeList:
			return m.updateList(msg)
		case modeAdd:
			return m.updateAdd(msg)
		case modeCheckout:
			return m.updateCheckout(msg)
		case modeConfirmDelete:
			return m.updateConfirmDelete(msg)
		}
	}

	return m, nil
}

func (m model) View() string {
	var b strings.Builder

	switch m.mode {
	case modeList:
		// Show error prominently if there's one
		if m.err != nil {
			b.WriteString(titleStyle.Render("Git Worktrees"))
			b.WriteString("\n\n")
			b.WriteString(errorStyle.Render(fmt.Sprintf("  ⚠ %v", m.err)))
			b.WriteString("\n\n")

			// Show helpful hints based on error type
			errStr := m.err.Error()
			if strings.Contains(errStr, "not a git repository") {
				b.WriteString(helpStyle.Render("  Please run this tool from within a git repository."))
			}
			b.WriteString("\n\n")
			b.WriteString(helpStyle.Render("r: retry • q: quit"))
		} else if len(m.list.Items()) == 0 {
			// No error but no items - show empty state
			b.WriteString(titleStyle.Render("Git Worktrees"))
			b.WriteString("\n\n")
			b.WriteString(helpStyle.Render("  No worktrees found."))
			b.WriteString("\n")
			b.WriteString(helpStyle.Render("  Press 'a' to create a new branch or 'c' to checkout existing!"))
			b.WriteString("\n\n")
			b.WriteString(helpStyle.Render("a: add new • c: checkout existing • r: refresh • q: quit"))
		} else {
			b.WriteString(m.list.View())
			b.WriteString("\n")
			b.WriteString(helpStyle.Render("enter: cd to worktree • a: add new • c: checkout existing • d: delete • r: refresh • q: quit"))
		}
	case modeAdd:
		b.WriteString(titleStyle.Render("Add New Worktree"))
		b.WriteString("\n\n")
		b.WriteString("  Branch: " + m.branchInput.View() + "\n")

		// Show auto-generated path preview
		pathPreview := m.pathInput.Value()
		if pathPreview == "" {
			pathPreview = "(will be auto-generated in .worktrees/)"
		}
		b.WriteString(fmt.Sprintf("  Path:   %s\n\n", pathPreview))
		b.WriteString(helpStyle.Render("enter: create • esc: cancel"))
	case modeCheckout:
		if m.err != nil {
			b.WriteString(titleStyle.Render("Select Branch"))
			b.WriteString("\n\n")
			b.WriteString(errorStyle.Render(fmt.Sprintf("  ⚠ %v", m.err)))
			b.WriteString("\n\n")
			b.WriteString(helpStyle.Render("esc: back"))
		} else if len(m.branchList.Items()) == 0 {
			b.WriteString(titleStyle.Render("Select Branch"))
			b.WriteString("\n\n")
			b.WriteString(helpStyle.Render("  Loading branches..."))
			b.WriteString("\n\n")
			b.WriteString(helpStyle.Render("esc: cancel"))
		} else {
			b.WriteString(m.branchList.View())
			b.WriteString("\n")
			b.WriteString(helpStyle.Render("enter: checkout • /: filter • esc: cancel"))
		}
	case modeConfirmDelete:
		b.WriteString(titleStyle.Render("Confirm Delete"))
		b.WriteString("\n\n")
		b.WriteString(fmt.Sprintf("  Delete worktree: %s?\n\n", m.selectedItem.Path))
		b.WriteString(helpStyle.Render("y: yes • n: no"))
	}

	// Show errors in other modes
	if m.err != nil && m.mode != modeList {
		b.WriteString("\n")
		b.WriteString(errorStyle.Render(fmt.Sprintf("Error: %v", m.err)))
	}

	if m.message != "" {
		b.WriteString("\n")
		b.WriteString(successStyle.Render(m.message))
	}

	return b.String()
}

func (m model) updateList(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q":
		return m, tea.Quit
	case "enter":
		// Change to selected worktree directory
		if len(m.list.Items()) > 0 {
			selected := m.list.SelectedItem().(Worktree)
			m.cdPath = selected.Path
			return m, tea.Quit
		}
		return m, nil
	case "a":
		m.mode = modeAdd
		m.pathInput.SetValue("")
		m.branchInput.SetValue("")
		m.branchInput.Focus()
		m.inputFocus = 0
		m.err = nil
		m.message = ""
		return m, nil
	case "c":
		m.mode = modeCheckout
		m.err = nil
		m.message = ""
		return m, loadBranches
	case "d":
		if len(m.list.Items()) > 0 {
			selected := m.list.SelectedItem().(Worktree)
			if selected.IsMain {
				m.err = fmt.Errorf("cannot delete main worktree")
				return m, nil
			}
			m.selectedItem = selected
			m.mode = modeConfirmDelete
			m.err = nil
			m.message = ""
		}
		return m, nil
	case "r":
		m.err = nil
		m.message = ""
		return m, loadWorktrees
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) updateAdd(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.mode = modeList
		m.err = nil
		return m, nil
	case "enter":
		branch := strings.TrimSpace(m.branchInput.Value())

		if branch == "" {
			m.err = fmt.Errorf("branch name cannot be empty")
			return m, nil
		}

		// Generate path automatically
		path, err := GenerateWorktreePath(branch)
		if err != nil {
			m.err = err
			return m, nil
		}

		// Create new branch by default
		err = AddWorktree(path, branch, true)
		if err != nil {
			m.err = err
			return m, nil
		}

		m.mode = modeList
		m.message = fmt.Sprintf("Worktree created: %s", path)
		m.err = nil
		return m, loadWorktrees
	}

	// Update branch input and auto-generate path preview
	var cmd tea.Cmd
	m.branchInput, cmd = m.branchInput.Update(msg)

	// Update path preview based on branch name
	branch := strings.TrimSpace(m.branchInput.Value())
	if branch != "" {
		if path, err := GenerateWorktreePath(branch); err == nil {
			m.pathInput.SetValue(path)
		}
	} else {
		m.pathInput.SetValue("")
	}

	return m, cmd
}

func (m model) updateCheckout(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.mode = modeList
		m.err = nil
		return m, nil
	case "enter":
		if len(m.branchList.Items()) == 0 {
			return m, nil
		}

		selectedBranch := m.branchList.SelectedItem().(Branch)

		// Create worktree from existing branch
		path, err := CreateWorktreeFromBranch(selectedBranch.Name)
		if err != nil {
			m.err = err
			return m, nil
		}

		m.mode = modeList
		m.message = fmt.Sprintf("Worktree created from branch '%s': %s", selectedBranch.Name, path)
		m.err = nil
		return m, loadWorktrees
	}

	// Update branch list
	var cmd tea.Cmd
	m.branchList, cmd = m.branchList.Update(msg)
	return m, cmd
}

func (m model) updateConfirmDelete(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y":
		err := RemoveWorktree(m.selectedItem.Path, false)
		if err != nil {
			m.err = err
			m.mode = modeList
			return m, nil
		}

		m.mode = modeList
		m.message = fmt.Sprintf("Worktree removed: %s", m.selectedItem.Path)
		m.err = nil
		return m, loadWorktrees
	case "n", "esc":
		m.mode = modeList
		m.err = nil
		return m, nil
	}

	return m, nil
}
