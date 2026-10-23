package main

import (
	"errors"
	"strings"
	"testing"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"
)

func TestSubstringFilter(t *testing.T) {
	tests := []struct {
		name     string
		term     string
		targets  []string
		expected int // expected number of matches
	}{
		{
			name:     "simple substring match",
			term:     "bug",
			targets:  []string{"bugfix", "feature", "debug"},
			expected: 2, // bugfix and debug
		},
		{
			name:     "case insensitive match",
			term:     "BUG",
			targets:  []string{"bugfix", "feature", "debug"},
			expected: 2,
		},
		{
			name:     "no matches",
			term:     "xyz",
			targets:  []string{"bugfix", "feature", "debug"},
			expected: 0,
		},
		{
			name:     "dash in term",
			term:     "-bugg",
			targets:  []string{"origin/ppoluektov/KR-7865-better-customer-facing-error-messages-in-direct-avs", "feature-buggy", "main"},
			expected: 1, // only feature-buggy contains "-bugg"
		},
		{
			name:     "partial match",
			term:     "feat",
			targets:  []string{"feature", "bugfix", "feature-auth"},
			expected: 2, // feature and feature-auth
		},
		{
			name:     "empty term",
			term:     "",
			targets:  []string{"feature", "bugfix"},
			expected: 2, // empty string matches everything
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ranks := substringFilter(tt.term, tt.targets)
			if len(ranks) != tt.expected {
				t.Errorf("substringFilter(%q, %v) returned %d matches, expected %d",
					tt.term, tt.targets, len(ranks), tt.expected)
				for _, rank := range ranks {
					t.Logf("  Matched: %s (index %d)", tt.targets[rank.Index], rank.Index)
				}
			}

			// Verify that matched indexes are set correctly
			for _, rank := range ranks {
				if len(rank.MatchedIndexes) == 0 && tt.term != "" {
					t.Errorf("Rank for %q has no matched indexes", tt.targets[rank.Index])
				}
			}
		})
	}
}


func key(code rune) tea.KeyPressMsg {
	k := tea.Key{Code: code}
	if code >= 32 && code < 127 {
		k.Text = string(code)
	}
	return tea.KeyPressMsg(k)
}

func send(t *testing.T, m model, msg tea.Msg) (model, tea.Cmd) {
	t.Helper()
	next, cmd := m.Update(msg)
	nm, ok := next.(model)
	if !ok {
		t.Fatalf("Update returned %T, want model", next)
	}
	return nm, cmd
}

func isQuit(cmd tea.Cmd) bool {
	if cmd == nil {
		return false
	}
	_, ok := cmd().(tea.QuitMsg)
	return ok
}

func loadedModel(t *testing.T) model {
	t.Helper()
	m := initialModel()
	m, _ = send(t, m, tea.WindowSizeMsg{Width: 100, Height: 30})
	m, _ = send(t, m, worktreesLoadedMsg{
		{Path: "/repo", Branch: "master", Commit: "abc1234", IsMain: true},
		{Path: "/repo/.worktrees/feat", Branch: "feat", Commit: "def5678"},
	})
	return m
}

func TestViewUsesAltScreen(t *testing.T) {
	v := initialModel().View()
	if !v.AltScreen {
		t.Error("View().AltScreen = false, want true")
	}
}

func TestKeyPressStrings(t *testing.T) {
	if got := key('q').String(); got != "q" {
		t.Errorf("q key String() = %q", got)
	}
	if got := key(tea.KeyEnter).String(); got != "enter" {
		t.Errorf("enter key String() = %q", got)
	}
	if got := key(tea.KeyEscape).String(); got != "esc" {
		t.Errorf("esc key String() = %q", got)
	}
	ctrlC := tea.KeyPressMsg(tea.Key{Code: 'c', Mod: tea.ModCtrl})
	if got := ctrlC.String(); got != "ctrl+c" {
		t.Errorf("ctrl+c String() = %q", got)
	}
}

func TestQuitKeys(t *testing.T) {
	m := loadedModel(t)
	if _, cmd := send(t, m, key('q')); !isQuit(cmd) {
		t.Error("q in list mode should quit")
	}
	ctrlC := tea.KeyPressMsg(tea.Key{Code: 'c', Mod: tea.ModCtrl})
	m.mode = modeAdd
	if _, cmd := send(t, m, ctrlC); !isQuit(cmd) {
		t.Error("ctrl+c should quit from any mode")
	}
}

