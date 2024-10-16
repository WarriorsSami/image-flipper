package main

import (
	"context"
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	imgproc "image_utils"
	"os"
	"strings"
)

const maxWidth = 100

var (
	red    = lipgloss.AdaptiveColor{Light: "#FE5F86", Dark: "#FE5F86"}
	indigo = lipgloss.AdaptiveColor{Light: "#5A56E0", Dark: "#7571F9"}
	green  = lipgloss.AdaptiveColor{Light: "#02BA84", Dark: "#02BF87"}
)

type Styles struct {
	Base,
	HeaderText,
	Status,
	StatusHeader,
	Highlight,
	ErrorHeaderText,
	Help lipgloss.Style
}

func NewStyles(lg *lipgloss.Renderer) *Styles {
	s := Styles{}
	s.Base = lg.NewStyle().
		Padding(1, 4, 0, 1)
	s.HeaderText = lg.NewStyle().
		Foreground(indigo).
		Bold(true).
		Padding(0, 1, 0, 2)
	s.Status = lg.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(indigo).
		PaddingLeft(1).
		MarginTop(1)
	s.StatusHeader = lg.NewStyle().
		Foreground(green).
		Bold(true)
	s.Highlight = lg.NewStyle().
		Foreground(lipgloss.Color("212"))
	s.ErrorHeaderText = s.HeaderText.
		Foreground(red)
	s.Help = lg.NewStyle().
		Foreground(lipgloss.Color("240"))
	return &s
}

type Model struct {
	lg     *lipgloss.Renderer
	styles *Styles
	form   *huh.Form
	width  int
	// TODO: Fix the loading spinner
	loading bool
	spinner spinner.Model
}

func NewModel() Model {
	m := Model{width: maxWidth}
	m.lg = lipgloss.DefaultRenderer()
	m.styles = NewStyles(m.lg)
	m.loading = false
	m.spinner = spinner.New()
	m.spinner.Spinner = spinner.Dot

	homeDir, _ := os.UserHomeDir()

	m.form = huh.NewForm(
		huh.NewGroup(
			huh.NewFilePicker().
				CurrentDirectory(homeDir).
				DirAllowed(true).
				FileAllowed(false).
				Key("inputFolderPath").
				Title("Select the folder with images").
				Description("This is where the images you want to flip are located").
				Validate(checkIfFolderExists),

			huh.NewFilePicker().
				CurrentDirectory(homeDir).
				DirAllowed(true).
				FileAllowed(false).
				Key("outputFolderPath").
				Title("Select the output folder").
				Description("This is where the flipped images will be saved").
				Validate(checkIfFolderExists),

			huh.NewSelect[string]().
				Key("flipMode").
				Options(huh.NewOptions("horizontal", "vertical", "both")...).
				Title("Choose your flipping mode").
				Description("This will determine how the images are flipped"),

			huh.NewConfirm().
				Key("done").
				Title("All done?").
				Validate(func(v bool) error {
					if !v {
						return fmt.Errorf("welp, finish up then")
					}
					return nil
				}).
				Affirmative("Yep").
				Negative("Wait, no"),
		),
	).
		WithWidth(70).
		WithShowHelp(false).
		WithShowErrors(false)

	return m
}

