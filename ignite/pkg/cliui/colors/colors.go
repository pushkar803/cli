package colors

import "github.com/fatih/color"

var (
	Info  = color.New(color.FgYellow).SprintFunc()
	Infof = color.New(color.FgYellow).SprintfFunc()
)
