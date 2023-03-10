package tui

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/atotto/clipboard"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pkg/browser"
	"github.com/pomdtr/sunbeam/schemas"
	"github.com/pomdtr/sunbeam/utils"
)

type PopPageMsg struct{}

type pushMsg struct {
	container Page
}

type Page interface {
	Init() tea.Cmd
	Update(tea.Msg) (Page, tea.Cmd)
	View() string
	SetSize(width, height int)
}

type SunbeamOptions struct {
	MaxHeight int
	Padding   int
}

type Model struct {
	width, height int
	options       SunbeamOptions
	exitCmd       *exec.Cmd

	root  Page
	pages []Page
	form  *Form

	hidden bool
}

func NewModel(root Page, options SunbeamOptions) *Model {
	return &Model{root: root, options: options}
}

func (m *Model) SetRoot(root Page) {
	m.root = root
}

func (m *Model) Init() tea.Cmd {
	return m.root.Init()
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEscape:
			fmt.Sprintln("Escape")
		case tea.KeyCtrlC:
			m.hidden = true
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.SetSize(msg.Width, msg.Height)
		return m, nil
	case OpenMsg:
		err := browser.OpenURL(msg.Action.Target)
		if err != nil {
			return m, func() tea.Msg {
				return err
			}
		}

		m.hidden = true
		return m, tea.Quit
	case CopyTextMsg:
		err := clipboard.WriteAll(msg.Action.Text)
		if err != nil {
			return m, func() tea.Msg {
				return fmt.Errorf("failed to copy text to clipboard: %s", err)
			}
		}

		m.hidden = true
		return m, tea.Quit
	case PushPageMsg:
		cmd := m.Push(NewCommandRunner(func(query string) ([]byte, error) {
			return os.ReadFile(msg.Action.Path)
		}))

		return m, cmd
	case RunCommandMsg:
		if len(msg.Action.Inputs) > 0 {
			formItems := make([]FormItem, len(msg.Action.Inputs))
			for i, input := range msg.Action.Inputs {
				item, err := NewFormItem(input)
				if err != nil {
					return m, func() tea.Msg {
						return fmt.Errorf("failed to create form input: %s", err)
					}
				}

				formItems[i] = item
			}

			form := NewForm(formItems, func(values map[string]string) tea.Cmd {
				command := make([]string, len(msg.Action.Command))
				for i, arg := range msg.Action.Command {
					for key, value := range values {
						arg = strings.ReplaceAll(arg, fmt.Sprintf("${input:%s}", key), value)
					}
					command[i] = arg
				}

				return func() tea.Msg {
					return RunCommandMsg{Action: schemas.Action{
						Command:   command,
						OnSuccess: msg.Action.OnSuccess,
					}}
				}
			})
			m.form = form
			form.SetSize(m.pageWidth(), m.pageHeight())
			return m, form.Init()
		}

		m.form = nil

		name, args := utils.SplitCommand(msg.Action.Command)

		if msg.Action.OnSuccess == "" {
			m.exitCmd = exec.Command(name, args...)
			m.hidden = true
			return m, tea.Quit
		}
		if msg.Action.OnSuccess == "push" {
			return m, m.Push(NewCommandRunner(func(query string) ([]byte, error) {
				cmd := exec.Command(name, args...)
				cmd.Stdin = strings.NewReader(query)
				output, err := cmd.Output()
				if err != nil {
					if err, ok := err.(*exec.ExitError); ok {
						return nil, fmt.Errorf("command exit with code %d: %s", err.ExitCode(), err.Stderr)
					}
					return nil, err
				}

				return output, nil
			}))
		}

		return m, func() tea.Msg {
			cmd := exec.Command(name, args...)
			output, err := cmd.Output()
			if err != nil {
				return fmt.Errorf("failed to run command: %s", err)
			}

			switch msg.Action.OnSuccess {
			case "copy":
				return CopyTextMsg{Action: schemas.Action{
					Text: string(output),
				}}
			case "open":
				return OpenMsg{Action: schemas.Action{
					Target: string(output),
				}}
			case "reload":
				return ReloadPageMsg{}
			default:
				return fmt.Errorf("unknown onSuccess action: %s", msg.Action.OnSuccess)
			}
		}
	case PopPageMsg:
		if m.form != nil {
			m.form = nil
			return m, nil
		}

		if len(m.pages) == 0 {
			m.hidden = true
			return m, tea.Quit
		} else {
			m.Pop()
			return m, nil
		}
	case error:
		detail := NewDetail("Error", msg.Error, []Action{
			{
				Title: "Copy error",
				Cmd: func() tea.Msg {
					return CopyTextMsg{schemas.Action{
						Text: msg.Error(),
					}}
				},
			},
			{
				Title: "Reload Page",
				Cmd: func() tea.Msg {
					return ReloadPageMsg{}
				},
			},
		})
		detail.SetSize(m.pageWidth(), m.pageHeight())

		if len(m.pages) == 0 {
			m.root = detail
		} else {
			m.pages[len(m.pages)-1] = detail
		}

		return m, detail.Init()
	}

	// Update the current page
	var cmd tea.Cmd

	if m.form != nil {
		m.form, cmd = m.form.Update(msg)
	} else if len(m.pages) == 0 {
		m.root, cmd = m.root.Update(msg)
	} else {
		currentPageIdx := len(m.pages) - 1
		m.pages[currentPageIdx], cmd = m.pages[currentPageIdx].Update(msg)
	}

	return m, cmd
}

func (m *Model) View() string {
	if m.hidden {
		return ""
	}

	if m.form != nil {
		return m.form.View()
	}

	var pageView string

	if len(m.pages) > 0 {
		currentPage := m.pages[len(m.pages)-1]
		pageView = currentPage.View()
	} else {
		pageView = m.root.View()
	}

	return lipgloss.NewStyle().Padding(m.options.Padding).Render(pageView)
}

func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height

	m.root.SetSize(m.pageWidth(), m.pageHeight())
	for _, page := range m.pages {
		page.SetSize(m.pageWidth(), m.pageHeight())
	}
}

func (m *Model) pageWidth() int {
	return m.width - 2*m.options.Padding
}

func (m *Model) pageHeight() int {
	if m.options.MaxHeight > 0 {
		return utils.Min(m.options.MaxHeight, m.height) - 2*m.options.Padding
	}
	return m.height - 2*m.options.Padding
}

func (m *Model) Push(page Page) tea.Cmd {
	page.SetSize(m.pageWidth(), m.pageHeight())
	m.pages = append(m.pages, page)
	return page.Init()
}

func (m *Model) Pop() {
	if len(m.pages) > 0 {
		m.pages = m.pages[:len(m.pages)-1]
	}
}

func (m *Model) Draw() (err error) {
	// Background detection before we start the program
	lipgloss.SetHasDarkBackground(lipgloss.HasDarkBackground())

	err = os.Setenv("SUNBEAM_RUNNER", "true")
	if err != nil {
		return fmt.Errorf("failed to set SUNBEAM_RUNNER env var: %s", err)
	}

	var p *tea.Program
	if m.options.MaxHeight == 0 {
		p = tea.NewProgram(m, tea.WithAltScreen())
	} else {
		p = tea.NewProgram(m)
	}

	res, err := p.Run()
	if err != nil {
		return err
	}

	model, ok := res.(*Model)
	if !ok {
		return fmt.Errorf("could not convert res back to *Model")
	}

	if model.exitCmd != nil {
		model.exitCmd.Stdin = os.Stdin
		model.exitCmd.Stdout = os.Stdout
		model.exitCmd.Stderr = os.Stderr

		return model.exitCmd.Run()
	}

	return nil
}
