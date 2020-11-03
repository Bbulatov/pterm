package pterm

import (
	"strings"

	"github.com/mattn/go-runewidth"

	"github.com/pterm/pterm/internal"
)

// Panel contains the data, which should be printed inside a PanelPrinter.
type Panel struct {
	Data string
}

// Panels is a two dimensional coordinate system for Panel.
type Panels [][]Panel

// DefaultPanel is the default PanelPrinter.
var DefaultPanel = PanelPrinter{
	Padding: 1,
}

// PanelPrinter prints content in boxes.
type PanelPrinter struct {
	Panels          Panels
	Padding         int
	BottomPadding   int
	SameColumnWidth bool
}

// WithPanels returns a new PanelPrinter with specific options.
func (p PanelPrinter) WithPanels(panels Panels) *PanelPrinter {
	p.Panels = panels
	return &p
}

// WithPadding returns a new PanelPrinter with specific options.
func (p PanelPrinter) WithPadding(padding int) *PanelPrinter {
	p.Padding = padding
	return &p
}

// WithBottomPadding returns a new PanelPrinter with specific options.
func (p PanelPrinter) WithBottomPadding(bottomPadding int) *PanelPrinter {
	if bottomPadding < 0 {
		bottomPadding = 0
	}
	p.BottomPadding = bottomPadding
	return &p
}

// WithSameColumnWidth returns a new PanelPrinter with specific options.
func (p PanelPrinter) WithSameColumnWidth(b ...bool) *PanelPrinter {
	b2 := internal.WithBoolean(b)
	p.SameColumnWidth = b2
	return &p
}

// Srender renders the Template as a string.
func (p PanelPrinter) Srender() (string, error) {
	var ret string

	columnMaxHeightMap := make(map[int]int)

	if p.SameColumnWidth {
		for _, panel := range p.Panels {
			for i, p2 := range panel {
				if columnMaxHeightMap[i] < internal.GetStringMaxWidth(p2.Data) {
					columnMaxHeightMap[i] = internal.GetStringMaxWidth(p2.Data)
				}
			}
		}
	}

	for j, boxLine := range p.Panels {
		var maxHeight int

		for _, box := range boxLine {
			height := len(strings.Split(box.Data, "\n"))
			if height > maxHeight {
				maxHeight = height
			}
		}

		var renderedPanels []string

		for _, box := range boxLine {
			renderedPanels = append(renderedPanels, box.Data)
		}

		for i := 0; i <= maxHeight; i++ {
			if maxHeight != i {
				for j, letter := range renderedPanels {
					var letterLine string
					letterLines := strings.Split(letter, "\n")
					var maxLetterWidth int
					if !p.SameColumnWidth {
						maxLetterWidth = internal.GetStringMaxWidth(letter)
					}
					if len(letterLines) > i {
						letterLine = letterLines[i]
					}
					letterLineLength := runewidth.StringWidth(letterLine)
					if !p.SameColumnWidth {
						if letterLineLength < maxLetterWidth {
							letterLine += strings.Repeat(" ", maxLetterWidth-letterLineLength)
						}
					} else {
						if letterLineLength < columnMaxHeightMap[j] {
							letterLine += strings.Repeat(" ", columnMaxHeightMap[j]-letterLineLength)
						}
					}
					letterLine += strings.Repeat(" ", p.Padding)
					ret += letterLine
				}
				ret += "\n"
			} else if j+1 != len(p.Panels) {
				ret += strings.Repeat("\n", p.BottomPadding)
			}
		}
	}

	return ret, nil
}

// Render prints the Template to the terminal.
func (p PanelPrinter) Render() error {
	s, _ := p.Srender()
	Println(s)

	return nil
}
