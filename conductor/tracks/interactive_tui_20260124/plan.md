# Implementation Plan - Interactive TUI Console

## Phase 1: Scaffolding
- [x] **Task: Install Dependencies**
    - `go get github.com/charmbracelet/bubbletea`
    - `go get github.com/charmbracelet/lipgloss`
    - `go get github.com/charmbracelet/bubbles`
- [x] **Task: CLI Command**
    - Create `internal/cli/console.go`.
    - Implement the `spectre console` command.
- [x] **Task: Basic Model**
    - Create `internal/tui/model.go` and `internal/tui/app.go`.
    - Implement the `tea.Model` interface.

## Phase 2: Navigation & Dashboard
- [x] **Task: Home View**
    - Implement a polished dashboard with ASCII art and system stats.
- [x] **Task: Sidebar/Tabs**
    - Add a way to switch between "Home", "Cases", and "Run".

## Phase 3: Case Management
- [x] **Task: Case List**
    - Use `bubbles/list` to show all investigation cases.
- [ ] **Task: Case Details**
    - Use `bubbles/table` to display entities when a case is selected.

## Phase 4: Interactive Collection
- [ ] **Task: Execution Engine**
    - Implement a way to trigger `collector.Run` and show a progress spinner in the TUI.

## Phase 5: Polishing
- [ ] **Task: Styling**
    - Use `lipgloss` to add the SPECTRE color theme (Blue/Purple/Dark).
- [ ] **Task: Refinement**
    - Ensure smooth transitions and error handling.