func TestEnterSelectsWorktree(t *testing.T) {
	m := loadedModel(t)
	m, _ = send(t, m, key('j'))
	m, cmd := send(t, m, key(tea.KeyEnter))
	if !isQuit(cmd) {
		t.Error("enter should quit")
	}
	if m.cdPath != "/repo/.worktrees/feat" {
		t.Errorf("cdPath = %q, want the second worktree", m.cdPath)
	}
}

func TestDeleteFlow(t *testing.T) {
	m := loadedModel(t)

	// Main worktree cannot be deleted.
	m, _ = send(t, m, key('d'))
	if m.mode != modeList || m.err == nil {
		t.Fatalf("deleting main: mode=%v err=%v, want list mode with error", m.mode, m.err)
	}

	m, _ = send(t, m, key('j'))
	m, _ = send(t, m, key('d'))
	if m.mode != modeConfirmDelete || m.selectedItem.Path != "/repo/.worktrees/feat" {
		t.Fatalf("mode=%v selected=%q, want confirm for feat", m.mode, m.selectedItem.Path)
	}
	if out := ansi.Strip(m.render()); !strings.Contains(out, "Delete worktree: /repo/.worktrees/feat?") {
		t.Errorf("confirm view missing prompt: %q", out)
	}

	m, _ = send(t, m, key('n'))
	if m.mode != modeList {
		t.Errorf("n should return to list mode, got %v", m.mode)
	}
}

func TestAddModeFlow(t *testing.T) {
	m := loadedModel(t)
	m, _ = send(t, m, key('a'))
	if m.mode != modeAdd {
		t.Fatalf("mode = %v, want modeAdd", m.mode)
	}

	for _, r := range "feat-x" {
		m, _ = send(t, m, key(r))
	}
	if got := m.branchInput.Value(); got != "feat-x" {
		t.Errorf("branchInput = %q, want feat-x", got)
	}
	out := ansi.Strip(m.render())
	if !strings.Contains(out, "Add New Worktree") || !strings.Contains(out, "feat-x") {
		t.Errorf("add view missing title or typed text: %q", out)
	}

	m, _ = send(t, m, key(tea.KeyEscape))
	if m.mode != modeList {
		t.Errorf("esc should return to list mode, got %v", m.mode)
	}
}

func TestAddEmptyBranchErrors(t *testing.T) {
	m := loadedModel(t)
	m, _ = send(t, m, key('a'))
	m, _ = send(t, m, key(tea.KeyEnter))
	if m.err == nil || m.mode != modeAdd {
		t.Errorf("empty branch: mode=%v err=%v, want modeAdd with error", m.mode, m.err)
	}
}

func TestListViewContent(t *testing.T) {
	out := ansi.Strip(loadedModel(t).render())
	for _, want := range []string{"Git Worktrees", "master", "feat", "enter: cd to worktree"} {
		if !strings.Contains(out, want) {
			t.Errorf("list view missing %q:\n%s", want, out)
		}
	}
}

func TestEmptyAndErrorViews(t *testing.T) {
	m := initialModel()
	out := ansi.Strip(m.render())
	if !strings.Contains(out, "No worktrees found.") {
		t.Errorf("empty view: %q", out)
	}

	m, _ = send(t, m, errMsg(errors.New("not a git repository")))
	out = ansi.Strip(m.render())
	if !strings.Contains(out, "not a git repository") || !strings.Contains(out, "retry") {
		t.Errorf("error view: %q", out)
	}
}

func TestCheckoutBranchFilter(t *testing.T) {
	m := loadedModel(t)
	m, _ = send(t, m, key('c'))
	if m.mode != modeCheckout {
		t.Fatalf("mode = %v, want modeCheckout", m.mode)
	}
	m, _ = send(t, m, branchesLoadedMsg{
		{Name: "main"}, {Name: "feature-auth"}, {Name: "bugfix"},
	})
	if got := len(m.branchList.Items()); got != 3 {
		t.Fatalf("branch items = %d, want 3", got)
	}

	m, _ = send(t, m, key('/'))
	if m.branchList.FilterState() != list.Filtering {
		t.Fatalf("filter state = %v, want Filtering", m.branchList.FilterState())
	}
	for _, r := range "auth" {
		m, _ = send(t, m, key(r))
	}
	if got := m.branchList.FilterValue(); got != "auth" {
		t.Errorf("filter value = %q, want auth", got)
	}
	m, _ = send(t, m, key(tea.KeyEscape))
	if m.mode != modeCheckout {
		t.Errorf("esc while filtering should clear filter, not leave checkout (mode=%v)", m.mode)
	}
}
