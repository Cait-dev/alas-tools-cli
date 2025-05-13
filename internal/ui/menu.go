package ui

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/Cait-dev/alas-tools-cli/internal/handlers"
)

type menuItem struct {
	title, desc string
}

func (m menuItem) Title() string       { return m.title }
func (m menuItem) Description() string { return m.desc }
func (m menuItem) FilterValue() string { return m.title }

type model struct {
	list     list.Model
	quitting bool
	action   int
}

func initialModel() model {
	items := []list.Item{
		menuItem{title: "Opción 1: Corregir X&Y", desc: "Herramienta para corregir coordenadas usando Google Places"},
		menuItem{title: "Opción 2: Mostrar una ruta optimizada de un pallet", desc: "Compara las rutas"},
		menuItem{title: "Opción 3: Obtener coordenadas", desc: "Extrae coordenadas de un pallet y las guarda en un archivo"},
		menuItem{title: "Opción 4: Generar mapa HTML", desc: "Crea un mapa interactivo a partir de un archivo de coordenadas"},
		menuItem{title: "Opción 5: Ayuda", desc: "Muestra la información de ayuda"},
		menuItem{title: "Salir", desc: "Salir de la aplicación"},
	}

	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = selectedItemStyle
	delegate.Styles.NormalTitle = itemStyle

	l := list.New(items, delegate, 80, 20)
	l.Title = "Menú Principal"
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	return model{
		list:     l,
		quitting: false,
		action:   -1,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			i, ok := m.list.SelectedItem().(menuItem)
			if ok {
				if strings.Contains(i.title, "Salir") {
					m.quitting = true
					return m, tea.Quit
				}

				m.action = m.list.Index()
				return m, tea.Quit
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.quitting {
		return quitTextStyle.Render("¡Gracias por usar Alas-Tools-Cli!")
	}

	return docStyle.Render(m.list.View())
}

func ShowStartScreen() {
	fmt.Print("\033[H\033[2J")
	asciiArt := `
/$$$$$$  /$$                         /$$$$$$$$                  /$$              /$$$$$$  /$$ /$$
/$$__  $$| $$                        |__  $$__/                 | $$             /$$__  $$| $$|__/
| $$  \ $$| $$ /$$   /$$ /$$$$$$$       | $$  /$$$$$$   /$$$$$$ | $$  /$$$$$$$ | $$  \__/| $$ /$$
| $$$$$$$$| $$|  $$ /$$/| $$__  $$      | $$ /$$__  $$ /$$__  $$| $$ /$$_____/ | $$      | $$| $$
| $$__  $$| $$ \  $$$$/ | $$  \ $$      | $$| $$  \ $$| $$  \ $$| $$|  $$$$$$ | $$      | $$| $$
| $$  | $$| $$  >$$  $$ | $$  | $$      | $$| $$  | $$| $$  | $$| $$ \____  $$| $$    $$| $$| $$
| $$  | $$| $$ /$$/\  $$| $$  | $$      | $$|  $$$$$$/|  $$$$$$/| $$ /$$$$$$$/|  $$$$$$/| $$| $$
|__/  |__/|__/|__/  \__/|__/  |__/      |__/ \______/  \______/ |__/|_______/  \______/ |__/|__/
`
	subtitulo := `
┌─────────────────────────────────┐
│ Entertainment - Development     │
│ Productivity - & much more!     │
└─────────────────────────────────┘
`
	verde := "\033[32m"
	reset := "\033[0m"
	fmt.Println(verde + asciiArt + subtitulo + reset)

	fmt.Println("\nBienvenido a Alas-Tools-Cli v1.1.1")
	fmt.Println("─────────────────────────────")
	fmt.Println("Use the arrow keys to navigate: ↑ ↓")
	time.Sleep(2 * time.Second)
}

func StartMainMenu() {
	salir := false
	for !salir {
		m := initialModel()
		p := tea.NewProgram(m)
		finalModel, err := p.Run()
		if err != nil {
			fmt.Printf("Error al iniciar la aplicación: %v\n", err)
			os.Exit(1)
		}

		if finalModel, ok := finalModel.(model); ok {
			if finalModel.quitting {
				salir = true
			} else {
				switch finalModel.action {
				case 0:
					handlers.CorregirXY()
				case 1:
					handlers.MostrarRutaOptimizada()
				case 2:
					handlers.ObtenerCoordenadas()
				case 3:
					handlers.GenerarMapaHTML("")
				case 4:
					handlers.MostrarAyuda()
				}

				fmt.Print("\033[H\033[2J")
				fmt.Println("\nVolviendo al menú principal...")
				time.Sleep(1 * time.Second)
			}
		}
	}
}
