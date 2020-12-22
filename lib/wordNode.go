package lib

import (
	"errors"
)

type syntaxNode struct {
	Type      string
	Word      string
	Length    int
	Level     int
	startEnd  rune // used to determine relative position of word during rendering
	Daughters []*syntaxNode
}

type sentence struct {
	S      string
	Syntax *syntaxNode
	levels int
}

func NewSentence(s string) (*sentence, error) {
	bCnt := 0
	for _, char := range s {
		switch char {
		case '[':
			bCnt++
		case ']':
			bCnt--
		}
	}
	if bCnt != 0 || s[0] != '[' {
		return nil, errors.New("missing bracket(s)")
	}

	ns := &sentence{S: s, levels: 0}
	ns.Syntax = ns.buildTree([]rune(s[1:len(s)-1]), 0)

	return ns, nil
}

// buildTree constructs the tree level by level
func (sen *sentence) buildTree(r []rune, level int) *syntaxNode {
	if level > sen.levels {
		sen.levels = level
	}

	const MARGIN = 4 // how much space between words when displayed
	var (
		bCnt, length, wLength, siblingCount int
		synType                             string
		buffer                              []rune
		nxt                                 *syntaxNode
	)
	branches := make([]*syntaxNode, 0)
	i := 0

	for ; i < len(r) && r[i] != ' ' && r[i] != '['; i++ {
		buffer = append(buffer, r[i])
		wLength++
	}
	if wLength != 0 {
		synType = string(buffer)
		buffer = buffer[:0]
		wLength = 0
	}
	for ; i < len(r); i++ {
		switch r[i] {
		case '[':
			i++
			bCnt++
			if len(buffer) != 0 {
				branches = append(branches, &syntaxNode{Word: string(buffer),
					Length: wLength + MARGIN,
					Level:  level + 1})
				length += wLength + MARGIN
				buffer = buffer[:0]
			}
		b:
			for ; i < len(r) && bCnt != 0; i++ {
				switch r[i] {
				case '[':
					bCnt++
				case ']':
					bCnt--
				}

				if bCnt == 0 {
					buffer = append(buffer, ' ') // add space at end to ensure loop does not end prematurely
					break b
				}
				buffer = append(buffer, r[i]) // buffer characters to add to recursion
				wLength += 1
			}
			if wLength != 0 {
				nxt = sen.buildTree(buffer, level+1)
				branches = append(branches, nxt)
				length += nxt.Length
				wLength = 0
			}
		case ' ':
			break
		default:
			buffer = append(buffer, r[i])
			wLength++
			continue
		}
		if wLength != 0 {
			w := string(buffer)
			branches = append(branches, &syntaxNode{Word: w,
				Length: wLength + MARGIN,
				Level:  level + 1})
			length += wLength + MARGIN
			wLength = 0
			siblingCount++

			if level+1 > sen.levels {
				sen.levels = level + 1
			}
		}

		buffer = buffer[:0]
	}
	if branches != nil {
		first, last := -1, 0
		for i := range branches {
			if len(branches[i].Word) != 0 {
				last = i
				if first == -1 {
					first = i
				}
			}
		}
		if first != -1 {
			if first == last {
				branches[0].startEnd = '|'
			} else {
				branches[first].startEnd = '('
				branches[last].startEnd = ')'
			}
		}
	}

	node := &syntaxNode{Type: synType,
		Length:    length,
		Daughters: branches,
		Level:     level}
	return node
}
