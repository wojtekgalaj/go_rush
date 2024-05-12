package main

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"syscall"

	tea "github.com/charmbracelet/bubbletea"
)

func isExecutable(path string) bool {
	info, err := os.Stat(path)

	if err != nil {
		return false
	}

	if (info.Mode() & 0111) == 0 {
		return false
	}

	return true
}

type exe struct {
	exePath string
	exeDesc string
}

type model struct {
	execs  []fs.DirEntry
	cursor int
}

func initialModel(list []fs.DirEntry) model {
	model := model{
		execs:  list,
		cursor: 0,
	}

	return model
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			m.cursor--
		case "down", "j":
			m.cursor++
		case "d":
			describe(m.execs[m.cursor].Name())
		case "enter":
			exePath := m.execs[m.cursor].Name()

			executable, lookupErr := exec.LookPath("./" + exePath)

			if lookupErr != nil {
				fmt.Println(lookupErr.Error())
			}

			execErr := syscall.Exec(executable, []string{exePath}, os.Environ())

			if execErr != nil {
				fmt.Println("The program has failed to run. Program's problem, not yours.")
			}

			fmt.Println("The progam has been run")
			os.Exit(0)
		}
	}

	return m, nil
}

func (m model) View() string {
	s := "Which executable would you like to run?\n\n"

	for i, f := range m.execs {
		selMarker := "....."
		if i == m.cursor {
			selMarker = ">...."
		}
		s += fmt.Sprintf("%s%s\n", selMarker, f.Name())
	}

	s += "\n\n Press enter to run the program or d to add a description to it. Press q to quit."

	return s
}

func describe(exeName string) {
	// get description from standard in
	var desc string
	fmt.Scanln(&desc)
	fmt.Println("Description added" + desc)
}

func main() {
	list, err := os.ReadDir(".")

	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	var e []fs.DirEntry

	for _, f := range list {
		if isExecutable(f.Name()) {
			e = append(e, f)
		}
	}

	if len(e) == 0 {
		fmt.Println("No executables found in the current directory")
	}

	p := tea.NewProgram(initialModel(e))

	if err := p.Start(); err != nil {
		fmt.Printf("Oh no! %v\n", err)
		os.Exit(1)
	}

	os.Exit(0)

}
