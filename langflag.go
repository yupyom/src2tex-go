package main

import (
	"path/filepath"
	"strings"
)

// Helpers
func SkipSpaces1(ptr *int) {
	for {
		c := byte(Buffer[*ptr])
		if c > '\x00' && c <= ' ' {
			*ptr = IncBufPtr(*ptr)
		} else {
			break
		}
	}
}
func SkipSpaces2(ptr *int) {
	for {
		c := byte(Buffer[*ptr])
		if c > '\x00' && c <= ' ' && c != '\n' {
			*ptr = IncBufPtr(*ptr)
		} else {
			break
		}
	}
}
func SkipSpaces1Str(s string, idx *int) {
	for *idx < len(s) {
		c := byte(s[*idx])
		if c > '\x00' && c <= ' ' {
			*idx++
		} else {
			break
		}
	}
}
func SkipSpaces2Str(s string, idx *int) {
	for *idx < len(s) {
		c := byte(s[*idx])
		if c > '\x00' && c <= ' ' && c != '\n' {
			*idx++
		} else {
			break
		}
	}
}
func SearchPhrase(bufPtr int, phrase string) int {
	bPtr1 := bufPtr
	flag := 1
	for bPtr1 != BufferSize/2 {
		bPtr2 := bPtr1
		SkipSpaces1(&bPtr2)
		cIdx := 0
		SkipSpaces1Str(phrase, &cIdx)
		for cIdx < len(phrase) {
			c := byte(Buffer[bPtr2])
			if c >= 'A' && c <= 'Z' {
				c += 0x20
			}
			flag = int(c) - int(phrase[cIdx])
			if flag != 0 {
				break
			}
			bPtr2 = IncBufPtr(bPtr2)
			SkipSpaces1(&bPtr2)
			cIdx++
			SkipSpaces1Str(phrase, &cIdx)
		}
		if flag == 0 {
			break
		}
		bPtr1++
	}
	if flag == 0 {
		return 1
	}
	return 0
}
func SearchTwoPhrases(bufPtr int, phrase1, phrase2 string) int {
	bPtr1 := bufPtr
	flag := 1
	for bPtr1 != BufferSize/2 {
		bPtr2 := bPtr1
		SkipSpaces1(&bPtr2)
		cIdx := 0
		SkipSpaces1Str(phrase1, &cIdx)
		for cIdx < len(phrase1) {
			c := byte(Buffer[bPtr2])
			if c >= 'A' && c <= 'Z' {
				c += 0x20
			}
			flag = int(c) - int(phrase1[cIdx])
			if flag != 0 {
				break
			}
			bPtr2 = IncBufPtr(bPtr2)
			SkipSpaces1(&bPtr2)
			cIdx++
			SkipSpaces1Str(phrase1, &cIdx)
		}
		if flag == 0 {
			break
		}
		bPtr1++
	}
	if flag != 0 {
		return 0
	}
	for i := 0; i < 40; i++ {
		bPtr2 := bPtr1
		SkipSpaces1(&bPtr2)
		cIdx := 0
		SkipSpaces1Str(phrase2, &cIdx)
		for cIdx < len(phrase2) {
			c := byte(Buffer[bPtr2])
			if c >= 'A' && c <= 'Z' {
				c += 0x20
			}
			flag = int(c) - int(phrase2[cIdx])
			if flag != 0 {
				break
			}
			bPtr2 = IncBufPtr(bPtr2)
			SkipSpaces1(&bPtr2)
			cIdx++
			SkipSpaces1Str(phrase2, &cIdx)
		}
		if flag == 0 {
			break
		}
		bPtr1++
	}
	if flag == 0 {
		return 1
	}
	return 0
}
func GetPhrase(bufPtr int, phrase string) int {
	bPtr1 := bufPtr
	flag := 1
	for i := 0; i < 64; i++ {
		bPtr2 := bPtr1
		cIdx := 0
		for cIdx < len(phrase) {
			c := byte(Buffer[bPtr2])
			if c >= 'A' && c <= 'Z' {
				c += 0x20
			}
			flag = int(c) - int(phrase[cIdx])
			if flag != 0 {
				break
			}
			bPtr2 = IncBufPtr(bPtr2)
			cIdx++
		}
		if flag == 0 {
			break
		}
		bPtr1++
	}
	if flag == 0 {
		return bPtr1
	}
	return -1
}
func SearchLine(bufPtr int, phrase string) int {
	bPtr := bufPtr
	SkipSpaces2(&bPtr)
	cIdx := 0
	SkipSpaces2Str(phrase, &cIdx)
	var line1 []byte
	var line2 []byte
	for i := 0; i < 255 && cIdx < len(phrase); i++ {
		c := byte(Buffer[bPtr])
		if c == '\n' {
			c = 0x00
		}
		if c >= 'A' && c <= 'Z' {
			c += 0x20
		}
		line1 = append(line1, c)
		line2 = append(line2, phrase[cIdx])
		bPtr = IncBufPtr(bPtr)
		SkipSpaces2(&bPtr)
		cIdx++
		SkipSpaces2Str(phrase, &cIdx)
	}
	flag := 1
	if len(line1) == len(line2) {
		flag = 0
		for i := 0; i < len(line2); i++ {
			if line1[i] != line2[i] {
				flag = int(line1[i]) - int(line2[i])
				break
			}
		}
	}
	if flag == 0 {
		return 1
	}
	return 0
}

func SetBasFlag(bufPtr int) {
	if MFlag != 0 {
		return
	}
	BASFlag = 10 * SearchPhrase(bufPtr, "' basic")
	BASFlag += 10 * SearchPhrase(bufPtr, "rem basic")
	BASFlag += SearchPhrase(bufPtr, "0 '")
	BASFlag += SearchPhrase(bufPtr, "1 '")
	BASFlag += SearchPhrase(bufPtr, "2 '")
	BASFlag += SearchPhrase(bufPtr, "3 '")
	BASFlag += SearchPhrase(bufPtr, "4 '")
	BASFlag += SearchPhrase(bufPtr, "5 '")
	BASFlag += SearchPhrase(bufPtr, "6 '")
	BASFlag += SearchPhrase(bufPtr, "7 '")
	BASFlag += SearchPhrase(bufPtr, "8 '")
	BASFlag += SearchPhrase(bufPtr, "9 '")
	BASFlag += 2 * SearchPhrase(bufPtr, "0 rem")
	BASFlag += 2 * SearchPhrase(bufPtr, "1 rem")
	BASFlag += 2 * SearchPhrase(bufPtr, "2 rem")
	BASFlag += 2 * SearchPhrase(bufPtr, "3 rem")
	BASFlag += 2 * SearchPhrase(bufPtr, "4 rem")
	BASFlag += 2 * SearchPhrase(bufPtr, "5 rem")
	BASFlag += 2 * SearchPhrase(bufPtr, "6 rem")
	BASFlag += 2 * SearchPhrase(bufPtr, "7 rem")
	BASFlag += 2 * SearchPhrase(bufPtr, "8 rem")
	BASFlag += 2 * SearchPhrase(bufPtr, "9 rem")
	BASFlag += SearchPhrase(bufPtr, "cls")
	BASFlag += SearchPhrase(bufPtr, "defdbl")
	BASFlag += SearchPhrase(bufPtr, "<>")
	BASFlag += 2 * SearchPhrase(bufPtr, "gosub")
	BASFlag += 3 * SearchPhrase(bufPtr, "then gosub")
	BASFlag += 3 * SearchPhrase(bufPtr, "then print")
	BASFlag += SearchPhrase(bufPtr, "wend")
	BASFlag += SearchPhrase(bufPtr, "input \"")
	BASFlag += SearchPhrase(bufPtr, "data \"")
	BASFlag += SearchPhrase(bufPtr, "print \"")
	BASFlag += 2 * SearchPhrase(bufPtr, "$=\"")
	BASFlag += 4 * SearchPhrase(bufPtr, "left$(")
	BASFlag += 4 * SearchPhrase(bufPtr, "mid$(")
	BASFlag += 4 * SearchPhrase(bufPtr, "right$(")
	BASFlag += SearchPhrase(bufPtr, "dir$(")
	BASFlag += SearchPhrase(bufPtr, "getattr(")
	BASFlag += SearchPhrase(bufPtr, "chdir \"")
	BASFlag += SearchPhrase(bufPtr, "curdir$")
	BASFlag += SearchPhrase(bufPtr, "select case")
	BASFlag += SearchPhrase(bufPtr, "end select")
	BASFlag += SearchPhrase(bufPtr, "exit proc")
	BASFlag += SearchPhrase(bufPtr, "end proc")
	BASFlag += SearchPhrase(bufPtr, "end sub")
}

