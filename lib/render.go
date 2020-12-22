package lib

import (
	"fmt"
	"runtime"
)

//for level, pass to print buffer
//if no children,
//	add width to lower level

var (
	Reset = "\033[0m"
	cType = "\033[31m" // red
	cWord = "\033[32m" // green
)

func (sen sentence) Render() {
	if runtime.GOOS == "windows" {
		cType = ""
		cWord = ""
	}

	lines := make([]string, sen.levels+1)
	queue := []*syntaxNode{sen.Syntax} // process nodes in order

	for len(queue) != 0 {
		current := queue[0]
		queue = queue[1:]

		if len(current.Type) == 0 && len(current.Word) == 0 && current.Daughters == nil {
			space := ""
			for i := 0; i < current.Length; i++ { // set padding to left of word
				space += " "
			}
			lines[current.Level] += space
			if current.Level < sen.levels {
				queue = append(queue, &syntaxNode{Length: current.Length, Level: current.Level + 1})

			}

			continue
		}

		var text, color string

		if len(current.Word) != 0 {
			color = cWord
			switch current.startEnd {
			case 0:
				text = current.Word
			case '(':
				text = "(" + current.Word
			case ')':
				text = current.Word + ")"
			default:
				text = "(" + current.Word + ")"
			}

		} else {
			text = current.Type
			color = cType
			if len(text) == 0 {
				text = "*"
			}
		}

		padding := current.Length - len(text)
		lPadding := ""
		rPadding := ""
		for i := 0; i < padding/2; i++ { // set padding to left of word
			lPadding += " "
		}
		for i := 0; i < (padding/2)+(padding%2); i++ { // set padding to right of word
			rPadding += " "
		}

		lines[current.Level] += lPadding + color + text + rPadding

		if current.Daughters == nil {
			if current.Level < sen.levels {
				queue = append(queue, &syntaxNode{Length: current.Length, Level: current.Level + 1})

			}
			continue
		} else {
			for _, child := range current.Daughters {
				queue = append(queue, child)
			}
		}

	}
	for _, line := range lines {
		fmt.Println(line)
	}
	fmt.Println(Reset)
}
