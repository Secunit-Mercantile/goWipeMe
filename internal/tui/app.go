package tui

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mat/gowipeme/internal/cleaner"
	"github.com/mat/gowipeme/internal/wiper"
)

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
)

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

type item string

func (i item) FilterValue() string { return string(i) }

type view int

const (
	menuView view = iota
	cleanerView
	wiperMethodView
	wiperConfirmView
	wiperProgressView
	resultsView
)

type model struct {
	list            list.Model
	currentView     view
	cleanerMgr      *cleaner.CleanerManager
	dryRunResults   map[string][]string
	cleanResults    []cleaner.CleanResult
	wiper           *wiper.Wiper
	wiperMethod     wiper.WipeMethod
	wiperProgress   wiper.Progress
	wiperComplete   bool
	wiperError      error
	progressBar     progress.Model
	freeSpace       int64
	totalSpace      int64
	methodSelection int
	err             error
	quitting        bool
}

func initialModel() model {
	items := []list.Item{
		item("Clear All History"),
		item("Secure Wipe Free Space"),
		item("Quit"),
	}

	const defaultWidth = 60
	const listHeight = 14

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "goWipeMe - Privacy Tool"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	// Initialize cleaner manager
	cm := cleaner.NewCleanerManager()
	cm.AddCleaner(cleaner.NewBrowserCleaner())
	cm.AddCleaner(cleaner.NewShellCleaner())
	cm.AddCleaner(cleaner.NewCacheCleaner())
	cm.AddCleaner(cleaner.NewRecentFilesCleaner())
	cm.AddCleaner(cleaner.NewClipboardCleaner())

	// Initialize progress bar
	pb := progress.New(progress.WithDefaultGradient())

	return model{
		list:            l,
		currentView:     menuView,
		cleanerMgr:      cm,
		progressBar:     pb,
		methodSelection: 0, // Default to SinglePassZeros
	}
}

type wiperProgressMsg wiper.Progress
type wiperCompleteMsg struct{}
type wiperErrorMsg error

func (m model) Init() tea.Cmd {
	return nil
}

func startWiping(w *wiper.Wiper) tea.Cmd {
	return func() tea.Msg {
		progressChan := make(chan wiper.Progress)
		errChan := make(chan error, 1)

		// Run wiping in goroutine
		go func() {
			err := w.WipeFreeSpace(progressChan)
			if err != nil {
				errChan <- err
			}
			close(progressChan)
		}()

		// Wait for progress updates or completion
		for {
			select {
			case prog, ok := <-progressChan:
				if !ok {
					// Channel closed, wiping complete
					return wiperCompleteMsg{}
				}
				return wiperProgressMsg(prog)
			case err := <-errChan:
				return wiperErrorMsg(err)
			}
		}
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case wiperProgressMsg:
		m.wiperProgress = wiper.Progress(msg)
		return m, startWiping(m.wiper) // Continue receiving updates

	case wiperCompleteMsg:
		m.wiperComplete = true
		m.currentView = resultsView
		return m, nil

	case wiperErrorMsg:
		m.wiperError = error(msg)
		m.currentView = resultsView
		return m, nil

	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "q":
			if m.currentView == menuView {
				m.quitting = true
				return m, tea.Quit
			}
			if m.currentView == wiperProgressView {
				// Don't allow quitting during wipe
				return m, nil
			}
			// Go back to menu from other views
			m.currentView = menuView
			m.dryRunResults = nil
			m.cleanResults = nil
			m.wiperComplete = false
			m.wiperError = nil
			return m, nil

		case "up", "k":
			if m.currentView == wiperMethodView && m.methodSelection > 0 {
				m.methodSelection--
			}

		case "down", "j":
			if m.currentView == wiperMethodView && m.methodSelection < 2 {
				m.methodSelection++
			}

		case "enter":
			if m.currentView == menuView {
				i, ok := m.list.SelectedItem().(item)
				if ok {
					switch string(i) {
					case "Clear All History":
						m.currentView = cleanerView
						// Run dry-run
						results, err := m.cleanerMgr.DryRunAll()
						if err != nil {
							m.err = err
							return m, nil
						}
						m.dryRunResults = results
						return m, nil

					case "Secure Wipe Free Space":
						m.currentView = wiperMethodView
						return m, nil

					case "Quit":
						m.quitting = true
						return m, tea.Quit
					}
				}
			} else if m.currentView == cleanerView {
				// User confirmed, run cleaning
				m.cleanResults = m.cleanerMgr.CleanAll()
				m.currentView = resultsView
				return m, nil
			} else if m.currentView == wiperMethodView {
				// User selected a wipe method
				m.wiperMethod = wiper.WipeMethod(m.methodSelection)

				// Get home directory as default volume
				homeDir, _ := wiper.GetHomeDir()
				w, err := wiper.NewWiper(homeDir, m.wiperMethod)
				if err != nil {
					m.err = err
					m.currentView = menuView
					return m, nil
				}
				m.wiper = w

				// Get volume info
				totalSpace, freeSpace, err := w.GetVolumeInfo()
				if err != nil {
					m.err = err
					m.currentView = menuView
					return m, nil
				}
				m.totalSpace = totalSpace
				m.freeSpace = freeSpace

				m.currentView = wiperConfirmView
				return m, nil
			} else if m.currentView == wiperConfirmView {
				// Start wiping
				m.currentView = wiperProgressView
				return m, startWiping(m.wiper)
			} else if m.currentView == resultsView {
				// Go back to menu
				m.currentView = menuView
				m.dryRunResults = nil
				m.cleanResults = nil
				m.wiperComplete = false
				m.wiperError = nil
				return m, nil
			}
		}
	}

	var cmd tea.Cmd
	if m.currentView == menuView {
		m.list, cmd = m.list.Update(msg)
	}
	return m, cmd
}