func (m Model) Init() tea.Cmd {
	return m.form.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = min(msg.Width, maxWidth) - m.styles.Base.GetHorizontalFrameSize()
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	var cmds []tea.Cmd

	// Process the form
	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
		cmds = append(cmds, cmd)
	}

	// Process the spinner
	m.spinner, cmd = m.spinner.Update(msg)
	cmds = append(cmds, cmd)

	// If the form is completed, flip the images
	if m.form.State == huh.StateCompleted {
		cmds = append(cmds, m.spinner.Tick)
		m.loading = true

		if err := m.flipImages(); err != nil {
			return m, tea.Batch(cmds...)
		}

		cmds = append(cmds, tea.Quit)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	s := m.styles

	switch m.form.State {
	case huh.StateCompleted:
		outcomeView := m.getFlipOutcomeView()
		outcomeView = s.Highlight.Render(outcomeView)
		var b strings.Builder
		fmt.Fprintf(&b, "Images have been flipped successfully!\n\n%s", outcomeView)
		return s.Status.Margin(0, 1).Padding(1, 2).Width(48).Render(b.String()) + "\n\n"
	default:
		// Form (left side)
		v := strings.TrimSuffix(m.form.View(), "\n\n")
		form := m.lg.NewStyle().Margin(1, 0).Render(v)

		// Status (right side)
		var status string
		{
			var (
				buildInfo        = m.getBuildInfoView()
				inputFolderPath  = ""
				outputFolderPath = ""
				flipMode         = imgproc.FlipHorizontal
			)

			if m.form.GetString("inputFolderPath") != "" {
				inputFolderPath = "Input Folder: " + m.form.GetString("inputFolderPath")
			}

			if m.form.GetString("outputFolderPath") != "" {
				outputFolderPath = "Output Folder: " + m.form.GetString("outputFolderPath")
			}

			if m.form.GetString("flipMode") != "" {
				flipMode.Set(m.form.GetString("flipMode"))
			}

			const statusWidth = 50
			statusMarginLeft := m.width - statusWidth - lipgloss.Width(form) - s.Status.GetMarginRight()
			status = s.Status.
				Height(lipgloss.Height(form)).
				Width(statusWidth).
				MarginLeft(statusMarginLeft).
				Render(s.StatusHeader.Render("Current Build") + "\n" +
					buildInfo +
					"\n\n" +
					inputFolderPath +
					"\n\n" +
					outputFolderPath +
					"\n\n" +
					flipMode.String() +
					"\n\n")
		}

		errors := m.form.Errors()
		header := m.appBoundaryView("Image Flipper")
		if len(errors) > 0 {
			header = m.appErrorBoundaryView(m.errorView())
		}
		body := lipgloss.JoinHorizontal(lipgloss.Top, form, status)

		footer := m.appBoundaryView(m.form.Help().ShortHelpView(m.form.KeyBinds()))
		if len(errors) > 0 {
			footer = m.appErrorBoundaryView("")
		}

		return s.Base.Render(header + "\n" + body + "\n\n" + footer)
	}
}

func (m Model) errorView() string {
	var s string
	for _, err := range m.form.Errors() {
		s += err.Error()
	}
	return s
}

func (m Model) appBoundaryView(text string) string {
	return lipgloss.PlaceHorizontal(
		m.width,
		lipgloss.Left,
		m.styles.HeaderText.Render(text),
		lipgloss.WithWhitespaceChars("/"),
		lipgloss.WithWhitespaceForeground(indigo),
	)
}

func (m Model) appErrorBoundaryView(text string) string {
	return lipgloss.PlaceHorizontal(
		m.width,
		lipgloss.Left,
		m.styles.ErrorHeaderText.Render(text),
		lipgloss.WithWhitespaceChars("/"),
		lipgloss.WithWhitespaceForeground(red),
	)
}

func (m Model) getFlipOutcomeView() string {
	inputFolderPath := m.form.GetString("inputFolderPath")
	outputFolderPath := m.form.GetString("outputFolderPath")
	flipMode := m.form.GetString("flipMode")

	return fmt.Sprintf("Input Folder: %s\nOutput Folder: %s\nFlip Mode: %s\n", inputFolderPath, outputFolderPath, flipMode)
}

func (m Model) getBuildInfoView() string {
	if m.loading {
		return fmt.Sprintf("Processing images %s", m.spinner.View())
	} else {
		return ""
	}
}

func checkIfFolderExists(folderPath string) error {
	if _, err := imgproc.CheckIfFolderExists(folderPath); err != nil {
		return err
	}

	return nil
}

func (m Model) flipImages() error {
	ctx := context.Background()
	input := m.form.GetString("inputFolderPath")
	output := m.form.GetString("outputFolderPath")

	var direction imgproc.FlipDirection
	if err := direction.Set(m.form.GetString("flipMode")); err != nil {
		return err
	}

	if _, err := imgproc.RunProcessImagesPipeline(ctx, input, output, direction); err != nil {
		return err
	}

	return nil
}
