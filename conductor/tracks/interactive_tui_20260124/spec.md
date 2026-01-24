# Specification - Interactive TUI Console

## Overview
SPECTRE needs a central "Command Center" in the terminal. The TUI will allow users to manage cases, explore entities, and trigger collectors without leaving the terminal or typing long commands.

## Requirements

### 1. Framework
- **Primary:** `github.com/charmbracelet/bubbletea` (TUI engine)
- **Styling:** `github.com/charmbracelet/lipgloss`
- **Components:** `github.com/charmbracelet/bubbles` (Lists, Tables, Text Inputs)

### 2. Main Views
- **Home/Dashboard:** ASCII art banner, summary of total cases/entities, and recent activity.
- **Case Explorer:** A list of all cases. Selecting one opens the Case Detail view.
- **Case Detail:**
    - Table of entities (Value, Type, Source).
    - Table of relationships.
    - Log view for recent evidence.
- **Collector Runner:** A simple form to input a target and select a collector (DNS, WHOIS, etc.).

### 3. Navigation
- Tab-based or Sidebar-based navigation.
- Hotkeys: `q` to quit, `esc` to go back, `ctrl+c` for emergency exit.

## Success Criteria
- Running `spectre console` launches a full-screen interactive interface.
- Users can browse cases and entities using arrow keys.
- Collectors can be triggered from within the TUI.
