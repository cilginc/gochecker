package ui

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

var (
	Err     = color.New(color.FgRed, color.Bold).SprintFunc()
	Title   = color.New(color.Bold).SprintFunc()
	Success = color.New(color.FgGreen).SprintFunc()
	Warn    = color.New(color.FgYellow).SprintFunc()
	Info    = color.New(color.FgCyan).SprintFunc()
)

func DisableColor() {
	color.NoColor = true
}

func PrintErr(a ...any) {
	fmt.Fprintln(os.Stderr, Err(a...))
}

func PrintTitle(a ...any) {
	fmt.Println(Title(a...))
}

func PrintSuccess(a ...any) {
	fmt.Println(Success(a...))
}

func PrintWarn(a ...any) {
	fmt.Println(Warn(a...))
}

func PrintInfo(a ...any) {
	fmt.Println(Info(a...))
}