func (m model) View() string {
	if m.quitting {
		return "Thanks for using goWipeMe!\n"
	}

	switch m.currentView {
	case menuView:
		return "\n" + m.list.View()

	case cleanerView:
		return m.renderCleanerView()

	case wiperMethodView:
		return m.renderWiperMethodView()

	case wiperConfirmView:
		return m.renderWiperConfirmView()

	case wiperProgressView:
		return m.renderWiperProgressView()

	case resultsView:
		return m.renderResultsView()

	default:
		return "Unknown view"
	}
}

func (m model) renderCleanerView() string {
	var s strings.Builder

	s.WriteString("\n  🧹 Clear All History - Dry Run\n\n")

	if m.err != nil {
		s.WriteString(fmt.Sprintf("  Error: %v\n\n", m.err))
		s.WriteString("  Press 'q' to go back\n")
		return s.String()
	}

	if len(m.dryRunResults) == 0 {
		s.WriteString("  ✓ Nothing to clean!\n\n")
		s.WriteString("  Press 'q' to go back\n")
		return s.String()
	}

	totalItems := 0
	for cleanerName, items := range m.dryRunResults {
		totalItems += len(items)
		s.WriteString(fmt.Sprintf("  %s (%d items):\n", cleanerName, len(items)))
		for _, item := range items {
			s.WriteString(fmt.Sprintf("    • %s\n", item))
		}
		s.WriteString("\n")
	}

	s.WriteString(fmt.Sprintf("  Total items to clean: %d\n\n", totalItems))
	s.WriteString("  ⚠️  WARNING: This action cannot be undone!\n\n")
	s.WriteString("  Press ENTER to confirm and clean\n")
	s.WriteString("  Press 'q' to cancel\n")

	return s.String()
}

func (m model) renderResultsView() string {
	var s strings.Builder

	// Check if this is wiper results or cleaner results
	if m.wiperComplete || m.wiperError != nil {
		s.WriteString("\n  ✨ Disk Wiping Complete\n\n")

		if m.wiperError != nil {
			s.WriteString(fmt.Sprintf("  ✗ Error: %v\n", m.wiperError))
		} else {
			s.WriteString(fmt.Sprintf("  ✓ Successfully wiped free space using %s\n", m.wiperMethod.String()))
			s.WriteString(fmt.Sprintf("  ✓ Wiped: %s\n", wiper.FormatBytes(m.wiperProgress.BytesWritten)))
			s.WriteString(fmt.Sprintf("  ✓ Time taken: %s\n", m.wiperProgress.TimeElapsed.Round(time.Second)))
		}
	} else {
		s.WriteString("\n  ✨ Cleaning Complete\n\n")

		for _, result := range m.cleanResults {
			if result.Error != nil {
				s.WriteString(fmt.Sprintf("  ✗ %s: %v\n", result.CleanerName, result.Error))
			} else {
				s.WriteString(fmt.Sprintf("  ✓ %s: cleaned %d items\n", result.CleanerName, result.ItemsCleaned))
			}
		}
	}

	s.WriteString("\n  Press ENTER or 'q' to return to menu\n")

	return s.String()
}