func SetCFlag(bufPtr int) {
	CFlag = 10 * SearchPhrase(bufPtr, "")
	CFlag += 10 * SearchPhrase(bufPtr, "// c")
	CFlag += 10 * SearchPhrase(bufPtr, "")
	CFlag += 10 * SearchPhrase(bufPtr, "// objective-c")
	CFlag += 5 * SearchPhrase(bufPtr, "// //")
	CFlag += 5 * SearchPhrase(bufPtr, "#include <")
	CFlag += 5 * SearchPhrase(bufPtr, "#include \"")
	CFlag += 5 * SearchPhrase(bufPtr, "#import <")
	CFlag += 5 * SearchPhrase(bufPtr, "#import \"")
	CFlag += SearchPhrase(bufPtr, "#define")
	CFlag += 4 * (SearchPhrase(bufPtr, "#if") * SearchPhrase(bufPtr, "#endif"))
	CFlag += 5 * SearchPhrase(bufPtr, "@implementation")
	CFlag += 5 * SearchPhrase(bufPtr, "@interface")
	CFlag += 5 * SearchPhrase(bufPtr, "@private")
	CFlag += 5 * SearchPhrase(bufPtr, "@protected")
	CFlag += 5 * SearchPhrase(bufPtr, "@public")
	CFlag += 5 * SearchPhrase(bufPtr, "@selector(")
	CFlag += 5 * SearchPhrase(bufPtr, "@def(")
	CFlag += 5 * SearchPhrase(bufPtr, "@encode(")
	CFlag += 3 * SearchPhrase(bufPtr, "main() {")
	CFlag += 5 * SearchPhrase(bufPtr, "main(argc,argv) int argc;")
	CFlag += SearchPhrase(bufPtr, "printf(\"")
	CFlag += SearchPhrase(bufPtr, "; } }")
	CFlag += SearchPhrase(bufPtr, "++")
	CFlag += SearchPhrase(bufPtr, "--")
	CFlag += 2 * SearchPhrase(bufPtr, "+=")
	CFlag += 2 * SearchPhrase(bufPtr, "-=")
	CFlag += 2 * SearchPhrase(bufPtr, "*=")
	CFlag += 2 * SearchPhrase(bufPtr, "/=")
	CFlag += SearchPhrase(bufPtr, ")**(")
	CFlag += SearchPhrase(bufPtr, ")++(")
}

func SetCblFlag(bufPtr int) {
	if MFlag != 0 {
		return
	}
	CBLFlag = 10 * SearchPhrase(bufPtr, "* cobol")
	CBLFlag += 10 * SearchPhrase(bufPtr, "/ cobol")
	CBLFlag += 4 * SearchPhrase(bufPtr, "identification division.")
	CBLFlag += 4 * SearchPhrase(bufPtr, "program-id.")
	CBLFlag += SearchPhrase(bufPtr, "author.")
	CBLFlag += 4 * SearchPhrase(bufPtr, "date-written.")
	CBLFlag += 4 * SearchPhrase(bufPtr, "date-written.")
	CBLFlag += 4 * SearchPhrase(bufPtr, "environment division.")
	CBLFlag += 4 * SearchPhrase(bufPtr, "configuration section.")
	CBLFlag += 4 * SearchPhrase(bufPtr, "source-computer.")
	CBLFlag += 4 * SearchPhrase(bufPtr, "object-computer.")
	CBLFlag += 4 * SearchPhrase(bufPtr, "special-names.")
	CBLFlag += 3 * SearchPhrase(bufPtr, "input-output section.")
	CBLFlag += 4 * SearchPhrase(bufPtr, "file-contral.")
	CBLFlag += 2 * SearchPhrase(bufPtr, "i-o-contral.")
	CBLFlag += 4 * SearchPhrase(bufPtr, "data division.")
	CBLFlag += 3 * SearchPhrase(bufPtr, "file section.")
	CBLFlag += 4 * SearchPhrase(bufPtr, "working-storage section.")
	CBLFlag += 4 * SearchPhrase(bufPtr, "procedure division.")
	CBLFlag += SearchPhrase(bufPtr, "section.")
	CBLFlag += 4 * SearchPhrase(bufPtr, "stop run.")
}

func SetF77Flag(bufPtr int) {
	if MFlag != 0 {
		return
	}
	F77Flag = 10 * SearchPhrase(bufPtr, "c fortran")
	F77Flag += 10 * SearchPhrase(bufPtr, "* fortran")
	F77Flag += SearchPhrase(bufPtr, "implicit")
	F77Flag += 2 * SearchPhrase(bufPtr, "logical*2")
	F77Flag += 2 * SearchPhrase(bufPtr, "logical*4")
	F77Flag += 2 * SearchPhrase(bufPtr, "integer*2")
	F77Flag += 2 * SearchPhrase(bufPtr, "integer*4")
	F77Flag += 2 * SearchPhrase(bufPtr, "real*4")
	F77Flag += 2 * SearchPhrase(bufPtr, "real*8")
	F77Flag += (SearchPhrase(bufPtr, "do") * SearchPhrase(bufPtr, "continue"))
	F77Flag += (SearchPhrase(bufPtr, "subroutine") * SearchPhrase(bufPtr, "end"))
	F77Flag += 3 * (SearchPhrase(bufPtr, "write(") * SearchPhrase(bufPtr, "format("))
	F77Flag += SearchPhrase(bufPtr, "stop end")
	F77Flag += 4 * SearchPhrase(bufPtr, ".lt.")
	F77Flag += 4 * SearchPhrase(bufPtr, ".le.")
	F77Flag += 4 * SearchPhrase(bufPtr, ".eq.")
	F77Flag += 4 * SearchPhrase(bufPtr, ".ne.")
	F77Flag += 4 * SearchPhrase(bufPtr, ".gt.")
	F77Flag += 4 * SearchPhrase(bufPtr, ".ge.")
	F77Flag += 4 * SearchPhrase(bufPtr, ".or.")
	F77Flag += 4 * SearchPhrase(bufPtr, ".and.")
	F77Flag += 4 * SearchPhrase(bufPtr, ".not.")
}

