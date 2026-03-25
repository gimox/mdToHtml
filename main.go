package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

const (
	inputDir  = "input_md"
	outputDir = "output_html"
)

var (
	successStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("82")).Bold(true)
	fileStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("214"))
	infoStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("39"))
	warnStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("201")).Bold(true)
)

// --- MESSAGGI ---
type fileDoneMsg string
type errMsg error

// --- LOGICA CORE ---
func convertMdToHtml(fileName string) error {
	time.Sleep(100 * time.Millisecond)
	md := goldmark.New(goldmark.WithExtensions(extension.GFM))
	source, err := os.ReadFile(filepath.Join(inputDir, fileName))
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	if err := md.Convert(source, &buf); err != nil {
		return err
	}
	baseName := strings.TrimSuffix(fileName, filepath.Ext(fileName))
	return os.WriteFile(filepath.Join(outputDir, baseName+".html"), buf.Bytes(), 0644)
}

func cleanOutputFolder() error {
	_ = os.RemoveAll(outputDir)
	return os.MkdirAll(outputDir, 0755)
}

func openDirectory(path string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("explorer", path)
	case "darwin":
		cmd = exec.Command("open", path)
	default:
		cmd = exec.Command("xdg-open", path)
	}
	return cmd.Start()
}

// --- TUI MODEL ---
type item struct{ title, desc string }

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type model struct {
	list           list.Model
	spinner        spinner.Model
	state          string // MENU, PICKER, CONFIRM_CLEAN, LOADING, DONE
	filesToProcess []string
	processedCount int
	err            error
	width, height  int
}

func initialModel() model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return model{
		list:    createMainMenu(),
		spinner: s,
		state:   "MENU",
	}
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

		// Gestione specifica durante la conferma pulizia
		if m.state == "CONFIRM_CLEAN" {
			switch strings.ToLower(msg.String()) {
			case "y":
				_ = cleanOutputFolder()
				m.state = "LOADING"
				return m, tea.Batch(m.spinner.Tick, m.nextFileCmd())
			case "n":
				m.state = "LOADING"
				return m, tea.Batch(m.spinner.Tick, m.nextFileCmd())
			case "esc":
				m.state = "MENU"
				m.list = createMainMenu()
				m.resetListSize()
				return m, nil
			}
			return m, nil
		}

		if m.state == "LOADING" {
			return m, nil
		}

		switch msg.String() {
		case "esc":
			if m.list.FilterState() == list.Filtering {
				break
			}
			if m.state == "PICKER" || m.state == "DONE" {
				m.state = "MENU"
				m.list = createMainMenu()
				m.resetListSize()
				m.processedCount = 0
				return m, nil
			}
		case "q":
			if m.list.FilterState() != list.Filtering {
				return m, tea.Quit
			}
		case "enter":
			selected, ok := m.list.SelectedItem().(item)
			if !ok {
				return m, nil
			}

			if m.state == "MENU" {
				if selected.title == "Processa tutto" {
					files, _ := filepath.Glob(filepath.Join(inputDir, "*.md"))
					if len(files) == 0 {
						return m, nil
					}
					m.filesToProcess = nil
					for _, f := range files {
						m.filesToProcess = append(m.filesToProcess, filepath.Base(f))
					}
					m.state = "CONFIRM_CLEAN"
					return m, nil
				} else {
					m.state = "PICKER"
					m.list = createFileList()
					m.resetListSize()
					return m, nil
				}
			} else if m.state == "PICKER" {
				m.filesToProcess = []string{selected.title}
				m.state = "CONFIRM_CLEAN"
				return m, nil
			}
		}

	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.resetListSize()

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case fileDoneMsg:
		m.processedCount++
		if m.processedCount < len(m.filesToProcess) {
			return m, m.nextFileCmd()
		}
		_ = openDirectory(outputDir)
		m.state = "DONE"
		return m, nil

	case errMsg:
		m.err = msg
		m.state = "DONE"
		return m, nil
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m *model) resetListSize() {
	h, v := lipgloss.NewStyle().Margin(1, 2).GetFrameSize()
	m.list.SetSize(m.width-h, m.height-v)
}

func (m model) nextFileCmd() tea.Cmd {
	return func() tea.Msg {
		err := convertMdToHtml(m.filesToProcess[m.processedCount])
		if err != nil {
			return errMsg(err)
		}
		return fileDoneMsg(m.filesToProcess[m.processedCount])
	}
}

func (m model) View() string {
	switch m.state {
	case "CONFIRM_CLEAN":
		return fmt.Sprintf("\n\n  %s\n\n  Vuoi svuotare la cartella /%s prima di iniziare?\n  (I file esistenti verranno eliminati definitivamente)\n\n  %s  %s  %s",
			warnStyle.Render("⚠️  ATTENZIONE: PULIZIA OUTPUT"),
			outputDir,
			successStyle.Render("[Y] Sì, pulisci"),
			infoStyle.Render("[N] No, mantieni"),
			fileStyle.Render("[ESC] Annulla"),
		)
	case "LOADING":
		cur := m.filesToProcess[m.processedCount]
		return fmt.Sprintf("\n\n  %s Elaborazione %d/%d: %s...\n\n",
			m.spinner.View(), m.processedCount+1, len(m.filesToProcess), fileStyle.Render(cur))
	case "DONE":
		res := successStyle.Render("✨ CONVERSIONE COMPLETATA!")
		return fmt.Sprintf("\n\n  %s\n  File processati: %d\n  %s\n\n  Premi 'esc' per tornare al menu o 'q' per uscire",
			res, m.processedCount, infoStyle.Render("📂 Cartella aperta."))
	default:
		return lipgloss.NewStyle().Margin(1, 2).Render(m.list.View())
	}
}

// --- HELPERS ---

func createMainMenu() list.Model {
	items := []list.Item{
		item{title: "Scegli file singolo", desc: "Seleziona un .md specifico"},
		item{title: "Processa tutto", desc: "Converti tutti i file della cartella"},
	}
	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "RAG Doc Converter"
	l.KeyMap.Quit.SetKeys("q")
	return l
}

func createFileList() list.Model {
	files, _ := filepath.Glob(filepath.Join(inputDir, "*.md"))
	var items []list.Item
	for _, f := range files {
		items = append(items, item{title: filepath.Base(f), desc: "Markdown"})
	}
	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Seleziona File (ESC per tornare)"
	l.KeyMap.Quit.SetKeys("q")
	return l
}

func main() {
	_ = os.MkdirAll(inputDir, 0755)
	_ = os.MkdirAll(outputDir, 0755)
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Errore: %v", err)
		os.Exit(1)
	}
}