func (m model) renderWiperMethodView() string {
	var s strings.Builder

	s.WriteString("\n  💾 Select Wipe Method\n\n")

	methods := []wiper.WipeMethod{
		wiper.SinglePassZeros,
		wiper.DoD522022M,
		wiper.Gutmann,
	}

	for i, method := range methods {
		cursor := "  "
		if i == m.methodSelection {
			cursor = "> "
		}
		s.WriteString(fmt.Sprintf("  %s%d. %s\n", cursor, i+1, method.String()))
		s.WriteString(fmt.Sprintf("     %s\n\n", method.Description()))
	}

	s.WriteString("  Use arrow keys or j/k to navigate\n")
	s.WriteString("  Press ENTER to select\n")
	s.WriteString("  Press 'q' to go back\n")

	return s.String()
}

func (m model) renderWiperConfirmView() string {
	var s strings.Builder

	s.WriteString("\n  💾 Secure Wipe Free Space - Confirmation\n\n")

	s.WriteString(fmt.Sprintf("  Volume: %s\n", m.wiper.VolumePath))
	s.WriteString(fmt.Sprintf("  Total Space: %s\n", wiper.FormatBytes(m.totalSpace)))
	s.WriteString(fmt.Sprintf("  Free Space: %s\n", wiper.FormatBytes(m.freeSpace)))
	s.WriteString(fmt.Sprintf("  Method: %s\n", m.wiperMethod.String()))
	s.WriteString(fmt.Sprintf("  Description: %s\n\n", m.wiperMethod.Description()))

	s.WriteString("  ⚠️  WARNING: This operation will:\n")
	s.WriteString("     • Fill all free space on the volume\n")
	s.WriteString("     • Take a significant amount of time\n")
	s.WriteString("     • Cannot be interrupted once started\n\n")

	if m.wiperMethod == wiper.DoD522022M {
		s.WriteString("  ℹ️  Note: DoD method will make 3 passes over the free space\n\n")
	} else if m.wiperMethod == wiper.Gutmann {
		s.WriteString("  ℹ️  Note: Gutmann method will make 35 passes (very slow!)\n\n")
	}

	s.WriteString("  Press ENTER to start wiping\n")
	s.WriteString("  Press 'q' to cancel\n")

	return s.String()
}

func (m model) renderWiperProgressView() string {
	var s strings.Builder

	s.WriteString("\n  💾 Wiping Free Space...\n\n")

	if m.wiperProgress.TotalBytes > 0 {
		percentage := m.wiperProgress.Percentage()
		s.WriteString(fmt.Sprintf("  Progress: %.1f%%\n\n", percentage))

		// Progress bar
		barWidth := 50
		filled := int(percentage / 100.0 * float64(barWidth))
		bar := strings.Repeat("█", filled) + strings.Repeat("░", barWidth-filled)
		s.WriteString(fmt.Sprintf("  [%s]\n\n", bar))

		s.WriteString(fmt.Sprintf("  Current Pass: %d/%d\n", m.wiperProgress.CurrentPass, m.wiperProgress.TotalPasses))
		s.WriteString(fmt.Sprintf("  Method: %s\n", m.wiperProgress.CurrentMethod))
		s.WriteString(fmt.Sprintf("  Written: %s / %s\n",
			wiper.FormatBytes(m.wiperProgress.BytesWritten),
			wiper.FormatBytes(m.wiperProgress.TotalBytes)))
		s.WriteString(fmt.Sprintf("  Elapsed: %s\n", m.wiperProgress.TimeElapsed.Round(time.Second)))

		if m.wiperProgress.EstimatedTime > 0 {
			s.WriteString(fmt.Sprintf("  Estimated Remaining: %s\n", m.wiperProgress.EstimatedTime.Round(time.Second)))
		}
	} else {
		s.WriteString("  Initializing...\n")
	}

	s.WriteString("\n  ⚠️  Please do not interrupt this process!\n")
	s.WriteString("  (Ctrl+C is disabled during wiping)\n")

	return s.String()
}

// Run starts the TUI application
func Run() error {
	p := tea.NewProgram(initialModel())
	_, err := p.Run()
	return err
}