func SetHtmlFlag(bufPtr int) {
	HTMLFlag = 10 * SearchPhrase(bufPtr, "<!- html ->/")
	HTMLFlag += 10 * SearchPhrase(bufPtr, "</a>")
	HTMLFlag += 5 * SearchPhrase(bufPtr, "</address>")
	HTMLFlag += 5 * SearchPhrase(bufPtr, "</b>")
	HTMLFlag += 10 * SearchPhrase(bufPtr, "</base href=")
	HTMLFlag += 5 * SearchPhrase(bufPtr, "</body>")
	HTMLFlag += 5 * SearchPhrase(bufPtr, "</caption>")
	HTMLFlag += 5 * SearchPhrase(bufPtr, "</cite>")
	HTMLFlag += 5 * SearchPhrase(bufPtr, "</dir>")
	HTMLFlag += 5 * SearchPhrase(bufPtr, "</dl>")
	HTMLFlag += 5 * SearchPhrase(bufPtr, "</em>")
	HTMLFlag += 5 * SearchPhrase(bufPtr, "</h1>")
	HTMLFlag += 5 * SearchPhrase(bufPtr, "</h2>")
	HTMLFlag += 5 * SearchPhrase(bufPtr, "</h3>")
	HTMLFlag += 5 * SearchPhrase(bufPtr, "</h4>")
	HTMLFlag += 5 * SearchPhrase(bufPtr, "</h5>")
	HTMLFlag += 5 * SearchPhrase(bufPtr, "</h6>")
	HTMLFlag += 5 * SearchPhrase(bufPtr, "</head>")
	HTMLFlag += 10 * SearchPhrase(bufPtr, "</html>")
	HTMLFlag += 5 * SearchPhrase(bufPtr, "</i>")
	HTMLFlag += 10 * SearchPhrase(bufPtr, "</img src=")
	HTMLFlag += 5 * SearchPhrase(bufPtr, "</kbd>")
	HTMLFlag += 5 * SearchPhrase(bufPtr, "</menu>")
	HTMLFlag += 5 * SearchPhrase(bufPtr, "</ol>")
	HTMLFlag += 5 * SearchPhrase(bufPtr, "</strong>")
	HTMLFlag += 5 * SearchPhrase(bufPtr, "</table>")
	HTMLFlag += 5 * SearchPhrase(bufPtr, "</title>")
	HTMLFlag += 5 * SearchPhrase(bufPtr, "</tr>")
	HTMLFlag += 5 * SearchPhrase(bufPtr, "</tt>")
	HTMLFlag += 5 * SearchPhrase(bufPtr, "</ul>")
}

func SetJavaFlag(bufPtr int) {
	JAVAFlag = CFlag
	JAVAFlag += 10 * SearchPhrase(bufPtr, "")
	JAVAFlag += 10 * SearchPhrase(bufPtr, "// java")
	JAVAFlag += 10 * SearchPhrase(bufPtr, "synchronized")
	JAVAFlag += 10 * SearchPhrase(bufPtr, "instanceof")
	JAVAFlag += 10 * SearchPhrase(bufPtr, "import java.")
	JAVAFlag += 5 * SearchPhrase(bufPtr, "@author")
	JAVAFlag += 5 * SearchPhrase(bufPtr, "@version")
}

func SetLispFlag(bufPtr int) {
	if MFlag != 0 {
		return
	}
	LISPFlag = 10 * SearchPhrase(bufPtr, "; lisp")
	LISPFlag += 10 * SearchPhrase(bufPtr, "; scheme")
	LISPFlag += 5 * SearchPhrase(bufPtr, ")))))")
	LISPFlag += 4 * SearchPhrase(bufPtr, "(car")
	LISPFlag += 4 * SearchPhrase(bufPtr, "(cdr")
	LISPFlag += 2 * SearchPhrase(bufPtr, "(cons")
	LISPFlag += SearchPhrase(bufPtr, "(list")
	LISPFlag += 4 * SearchPhrase(bufPtr, "(setq")
	LISPFlag += 4 * SearchPhrase(bufPtr, "(set!")
	LISPFlag += 4 * SearchPhrase(bufPtr, "(lambda")
	LISPFlag += SearchPhrase(bufPtr, "(def")
	LISPFlag += 4 * SearchPhrase(bufPtr, "(defun")
	LISPFlag += SearchPhrase(bufPtr, "(define")
	LISPFlag += SearchPhrase(bufPtr, "(let")
	LISPFlag += SearchPhrase(bufPtr, "(+")
	LISPFlag += 3 * SearchPhrase(bufPtr, "(/")
	LISPFlag += 3 * SearchPhrase(bufPtr, "(cond (")
	LISPFlag += SearchPhrase(bufPtr, "(if (")
}

func SetMakeFlag(bufPtr int) {
	if MFlag != 0 {
		return
	}
	for bPtr := bufPtr; bPtr != BufferSize/2; bPtr++ {
		if Buffer[bPtr] == '\n' && Buffer[(bPtr+1)%BufferSize] == '#' && (Buffer[(bPtr+2)%BufferSize] == '\t' || Buffer[(bPtr+2)%BufferSize] == ' ') {
			MAKEFlag++
		}
	}
	MAKEFlag += 10 * SearchPhrase(bufPtr, "# makefile")
	MAKEFlag += 5 * SearchPhrase(bufPtr, "all:")
	MAKEFlag += 5 * SearchPhrase(bufPtr, "install:")
	MAKEFlag += 5 * SearchPhrase(bufPtr, "install.man:")
	MAKEFlag += 5 * SearchPhrase(bufPtr, "depend:")
	MAKEFlag += 5 * SearchPhrase(bufPtr, "clean:")
	MAKEFlag += 3 * SearchPhrase(bufPtr, "lint:")
	MAKEFlag += 3 * SearchPhrase(bufPtr, "diff:")
	MAKEFlag += SearchPhrase(bufPtr, "cc=")
	MAKEFlag += SearchPhrase(bufPtr, "$(cc)")
	MAKEFlag += SearchPhrase(bufPtr, "objs=")
	MAKEFlag += SearchPhrase(bufPtr, "obj=")
	MAKEFlag += SearchPhrase(bufPtr, "srcs=")
	MAKEFlag += SearchPhrase(bufPtr, "src=")
	MAKEFlag += SearchPhrase(bufPtr, "missing=")
	MAKEFlag += SearchPhrase(bufPtr, "optimize=")
	MAKEFlag += SearchPhrase(bufPtr, "parser=")
	MAKEFlag += SearchPhrase(bufPtr, "flags=")
	MAKEFlag += 2 * SearchPhrase(bufPtr, "ldflags=")
	MAKEFlag += 2 * SearchPhrase(bufPtr, "libdir=")
	MAKEFlag += 2 * SearchPhrase(bufPtr, "bindir=")
	MAKEFlag += 2 * SearchPhrase(bufPtr, "destdir=")
	MAKEFlag += 2 * SearchPhrase(bufPtr, "mandir=")
	MAKEFlag += 2 * SearchPhrase(bufPtr, "docdir=")
	MAKEFlag += 5 * SearchPhrase(bufPtr, "@echo")
	MAKEFlag += 5 * SearchPhrase(bufPtr, "@-echo")
}

func SetPasFlag(bufPtr int) {
	if MFlag != 0 {
		return
	}
	for bPtr := bufPtr; bPtr != BufferSize/2; bPtr++ {
		if Buffer[bPtr] == '\n' && Buffer[(bPtr+1)%BufferSize] == '(' && Buffer[(bPtr+2)%BufferSize] == '*' && (Buffer[(bPtr+3)%BufferSize] == ' ' || Buffer[(bPtr+3)%BufferSize] == '*') {
			PASFlag++
		}
	}
	PASFlag += 10 * SearchPhrase(bufPtr, "{ pascal }")
	PASFlag += 10 * SearchPhrase(bufPtr, "(* pascal *)")
	PASFlag += 2 * SearchPhrase(bufPtr, "*)(*")
	PASFlag += 3 * (SearchPhrase(bufPtr, "program") + SearchPhrase(bufPtr, "procedure"))
	PASFlag += (SearchPhrase(bufPtr, "begin") * SearchPhrase(bufPtr, "end"))
	PASFlag += (SearchPhrase(bufPtr, "const") * SearchPhrase(bufPtr, "var"))
	PASFlag += 3 * SearchPhrase(bufPtr, "keypressed(")
	PASFlag += SearchPhrase(bufPtr, "blockread(")
	PASFlag += SearchPhrase(bufPtr, "blockwrite(")
	PASFlag += 4 * SearchPhrase(bufPtr, "readln(")
	PASFlag += 4 * SearchPhrase(bufPtr, "writeln(")
	PASFlag += SearchPhrase(bufPtr, "write('")
	PASFlag += SearchPhrase(bufPtr, ":=")
}

