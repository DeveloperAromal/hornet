package cli

import (
	"fmt"
	"os"
	"strings"
	"unicode/utf8"

	"golang.org/x/term"
)

func Banner() {
	lines := []string{
		`                  :                                                  `,
		`                 t#,               L.                     ,;         `,
		`  .    .        ;##W.   j.         EW:        ,ft       f#i          `,
		`  Di   Dt      :#L:WE   EW,        E##;       t#E     .E#t  GEEEEEEEL`,
		`  E#i  E#i    .KG  ,#D  E##j       E###t      t#E    i#W,   ,;;L#K;;.`,
		`  E#t  E#t    EE    ;#f E###D.     E#fE#f     t#E   L#D.       t#E   `,
		`  E#t  E#t   f#.     t#iE#jG#W;    E#t D#G    t#E :K#Wfff;     t#E   `,
		`  E########f.:#G     GK E#t t##f   E#t  f#E.  t#E i##WLLLLt    t#E   `,
		`  E#j..K#j... ;#L   LW. E#t  :K#E: E#t   t#K: t#E  .E#L        t#E   `,
		`  E#t  E#t     t#f f#:  E#KDDDD###iE#t    ;#W,t#E    f#E:      t#E   `,
		`  E#t  E#t      f#D#;   E#f,t#Wi,,,E#t     :K#D#E     ,WW;     t#E   `,
		`  f#t  f#t       G#t    E#t  ;#W:  E#t      .E##E      .D#;    t#E   `,
		`   ii   ii        t     DWi   ,KK: ..         G#E        tt     fE   `,
		`                                               fE                :   `,
		`                                                ,                    `,
	}

	width := 80
	if w, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil && w > 0 {
		width = w
	}

	maxLen := 0
	for _, line := range lines {
		if n := utf8.RuneCountInString(line); n > maxLen {
			maxLen = n
		}
	}

	leftPad := (width - maxLen) / 2
	if leftPad < 0 {
		leftPad = 0
	}
	padding := strings.Repeat(" ", leftPad)

	const (
		colorReset  = "\033[0m"
		colorOrange = "\033[38;5;208m"
	)

	fmt.Println()
	for _, line := range lines {
		fmt.Printf("%s%s%s%s\n", padding, colorOrange, line, colorReset)
	}
	fmt.Println()
}