func SetPerlFlag(bufPtr int) {
	if MFlag != 0 {
		return
	}
	for bPtr := bufPtr; bPtr != BufferSize/2; bPtr++ {
		if Buffer[bPtr] == '\n' && Buffer[(bPtr+1)%BufferSize] == '#' && (Buffer[(bPtr+2)%BufferSize] == '\t' || Buffer[(bPtr+2)%BufferSize] == ' ') {
			PERLFlag++
		}
	}
	PERLFlag += 10 * SearchPhrase(bufPtr, "# perl")
	PERLFlag += 10 * SearchTwoPhrases(bufPtr, "#!", "/bin/perl")
	PERLFlag += 10 * SearchTwoPhrases(bufPtr, "#!", "/usr/bin/perl")
	PERLFlag += 10 * SearchTwoPhrases(bufPtr, "#!", "/usr/local/bin/perl")
	PERLFlag += 10 * SearchTwoPhrases(bufPtr, "#!", "/bin/jperl")
	PERLFlag += 10 * SearchTwoPhrases(bufPtr, "#!", "/usr/bin/jperl")
	PERLFlag += 10 * SearchTwoPhrases(bufPtr, "#!", "/usr/local/bin/jperl")
	PERLFlag += 2 * SearchPhrase(bufPtr, ".+")
	PERLFlag += 2 * SearchPhrase(bufPtr, ".-")
	PERLFlag += 3 * SearchPhrase(bufPtr, ".*")
	PERLFlag += 2 * SearchPhrase(bufPtr, "./")
	PERLFlag += 3 * SearchPhrase(bufPtr, ".\\")
	PERLFlag += 3 * SearchPhrase(bufPtr, ".**")
	PERLFlag += 3 * SearchPhrase(bufPtr, ".^")
	PERLFlag += 3 * SearchPhrase(bufPtr, "~=")
	PERLFlag += 3 * SearchPhrase(bufPtr, "=~")
	PERLFlag += 3 * SearchPhrase(bufPtr, ".=")
	PERLFlag += 4 * SearchPhrase(bufPtr, ".'")
	PERLFlag += 4 * SearchPhrase(bufPtr, "/^")
	PERLFlag += 3 * SearchPhrase(bufPtr, "/\\")
	PERLFlag += SearchPhrase(bufPtr, "print \"")
	PERLFlag += SearchPhrase(bufPtr, "<>")
	PERLFlag += 4 * SearchPhrase(bufPtr, "while(<")
	PERLFlag += 4 * SearchPhrase(bufPtr, "$_")
	PERLFlag += 3 * SearchPhrase(bufPtr, "$.")
	PERLFlag += 4 * SearchPhrase(bufPtr, "@_")
	PERLFlag += 2 * SearchPhrase(bufPtr, "__end__")
}

func SetShellFlag(bufPtr int) {
	if MFlag != 0 {
		return
	}
	for bPtr := bufPtr; bPtr != BufferSize/2; bPtr++ {
		if Buffer[bPtr] == '\n' && Buffer[(bPtr+1)%BufferSize] == '#' && (Buffer[(bPtr+2)%BufferSize] == '\t' || Buffer[(bPtr+2)%BufferSize] == ' ') {
			SHFlag++
		}
	}
	SHFlag += 10 * SearchPhrase(bufPtr, "# shell")
	SHFlag += 10 * SearchTwoPhrases(bufPtr, "#!", "/bin/sh")
	SHFlag += 10 * SearchTwoPhrases(bufPtr, "#!", "/usr/bin/sh")
	SHFlag += 10 * SearchTwoPhrases(bufPtr, "#!", "/bin/csh")
	SHFlag += 10 * SearchTwoPhrases(bufPtr, "#!", "/usr/bin/csh")
	SHFlag += 10 * SearchTwoPhrases(bufPtr, "#!", "/bin/ksh")
	SHFlag += 10 * SearchTwoPhrases(bufPtr, "#!", "/usr/bin/ksh")
	SHFlag += 10 * SearchTwoPhrases(bufPtr, "#!", "/bin/tcsh")
	SHFlag += 10 * SearchTwoPhrases(bufPtr, "#!", "/usr/bin/tcsh")
	SHFlag += 10 * SearchTwoPhrases(bufPtr, "#!", "/usr/local/bin/tcsh")
	SHFlag += 10 * SearchTwoPhrases(bufPtr, "#!", "/bin/zsh")
	SHFlag += 10 * SearchTwoPhrases(bufPtr, "#!", "/usr/bin/zsh")
	SHFlag += 10 * SearchTwoPhrases(bufPtr, "#!", "/usr/local/bin/zsh")
	SHFlag += 10 * SearchTwoPhrases(bufPtr, "#!", "/bin/bash")
	SHFlag += 10 * SearchTwoPhrases(bufPtr, "#!", "/usr/bin/bash")
	SHFlag += 10 * SearchTwoPhrases(bufPtr, "#!", "/usr/local/bin/bash")
}

func SetTclFlag(bufPtr int) {
	if MFlag != 0 {
		return
	}
	for bPtr := bufPtr; bPtr != BufferSize/2; bPtr++ {
		if Buffer[bPtr] == '\n' && Buffer[(bPtr+1)%BufferSize] == '#' && (Buffer[(bPtr+2)%BufferSize] == '\t' || Buffer[(bPtr+2)%BufferSize] == ' ') {
			TCLFlag++
		}
	}
	TCLFlag += 10 * SearchPhrase(bufPtr, "# tcl/tk")
	TCLFlag += 10 * SearchTwoPhrases(bufPtr, "#!", "/bin/tclsh")
	TCLFlag += 10 * SearchTwoPhrases(bufPtr, "#!", "/usr/bin/tclsh")
	TCLFlag += 10 * SearchTwoPhrases(bufPtr, "#!", "/usr/local/bin/tclsh")
	TCLFlag += 10 * SearchTwoPhrases(bufPtr, "#!", "/bin/wish")
	TCLFlag += 10 * SearchTwoPhrases(bufPtr, "#!", "/usr/bin/wish")
	TCLFlag += 10 * SearchTwoPhrases(bufPtr, "#!", "/usr/local/bin/wish")
	TCLFlag += 2 * SearchPhrase(bufPtr, "tcl_")
	TCLFlag += 2 * SearchPhrase(bufPtr, "tk_")
	TCLFlag += 2 * SearchPhrase(bufPtr, "button .")
	TCLFlag += 2 * SearchPhrase(bufPtr, "pack .")
	TCLFlag += 2 * SearchPhrase(bufPtr, "destroy .")
	TCLFlag += 2 * SearchPhrase(bufPtr, "entry .")
	TCLFlag += 2 * SearchPhrase(bufPtr, "label .")
	TCLFlag += 2 * SearchPhrase(bufPtr, "bind .")
	TCLFlag += 2 * SearchPhrase(bufPtr, "radiobutton .")
	TCLFlag += 2 * SearchPhrase(bufPtr, "scale .")
	TCLFlag += 2 * SearchPhrase(bufPtr, "frame .")
	TCLFlag += SearchPhrase(bufPtr, "menu .")
	TCLFlag += 2 * SearchPhrase(bufPtr, "menubutton .")
	TCLFlag += 2 * SearchPhrase(bufPtr, "focus .")
	TCLFlag += 5 * SearchPhrase(bufPtr, "tk_menuBar .")
	TCLFlag += 2 * SearchPhrase(bufPtr, "geometry .")
	TCLFlag += SearchPhrase(bufPtr, "title .")
	TCLFlag += 2 * SearchPhrase(bufPtr, "tkwait")
	TCLFlag += 2 * SearchPhrase(bufPtr, "blink .")
}

func SetAsrFlag(bufPtr int) {
	if MFlag != 0 {
		return
	}
	ASRFlag = CFlag
	ASRFlag += 10 * SearchPhrase(bufPtr, "")
	ASRFlag += SearchPhrase(bufPtr, "fujitsu lab")
	ASRFlag += SearchPhrase(bufPtr, "computer algebra")
	ASRFlag += SearchPhrase(bufPtr, "algebraic")
	ASRFlag += SearchPhrase(bufPtr, "symbolic")
	ASRFlag += SearchPhrase(bufPtr, "ctrl(\"")
	ASRFlag += SearchPhrase(bufPtr, "} def")
	ASRFlag += 2 * SearchPhrase(bufPtr, "*/ def")
	ASRFlag += SearchPhrase(bufPtr, "pari(")
	ASRFlag += SearchPhrase(bufPtr, "setprec(")
	ASRFlag += SearchPhrase(bufPtr, "newvect(")
	ASRFlag += SearchPhrase(bufPtr, "newmat(")
	ASRFlag += SearchPhrase(bufPtr, "functor(")
	ASRFlag += SearchPhrase(bufPtr, "funargs(")
	ASRFlag += SearchPhrase(bufPtr, "rtostr(")
	ASRFlag += SearchPhrase(bufPtr, "strtov(")
	ASRFlag += SearchPhrase(bufPtr, "which(\"")
	ASRFlag += SearchPhrase(bufPtr, "plotover(")
	ASRFlag += SearchPhrase(bufPtr, "conplot(")
	ASRFlag += SearchPhrase(bufPtr, "ifplot(")
	ASRFlag += SearchPhrase(bufPtr, "grm(")
	ASRFlag += SearchPhrase(bufPtr, "hgr(")
	ASRFlag += SearchPhrase(bufPtr, "hgrm(")
	ASRFlag += SearchPhrase(bufPtr, "newalg(")
	ASRFlag += SearchPhrase(bufPtr, "defpoly(")
	ASRFlag += SearchPhrase(bufPtr, "ratint(")
	if ASRFlag == CFlag {
		ASRFlag = 0
	}
}

func SetMacFlag(bufPtr int) {
	MACFlag = 10 * SearchPhrase(bufPtr, "")
	MACFlag += 10 * SearchPhrase(bufPtr, "")
	MACFlag += 5 * SearchPhrase(bufPtr, "?round")
	MACFlag += 5 * SearchPhrase(bufPtr, "?truncate")
	MACFlag += 5 * SearchPhrase(bufPtr, "::= ")
	MACFlag += SearchPhrase(bufPtr, ":=")
	MACFlag += SearchPhrase(bufPtr, "cabs(")
	MACFlag += 2 * SearchPhrase(bufPtr, "declare(")
	MACFlag += 5 * SearchPhrase(bufPtr, "define_variable(")
	MACFlag += SearchPhrase(bufPtr, "diff(")
	MACFlag += SearchPhrase(bufPtr, "eigenvals(")
	MACFlag += SearchPhrase(bufPtr, "eigenvects(")
	MACFlag += 5 * SearchPhrase(bufPtr, "eval_when")
	MACFlag += 3 * SearchPhrase(bufPtr, "gendiff(")
	MACFlag += SearchPhrase(bufPtr, "graph2(")
	MACFlag += SearchPhrase(bufPtr, "graph3d(")
	MACFlag += 3 * SearchPhrase(bufPtr, "file_")
	MACFlag += 3 * SearchPhrase(bufPtr, "grobner_basis")
	MACFlag += 5 * SearchPhrase(bufPtr, "mode_check_")
	MACFlag += 5 * SearchPhrase(bufPtr, "mode_declare(")
	MACFlag += SearchPhrase(bufPtr, "plot2(")
	MACFlag += SearchPhrase(bufPtr, "plot3d(")
	MACFlag += 3 * SearchPhrase(bufPtr, "resolvante_")
	MACFlag += 5 * SearchPhrase(bufPtr, "setup_autoload(")
	MACFlag += SearchPhrase(bufPtr, "solve(")
	MACFlag += SearchPhrase(bufPtr, "sum(")
	MACFlag += SearchPhrase(bufPtr, "taylor(")
	MACFlag += SearchPhrase(bufPtr, "taylor_")
	MACFlag += SearchPhrase(bufPtr, "tr_")
	MACFlag += 3 * SearchPhrase(bufPtr, "undiff(")
	MACFlag += 5 * SearchPhrase(bufPtr, "with_stdout(")
	MACFlag += 3 * SearchPhrase(bufPtr, "writefile(")
}

func SetMapFlag(bufPtr int) {
	for bPtr := bufPtr; bPtr != BufferSize/2; bPtr++ {
		if Buffer[bPtr] == '\n' && Buffer[(bPtr+1)%BufferSize] == '#' && (Buffer[(bPtr+2)%BufferSize] == '\t' || Buffer[(bPtr+2)%BufferSize] == ' ') {
			MAPFlag++
		}
	}
	MAPFlag += 10 * SearchPhrase(bufPtr, "# maple")
	MAPFlag += SearchPhrase(bufPtr, "##")
	MAPFlag += SearchPhrase(bufPtr, "###")
	MAPFlag += 2 * SearchPhrase(bufPtr, "####")
	MAPFlag += SearchPhrase(bufPtr, "university of waterloo")
	MAPFlag += SearchPhrase(bufPtr, "computer algebra")
	MAPFlag += SearchPhrase(bufPtr, "algebraic")
	MAPFlag += SearchPhrase(bufPtr, "symbolic")
	MAPFlag += SearchPhrase(bufPtr, "mvcal")
	MAPFlag += SearchPhrase(bufPtr, "daub")
	MAPFlag += SearchPhrase(bufPtr, "calcp")
	MAPFlag += SearchPhrase(bufPtr, "calcplot")
	MAPFlag += SearchPhrase(bufPtr, "posets")
	MAPFlag += SearchPhrase(bufPtr, "coxeter")
	MAPFlag += SearchPhrase(bufPtr, "matthews")
	MAPFlag += SearchPhrase(bufPtr, "casa")
	MAPFlag += SearchPhrase(bufPtr, "macro(")
	MAPFlag += SearchPhrase(bufPtr, "with(")
	MAPFlag += 4 * SearchPhrase(bufPtr, "fi;")
	MAPFlag += 4 * SearchPhrase(bufPtr, "od;")
	MAPFlag += 4 * SearchPhrase(bufPtr, "error(`")
	MAPFlag += 3 * SearchPhrase(bufPtr, "`)")
	MAPFlag += SearchPhrase(bufPtr, "='")
	MAPFlag += 3 * SearchPhrase(bufPtr, "=`")
	MAPFlag += 3 * SearchPhrase(bufPtr, "`=")
	MAPFlag += SearchPhrase(bufPtr, ":=")
	MAPFlag += 3 * SearchPhrase(bufPtr, "`:=")
	MAPFlag += 3 * SearchPhrase(bufPtr, ":=`")
	MAPFlag += SearchPhrase(bufPtr, "<>")
	MAPFlag += SearchPhrase(bufPtr, "proc(")
	MAPFlag += 2 * SearchPhrase(bufPtr, "readlib(")
	MAPFlag += 4 * SearchPhrase(bufPtr, "readlib(`")
	MAPFlag += SearchPhrase(bufPtr, "eigenvals(")
	MAPFlag += SearchPhrase(bufPtr, "eigenvects(")
	MAPFlag += SearchPhrase(bufPtr, "simplify(")
	MAPFlag += SearchPhrase(bufPtr, "sum(")
	MAPFlag += SearchPhrase(bufPtr, "diff(")
	MAPFlag += SearchPhrase(bufPtr, "int(")
	MAPFlag += SearchPhrase(bufPtr, "solve(")
	MAPFlag += SearchPhrase(bufPtr, "draw(")
	MAPFlag += SearchPhrase(bufPtr, "plot(")
	MAPFlag += SearchPhrase(bufPtr, "plot3d(")
	MAPFlag += SearchPhrase(bufPtr, "animate(")
	MAPFlag += SearchPhrase(bufPtr, "animate3d(")
	MAPFlag += 4 * SearchPhrase(bufPtr, "save `")
}

func SetMatFlag(bufPtr int) {
	for bPtr := bufPtr; bPtr != BufferSize/2; bPtr++ {
		if Buffer[bPtr] == '\n' && Buffer[(bPtr+1)%BufferSize] == '%' && (Buffer[(bPtr+2)%BufferSize] == '\t' || Buffer[(bPtr+2)%BufferSize] == ' ') {
			MATFlag++
		}
	}
	MATFlag += 10 * SearchPhrase(bufPtr, "(* mathematica *)")
	MATFlag += 2 * SearchPhrase(bufPtr, "*)(*")
	MATFlag += SearchPhrase(bufPtr, "wolfram research")
	MATFlag += SearchPhrase(bufPtr, "computer algebra")
	MATFlag += SearchPhrase(bufPtr, "algebraic")
	MATFlag += SearchPhrase(bufPtr, "symbolic")
	MATFlag += 3 * SearchPhrase(bufPtr, "beginpackage[")
	MATFlag += 3 * SearchPhrase(bufPtr, "endpackage[")
	MATFlag += 3 * SearchPhrase(bufPtr, "begin[")
	MATFlag += 3 * SearchPhrase(bufPtr, "module[")
	MATFlag += SearchPhrase(bufPtr, "block[")
	MATFlag += 3 * SearchPhrase(bufPtr, "end[")
	MATFlag += SearchPhrase(bufPtr, "on[")
	MATFlag += SearchPhrase(bufPtr, "off[")
	MATFlag += SearchPhrase(bufPtr, "array[")
	MATFlag += SearchPhrase(bufPtr, "table[")
	MATFlag += SearchPhrase(bufPtr, "vectorq[")
	MATFlag += SearchPhrase(bufPtr, "matrixq[")
	MATFlag += SearchPhrase(bufPtr, "list[")
	MATFlag += 2 * SearchPhrase(bufPtr, "evaluate[")
	MATFlag += SearchPhrase(bufPtr, "function[")
	MATFlag += 4 * SearchPhrase(bufPtr, "if[")
	MATFlag += 4 * SearchPhrase(bufPtr, "switch[")
	MATFlag += 4 * SearchPhrase(bufPtr, "do[")
	MATFlag += 4 * SearchPhrase(bufPtr, "while[")
	MATFlag += 4 * SearchPhrase(bufPtr, "for[")
	MATFlag += 4 * SearchPhrase(bufPtr, "break[")
	MATFlag += 4 * SearchPhrase(bufPtr, "continue[")
	MATFlag += 4 * SearchPhrase(bufPtr, "return[")
	MATFlag += 4 * SearchPhrase(bufPtr, "label[")
	MATFlag += 4 * SearchPhrase(bufPtr, "goto[")
	MATFlag += 2 * SearchPhrase(bufPtr, "sum[")
	MATFlag += 2 * SearchPhrase(bufPtr, "product[")
	MATFlag += 3 * SearchPhrase(bufPtr, "expand[")
	MATFlag += 2 * SearchPhrase(bufPtr, "factor[")
	MATFlag += 2 * SearchPhrase(bufPtr, "simplify[")
	MATFlag += 2 * SearchPhrase(bufPtr, "limit[")
	MATFlag += 2 * (SearchPhrase(bufPtr, "d[") + SearchPhrase(bufPtr, "dt["))
	MATFlag += 3 * SearchPhrase(bufPtr, "integrate[")
	MATFlag += 2 * SearchPhrase(bufPtr, "solve[")
	MATFlag += SearchPhrase(bufPtr, "series[")
	MATFlag += SearchPhrase(bufPtr, "show[")
	MATFlag += SearchPhrase(bufPtr, "plot[")
	MATFlag += SearchPhrase(bufPtr, "plot3d[")
}

func SetMlabFlag(bufPtr int) {
	MLABFlag += 10 * SearchPhrase(bufPtr, "# matlab")
	MLABFlag += 10 * SearchPhrase(bufPtr, "% matlab")
	MLABFlag += 10 * SearchPhrase(bufPtr, "# octave")
	MLABFlag += 10 * SearchPhrase(bufPtr, "% octave")
	MLABFlag += SearchPhrase(bufPtr, "# mathworks")
	MLABFlag += SearchPhrase(bufPtr, "% mathworks")
	MLABFlag += 3 * SearchPhrase(bufPtr, "endfor")
	MLABFlag += SearchPhrase(bufPtr, "function[")
	MLABFlag += 3 * SearchPhrase(bufPtr, "endfunction")
	MLABFlag += 3 * SearchPhrase(bufPtr, "endwhile")
	MLABFlag += 2 * SearchPhrase(bufPtr, "end end")
	MLABFlag += SearchPhrase(bufPtr, "plot(")
	MLABFlag += SearchPhrase(bufPtr, "plot3d(")
	MLABFlag += 2 * SearchPhrase(bufPtr, "surf(")
	MLABFlag += 2 * SearchPhrase(bufPtr, "shg")
	MLABFlag += 2 * SearchPhrase(bufPtr, "clg")
	MLABFlag += SearchPhrase(bufPtr, "gplot")
	MLABFlag += SearchPhrase(bufPtr, "gsplot")
	MLABFlag += 3 * SearchPhrase(bufPtr, "casesen")
	MLABFlag += 3 * SearchPhrase(bufPtr, "edit_history")
	MLABFlag += 3 * SearchPhrase(bufPtr, "run_history")
	MLABFlag += 2 * SearchPhrase(bufPtr, ".+")
	MLABFlag += 2 * SearchPhrase(bufPtr, ".-")
	MLABFlag += 3 * SearchPhrase(bufPtr, ".*")
	MLABFlag += 2 * SearchPhrase(bufPtr, "./")
	MLABFlag += 3 * SearchPhrase(bufPtr, ".\\")
	MLABFlag += 3 * SearchPhrase(bufPtr, ".**")
	MLABFlag += 3 * SearchPhrase(bufPtr, ".^")
	MLABFlag += 3 * SearchPhrase(bufPtr, "~=")
	MLABFlag += 4 * SearchPhrase(bufPtr, ".'")
	MLABFlag += SearchPhrase(bufPtr, "<>")
	MLABFlag += SearchPhrase(bufPtr, "printf(\"")
	MLABFlag += 4 * SearchPhrase(bufPtr, "printf('")
	MLABFlag += 4 * SearchPhrase(bufPtr, "disp('")
	MLABFlag += SearchPhrase(bufPtr, "error(\"")
	MLABFlag += 4 * SearchPhrase(bufPtr, "error('")
	MLABFlag += 3 * SearchPhrase(bufPtr, "nargin")
	MLABFlag += 3 * SearchPhrase(bufPtr, "nargout")
	MLABFlag += 4 * SearchPhrase(bufPtr, "(:")
	MLABFlag += 4 * SearchPhrase(bufPtr, ":)")
	MLABFlag += SearchPhrase(bufPtr, "zeros(")
	MLABFlag += 3 * SearchPhrase(bufPtr, "rot90")
	MLABFlag += 3 * SearchPhrase(bufPtr, "fliplr")
	MLABFlag += 3 * SearchPhrase(bufPtr, "flipud")
	MLABFlag += SearchPhrase(bufPtr, "diag")
	MLABFlag += 3 * SearchPhrase(bufPtr, "tril")
	MLABFlag += 3 * SearchPhrase(bufPtr, "triu")
	MLABFlag += 2 * SearchPhrase(bufPtr, "title('")
	MLABFlag += 3 * SearchPhrase(bufPtr, "gtext(")
	MLABFlag += 3 * SearchPhrase(bufPtr, "ginput(")
	MLABFlag += SearchPhrase(bufPtr, "input(")
}

func SetMpadFlag(bufPtr int) {
	MPADFlag += 10 * SearchPhrase(bufPtr, "# mupad #")
	MPADFlag += SearchPhrase(bufPtr, "##")
	MPADFlag += SearchPhrase(bufPtr, "###")
	MPADFlag += 2 * SearchPhrase(bufPtr, "####")
	MPADFlag += SearchPhrase(bufPtr, "computer algebra")
	MPADFlag += SearchPhrase(bufPtr, "algebraic")
	MPADFlag += SearchPhrase(bufPtr, "symbolic")
	MPADFlag += 5 * SearchPhrase(bufPtr, "end_case")
	MPADFlag += 5 * SearchPhrase(bufPtr, "end_for")
	MPADFlag += 5 * SearchPhrase(bufPtr, "end_if")
	MPADFlag += 3 * SearchPhrase(bufPtr, "parbegin")
	MPADFlag += 5 * SearchPhrase(bufPtr, "end_par")
	MPADFlag += 5 * SearchPhrase(bufPtr, "end_proc")
	MPADFlag += 5 * SearchPhrase(bufPtr, "end_repeat")
	MPADFlag += 3 * SearchPhrase(bufPtr, "seqbegin")
	MPADFlag += 5 * SearchPhrase(bufPtr, "end_seq")
	MPADFlag += 5 * SearchPhrase(bufPtr, "end_while")
	MPADFlag += SearchPhrase(bufPtr, ":=")
	MPADFlag += SearchPhrase(bufPtr, "diff(")
	MPADFlag += SearchPhrase(bufPtr, "int(")
	MPADFlag += SearchPhrase(bufPtr, "solve(")
	MPADFlag += SearchPhrase(bufPtr, "plot2d(")
	MPADFlag += SearchPhrase(bufPtr, "plot3d(")
	MPADFlag += SearchPhrase(bufPtr, "simplify(")
	MPADFlag += SearchPhrase(bufPtr, "sum(")
}

func SetRedFlag(bufPtr int) {
	if MFlag != 0 {
		return
	}
	REDFlag += 10 * SearchPhrase(bufPtr, "% reduce")
	REDFlag += SearchPhrase(bufPtr, ". hearn")
	REDFlag += SearchPhrase(bufPtr, "computer algebra")
	REDFlag += SearchPhrase(bufPtr, "algebraic")
	REDFlag += SearchPhrase(bufPtr, "symbolic")
	REDFlag += SearchPhrase(bufPtr, "arnum")
	REDFlag += SearchPhrase(bufPtr, "compact")
	REDFlag += SearchPhrase(bufPtr, "excalc")
	REDFlag += SearchPhrase(bufPtr, "gentran")
	REDFlag += SearchPhrase(bufPtr, "orthovec")
	REDFlag += SearchPhrase(bufPtr, "specfn")
	REDFlag += SearchPhrase(bufPtr, "fide")
	REDFlag += SearchPhrase(bufPtr, "physop")
	REDFlag += SearchPhrase(bufPtr, "reacteqn")
	REDFlag += SearchPhrase(bufPtr, "rlfi")
	REDFlag += SearchPhrase(bufPtr, "ghyper")
	REDFlag += SearchPhrase(bufPtr, "linalg")
	REDFlag += SearchPhrase(bufPtr, "ncpoly")
	REDFlag += SearchPhrase(bufPtr, "xideal")
	REDFlag += SearchPhrase(bufPtr, "zeilberg")
	REDFlag += SearchPhrase(bufPtr, "out \"")
	REDFlag += SearchPhrase(bufPtr, "begin scalar")
	REDFlag += (SearchPhrase(bufPtr, "off echo") + SearchPhrase(bufPtr, "on echo"))
	REDFlag += (SearchPhrase(bufPtr, "on rational") + SearchPhrase(bufPtr, "off rational"))
	REDFlag += (SearchPhrase(bufPtr, "on complex") + SearchPhrase(bufPtr, "off complex"))
	REDFlag += (SearchPhrase(bufPtr, "on round") + SearchPhrase(bufPtr, "off round"))
	REDFlag += (SearchPhrase(bufPtr, "on factor") + SearchPhrase(bufPtr, "off factor"))
	REDFlag += (SearchPhrase(bufPtr, "on nat") + SearchPhrase(bufPtr, "off nat"))
	REDFlag += SearchPhrase(bufPtr, ":=")
	REDFlag += SearchPhrase(bufPtr, "part(")
	REDFlag += SearchPhrase(bufPtr, "factorize(")
	REDFlag += SearchPhrase(bufPtr, "remainder(")
	REDFlag += SearchPhrase(bufPtr, "mat((")
	REDFlag += 2 * SearchPhrase(bufPtr, "then<<")
	REDFlag += 2 * SearchPhrase(bufPtr, "do<<")
	REDFlag += 2 * SearchPhrase(bufPtr, "repeat<<")
	REDFlag += 2 * SearchPhrase(bufPtr, "else<<")
	REDFlag += SearchPhrase(bufPtr, "df(")
	REDFlag += SearchPhrase(bufPtr, "int(")
	REDFlag += SearchPhrase(bufPtr, "defint(")
	REDFlag += SearchPhrase(bufPtr, "solve(")
	REDFlag += SearchPhrase(bufPtr, "taylor(")
	REDFlag += SearchPhrase(bufPtr, "groebner(")
	REDFlag += SearchPhrase(bufPtr, "odesolve(")
	REDFlag += SearchPhrase(bufPtr, "root(")
	REDFlag += SearchPhrase(bufPtr, "plot(")
	REDFlag += SearchPhrase(bufPtr, "linineq(")
}

func SetTxtFlag(bufPtr int) {
	if SearchPhrase(bufPtr, "from:") != 0 && SearchPhrase(bufPtr, "newsgroups:") != 0 && SearchPhrase(bufPtr, "subject:") != 0 && SearchPhrase(bufPtr, "date:") != 0 && SearchPhrase(bufPtr, "organization:") != 0 && SearchPhrase(bufPtr, "path:") != 0 {
		TXTFlag = 1
	}
	if SearchPhrase(bufPtr, "from:") != 0 && SearchPhrase(bufPtr, "subject:") != 0 && SearchPhrase(bufPtr, "date:") != 0 && SearchPhrase(bufPtr, "to:") != 0 && SearchPhrase(bufPtr, "return-path:") != 0 {
		TXTFlag = 1
	}
}

func SetLangFlag(bufPtr int) {
}

func SetPythonFlag(bufPtr int) {
	PythonFlag += 10 * SearchPhrase(bufPtr, "# python")
}

func SetRubyFlag(bufPtr int) {
	RubyFlag += 10 * SearchPhrase(bufPtr, "# ruby")
}

func SetRustFlag(bufPtr int) {
	RustFlag += 10 * SearchPhrase(bufPtr, "// rust")
}

func SetGoFlag(bufPtr int) {
	GoFlag += 10 * SearchPhrase(bufPtr, "// go")
}

func SetJsFlag(bufPtr int) {
	JSFlag += 10 * SearchPhrase(bufPtr, "// javascript")
}

func SetTsFlag(bufPtr int) {
	TSFlag += 10 * SearchPhrase(bufPtr, "// typescript")
}

func SetKotlinFlag(bufPtr int) {
	KotlinFlag += 10 * SearchPhrase(bufPtr, "// kotlin")
}

func SetSwiftFlag(bufPtr int) {
	SwiftFlag += 10 * SearchPhrase(bufPtr, "// swift")
}

func InitLangFlag() {
	bufPtr := 0
	ext := strings.ToLower(filepath.Ext(inputFileName))

	if ext == ".tex" || ext == ".txt" {
		TXTFlag = 1
		return
	}
	if ext == ".bas" || ext == ".vb" {
		BASFlag = 1
		return
	}
	if ext == ".c" || ext == ".cpp" || ext == ".vc" || ext == ".h" || ext == ".hpp" {
		CFlag = 1
		return
	}
	if ext == ".cbl" || ext == ".cob" {
		CBLFlag = 1
		return
	}
	if ext == ".f" || ext == ".for" {
		F77Flag = 1
		return
	}
	if ext == ".html" {
		HTMLFlag = 1
		return
	}
	if ext == ".java" {
		JAVAFlag = 1
		return
	}
	if ext == ".el" || ext == ".lsp" || ext == ".sc" || ext == ".scm" {
		LISPFlag = 1
		return
	}
	if strings.ToLower(filepath.Base(inputFileName)) == "makefile" {
		MAKEFlag = 1
		return
	}
	if ext == ".p" || ext == ".pas" || ext == ".tp" {
		PASFlag = 1
		return
	}
	if ext == ".pl" || ext == ".prl" {
		PERLFlag = 1
		return
	}
	if ext == ".sh" || ext == ".csh" || ext == ".ksh" {
		SHFlag = 1
		return
	}
	if ext == ".tcl" || ext == ".tk" {
		TCLFlag = 1
		return
	}
	if ext == ".asi" || ext == ".asir" || ext == ".asr" {
		ASRFlag = 1
		return
	}
	if ext == ".mac" || ext == ".max" {
		MACFlag = 1
		return
	}
	if ext == ".map" || ext == ".mpl" {
		MAPFlag = 1
		return
	}
	if ext == ".mat" || ext == ".mma" {
		MATFlag = 1
		return
	}
	if ext == ".ml" || ext == ".mtlb" || ext == ".oct" {
		MLABFlag = 1
		return
	}
	if ext == ".mu" {
		MPADFlag = 1
		return
	}
	if ext == ".red" || ext == ".rdc" {
		REDFlag = 1
		return
	}
	if ext == ".m" || ext == ".M" {
		MFlag = 1
		return
	}
	if ext == ".py" {
		PythonFlag = 1
		return
	}
	if ext == ".rb" {
		RubyFlag = 1
		return
	}
	if ext == ".rs" {
		RustFlag = 1
		return
	}
	if ext == ".go" {
		GoFlag = 1
		return
	}
	if ext == ".js" {
		JSFlag = 1
		return
	}
	if ext == ".ts" {
		TSFlag = 1
		return
	}
	if ext == ".kt" {
		KotlinFlag = 1
		return
	}
	if ext == ".swift" {
		SwiftFlag = 1
		return
	}

	SetBasFlag(bufPtr)
	SetCFlag(bufPtr)
	SetCblFlag(bufPtr)
	SetF77Flag(bufPtr)
	SetHtmlFlag(bufPtr)
	SetJavaFlag(bufPtr)
	SetLispFlag(bufPtr)
	SetMakeFlag(bufPtr)
	SetPasFlag(bufPtr)
	SetPerlFlag(bufPtr)
	SetShellFlag(bufPtr)
	SetTclFlag(bufPtr)

	SetAsrFlag(bufPtr)
	SetMacFlag(bufPtr)
	SetMapFlag(bufPtr)
	SetMatFlag(bufPtr)
	SetMlabFlag(bufPtr)
	SetMpadFlag(bufPtr)
	SetRedFlag(bufPtr)

	SetPythonFlag(bufPtr)
	SetRubyFlag(bufPtr)
	SetRustFlag(bufPtr)
	SetGoFlag(bufPtr)
	SetJsFlag(bufPtr)
	SetTsFlag(bufPtr)
	SetKotlinFlag(bufPtr)
	SetSwiftFlag(bufPtr)

	maxFlag := BASFlag
	maxName := "BAS"

	flags := []struct {
		name string
		val  *int
	}{
		{"C", &CFlag}, {"CBL", &CBLFlag}, {"F77", &F77Flag},
		{"HTML", &HTMLFlag}, {"JAVA", &JAVAFlag}, {"LISP", &LISPFlag},
		{"MAKE", &MAKEFlag}, {"PAS", &PASFlag}, {"PERL", &PERLFlag},
		{"SH", &SHFlag}, {"TCL", &TCLFlag}, {"ASR", &ASRFlag},
		{"MAC", &MACFlag}, {"MAP", &MAPFlag}, {"MAT", &MATFlag},
		{"MLAB", &MLABFlag}, {"MPAD", &MPADFlag}, {"RED", &REDFlag},
		{"PYTHON", &PythonFlag}, {"RUBY", &RubyFlag}, {"RUST", &RustFlag},
		{"GO", &GoFlag}, {"JS", &JSFlag}, {"TS", &TSFlag},
		{"KOTLIN", &KotlinFlag}, {"SWIFT", &SwiftFlag},
	}

	for _, f := range flags {
		if *f.val > maxFlag {
			maxFlag = *f.val
			maxName = f.name
		}
	}

	if maxFlag == 0 {
		if MFlag != 0 {
			MLABFlag = 1
		} else {
			TXTFlag = 1
		}
	} else {
		BASFlag = 0
		CFlag = 0
		CBLFlag = 0
		F77Flag = 0
		HTMLFlag = 0
		JAVAFlag = 0
		LISPFlag = 0
		MAKEFlag = 0
		PASFlag = 0
		PERLFlag = 0
		SHFlag = 0
		TCLFlag = 0
		ASRFlag = 0
		MACFlag = 0
		MAPFlag = 0
		MATFlag = 0
		MLABFlag = 0
		MPADFlag = 0
		REDFlag = 0
		PythonFlag = 0
		RubyFlag = 0
		RustFlag = 0
		GoFlag = 0
		JSFlag = 0
		TSFlag = 0
		KotlinFlag = 0
		SwiftFlag = 0

		switch maxName {
		case "BAS":
			BASFlag = 1
		case "C":
			CFlag = 1
		case "CBL":
			CBLFlag = 1
		case "F77":
			F77Flag = 1
		case "HTML":
			HTMLFlag = 1
		case "JAVA":
			JAVAFlag = 1
		case "LISP":
			LISPFlag = 1
		case "MAKE":
			MAKEFlag = 1
		case "PAS":
			PASFlag = 1
		case "PERL":
			PERLFlag = 1
		case "SH":
			SHFlag = 1
		case "TCL":
			TCLFlag = 1
		case "ASR":
			ASRFlag = 1
		case "MAC":
			MACFlag = 1
		case "MAP":
			MAPFlag = 1
		case "MAT":
			MATFlag = 1
		case "MLAB":
			MLABFlag = 1
		case "MPAD":
			MPADFlag = 1
		case "RED":
			REDFlag = 1
		case "PYTHON":
			PythonFlag = 1
		case "RUBY":
			RubyFlag = 1
		case "RUST":
			RustFlag = 1
		case "GO":
			GoFlag = 1
		case "JS":
			JSFlag = 1
		case "TS":
			TSFlag = 1
		case "KOTLIN":
			KotlinFlag = 1
		case "SWIFT":
			SwiftFlag = 1
		}
	}
}
