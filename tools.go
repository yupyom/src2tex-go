package main

import (
	"fmt"
	"os"
	"strings"
)

// readCustomFile reads the entire content of a file and returns it as a string.
// If the file cannot be read, it prints an error and exits.
func readCustomFile(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: cannot read custom file '%s': %v\n", path, err)
		os.Exit(1)
	}
	content := string(data)
	// Ensure trailing newline
	if !strings.HasSuffix(content, "\n") {
		content += "\n"
	}
	return content
}

// merge_ntt_ascii absorbs differences of NTT and ASCII JTeXs
func MergeNttAscii(outFile *os.File) {
	if IsUnicodeTeX {
		fmt.Fprintf(outFile, "\\def\\mc{\\relax}\n")
		fmt.Fprintf(outFile, "\\def\\gt{\\relax}\n")
		fmt.Fprintf(outFile, "\\def\\sc{\\scshape}\n")
		return
	}

	// Original C code had #ifdef ASCII. Here we default to ASCII JTeX style assuming modern environments or no-op
	// but we port the ASCII block verbatim.
	fmt.Fprintf(outFile, "\\ifx\\gtfam\\undefined\n")
	fmt.Fprintf(outFile, "  \\ifx\\dm\\undefined\n")
	fmt.Fprintf(outFile, "    \\ifx\\tendm\\undefined\n")
	fmt.Fprintf(outFile, "      \\def\\mc{\\null}\n")
	fmt.Fprintf(outFile, "    \\else\n")
	fmt.Fprintf(outFile, "      \\def\\mc{\\tendm}\n")
	fmt.Fprintf(outFile, "    \\fi\n")
	fmt.Fprintf(outFile, "  \\else\n")
	fmt.Fprintf(outFile, "    \\def\\mc{\\dm}\n")
	fmt.Fprintf(outFile, "  \\fi\n")
	fmt.Fprintf(outFile, "  \\ifx\\dg\\undefined\n")
	fmt.Fprintf(outFile, "    \\ifx\\tendg\\undefined\n")
	fmt.Fprintf(outFile, "      \\def\\gt{\\null}\n")
	fmt.Fprintf(outFile, "    \\else\n")
	fmt.Fprintf(outFile, "      \\def\\gt{\\tendg}\n")
	fmt.Fprintf(outFile, "    \\fi\n")
	fmt.Fprintf(outFile, "  \\else\n")
	fmt.Fprintf(outFile, "    \\def\\gt{\\dg}\n")
	fmt.Fprintf(outFile, "  \\fi\n")
	fmt.Fprintf(outFile, "\\fi\n")
	fmt.Fprintf(outFile, "\\ifx\\sc\\undefined\n")
	fmt.Fprintf(outFile, "  \\def\\sc{\\null}\n")
	fmt.Fprintf(outFile, "\\fi\n")
}

var gfPrevFlag = 0
var gFlagChar FlagChar // Single shared FlagChar like static in C

func GetFlagChar(inFile *os.File) *FlagChar {
	bufPtr := Fgetc2Buffer(inFile)
	flag := GetCommentFlag(bufPtr)
	if gfPrevFlag != 0 && flag != 0 && TXTFlag == 0 {
		if GetTexFlag(bufPtr) != 0 {
			flag++
		}
	}
	gFlagChar.Flag = flag
	gFlagChar.Character = Buffer[bufPtr]
	gFlagChar.BufferIdx = bufPtr
	gfPrevFlag = flag
	return &gFlagChar
}


// fprintfLegacyFallback outputs a minimal legacy preamble when font lookup fails.
func fprintfLegacyFallback(outFile *os.File) {
	fmt.Fprintf(outFile, "\\usepackage{iftex}\n")
	fmt.Fprintf(outFile, "\\ifXeTeX\n")
	fmt.Fprintf(outFile, "  \\usepackage{fontspec}\n")
	fmt.Fprintf(outFile, "  \\usepackage{xeCJK}\n")
	fmt.Fprintf(outFile, "  \\usepackage[haranoaji]{zxjafont}\n")
	fmt.Fprintf(outFile, "  \\XeTeXlinebreaklocale \"\"\n")
	fmt.Fprintf(outFile, "\\fi\n")
	fmt.Fprintf(outFile, "\\ifLuaTeX\n")
	fmt.Fprintf(outFile, "  \\usepackage[haranoaji]{luatexja-preset}\n")
	fmt.Fprintf(outFile, "\\fi\n")
}

// fprintfUnicodePreamble outputs the Unicode LaTeX preamble with all necessary
// packages and plain TeX compatibility macros.
// The preamble uses iftex to auto-detect the TeX engine so that the same .tex
// file compiles correctly with both XeLaTeX and LuaLaTeX.
func fprintfUnicodePreamble(fileName string, outFile *os.File) {
	fmt.Fprintf(outFile, "\\documentclass")
	if LatexOption != "" {
		fmt.Fprintf(outFile, "[%s]", LatexOption)
	}
	fmt.Fprintf(outFile, "{%s}\n", LatexStyle)

	// If -font was not explicitly specified, use simple legacy fallback.
	// This makes -unicode and -font independent options.
	if !FontExplicitlySet && CommentFontName == "" {
		fprintfLegacyFallback(outFile)
	} else {
		// Look up the selected font
		fd := LookupFont(CJKFontName)
		if fd == nil {
			fmt.Fprintf(os.Stderr, "\nWarning: unknown font '%s'. Falling back to legacy mode.\n", CJKFontName)
			fmt.Fprintf(os.Stderr, "Use -font list to see available fonts.\n")
			fprintfLegacyFallback(outFile)
		} else {
			// Check if the font is installed (skip for ipaex which is in TeX Live)
			if fd.Name != "ipaex" && !IsFontInstalled(fd.Name) {
				fmt.Fprintf(os.Stderr, "\nWarning: font '%s' (%s) is not installed.\n", fd.Name, fd.DisplayName)
				fmt.Fprintf(os.Stderr, "Install it with: src2latexg -font install %s\n", fd.Name)
				fmt.Fprintf(os.Stderr, "Falling back to IPAex.\n")
				fd = GetFontDef("ipaex")
			}

			// For ipaex, check TeX Live availability
			if fd.Name == "ipaex" && !IsIPAexAvailable() {
				fmt.Fprintf(os.Stderr, "\nWarning: IPAex font not found in TeX Live.\n")
				fmt.Fprintf(os.Stderr, "For proper CJK monospace alignment, install a monospace font.\n")
				fmt.Fprintf(os.Stderr, "See README.md for installation instructions, or run:\n")
				fmt.Fprintf(os.Stderr, "  src2latexg -font list\n")
				fmt.Fprintf(os.Stderr, "  src2latexg -font install hackgen\n")
				fmt.Fprintf(os.Stderr, "Falling back to legacy (Harano Aji) preamble.\n\n")
				fprintfLegacyFallback(outFile)
			} else {
				// Set unified font mode flag for text2tex output
				UseUnifiedFont = fd.Unified

				// Auto-detect mincho font for comments if not explicitly specified
				AutoDetectCommentFont()

				// Engine detection: XeLaTeX vs LuaLaTeX with selected font
				fmt.Fprintf(outFile, "\\usepackage{iftex}\n")
				fmt.Fprintf(outFile, "\\ifXeTeX\n")
				fmt.Fprintf(outFile, "%s", GenerateFontPreambleXeLaTeX(fd))
				fmt.Fprintf(outFile, "\\fi\n")
				fmt.Fprintf(outFile, "\\ifLuaTeX\n")
				fmt.Fprintf(outFile, "%s", GenerateFontPreambleLuaLaTeX(fd))
				fmt.Fprintf(outFile, "\\fi\n")
			}
		}
	}

	fmt.Fprintf(outFile, "\\usepackage[a4paper, margin=2cm]{geometry}\n")
	fmt.Fprintf(outFile, "\\usepackage{graphicx}\n")
	// Define \charwd = width of one \tt character for robust monospace alignment
	fmt.Fprintf(outFile, "\\newdimen\\charwd\n")
	fmt.Fprintf(outFile, "{\\tt\\global\\setbox0=\\hbox{x}\\global\\charwd=\\wd0}\n")
	// Uniform inter-word spacing (no extra space after periods)
	fmt.Fprintf(outFile, "\\frenchspacing\n")
	if HeaderFile != "" {
		// Custom header: user provides entire fancyhdr header block
		fmt.Fprintf(outFile, "\\usepackage{fancyhdr}\n")
		fmt.Fprintf(outFile, "\\pagestyle{fancy}\n")
		fmt.Fprintf(outFile, "%s", readCustomFile(HeaderFile))
	} else {
		fmt.Fprintf(outFile, "\\usepackage{fancyhdr}\n")
		fmt.Fprintf(outFile, "\\pagestyle{fancy}\n")
		fmt.Fprintf(outFile, "\\renewcommand{\\headrulewidth}{0pt}\n")
		fmt.Fprintf(outFile, "\\fancyhf{}\n")
		fmt.Fprintf(outFile, "\\fancyhead[R]{\\rm src2tex-go version 2.23}\n")
	}
	// plain TeX compatibility macros for LaTeX
	// (do NOT use amsmath — it redefines \cases and \pmatrix in incompatible ways)
	fmt.Fprintf(outFile, "%% plain TeX compatibility\n")
	fmt.Fprintf(outFile, "\\makeatletter\n")
	fmt.Fprintf(outFile, "\\providecommand{\\eqalign}[1]{\\vcenter{\\openup1\\jot\\m@th\n")
	fmt.Fprintf(outFile, "  \\ialign{\\strut\\hfil$\\displaystyle{##}$&$\\displaystyle{{}##}$\\hfil\\crcr#1\\crcr}}}\n")
	fmt.Fprintf(outFile, "\\makeatother\n")

	if FooterFile != "" {
		// Custom footer: user provides entire fancyhdr footer block
		fmt.Fprintf(outFile, "%s", readCustomFile(FooterFile))
	} else {
		// Escape the filename for fancyfoot
		fmt.Fprintf(outFile, "\\fancyfoot[R]{\\rm\\hfill ")
		fprintfEscapedFileName(fileName, outFile)
		fmt.Fprintf(outFile, "\\qquad page \\thepage}\n")
	}
}

// fprintfEscapedFileName writes a TeX-escaped version of the filename.
func fprintfEscapedFileName(fileName string, outFile *os.File) {
	for i := 0; i < len(fileName); i++ {
		b := fileName[i]
		switch b {
		case '"':
			fmt.Fprintf(outFile, "{\\tt \"}")
		case '#':
			fmt.Fprintf(outFile, "{\\tt\\#}")
		case '$':
			fmt.Fprintf(outFile, "{\\tt\\$}")
		case '%':
			fmt.Fprintf(outFile, "{\\tt\\%%}")
		case '&':
			fmt.Fprintf(outFile, "{\\tt\\&}")
		case '*':
			fmt.Fprintf(outFile, "{\\tt *}")
		case '-':
			fmt.Fprintf(outFile, "{\\tt -}")
		case '/':
			fmt.Fprintf(outFile, "{\\tt /}")
		case '<':
			fmt.Fprintf(outFile, "{\\tt <}")
		case '>':
			fmt.Fprintf(outFile, "{\\tt >}")
		case '\\':
			fmt.Fprintf(outFile, "$\\backslash$")
		case '^':
			fmt.Fprintf(outFile, "$\\hat{\\ }$")
		case '_':
			fmt.Fprintf(outFile, "{\\tt\\_}")
		case '{':
			fmt.Fprintf(outFile, "$\\{$")
		case '|':
			fmt.Fprintf(outFile, "{\\tt |}")
		case '}':
			fmt.Fprintf(outFile, "$\\}$")
		case '~':
			fmt.Fprintf(outFile, "$\\tilde{\\ }$")
		default:
			fmt.Fprintf(outFile, "%c", b)
		}
	}
}

// fprintf_documentstyle outputs documentstyle.
// Returns true if a documentstyle line was found and neutralized in Unicode mode
// (the caller should skip the rest of the current line).
func FprintfDocumentstyle(fileName string, bufPtr int, outFile *os.File) bool {
	bPtr := GetPhrase(bufPtr, "{\\documentstyle")
	if bPtr == -1 {
		// No documentstyle found in source — output default
		if IsLatexMode {
			if IsUnicodeTeX {
				fprintfUnicodePreamble(fileName, outFile)
			} else {
				fmt.Fprintf(outFile, "\\documentstyle")
				if LatexOption != "" {
					fmt.Fprintf(outFile, "[%s]", LatexOption)
				}
				fmt.Fprintf(outFile, "{%s}\n", LatexStyle)
			}
		}
		fmt.Fprintf(outFile, "\\begin{document}\n\n")
		fmt.Fprintf(outFile, "\\ifx\\sevenrm\\undefined\n")
		fmt.Fprintf(outFile, "  \\font\\sevenrm=cmr7 scaled \\magstep0\n")
		fmt.Fprintf(outFile, "\\fi\n")
		return false
	} else {
		// Source file contains a documentstyle directive
		if IsUnicodeTeX {
			// In Unicode mode, always use the Unicode preamble
			// (ignore the source's legacy documentstyle)
			fprintfUnicodePreamble(fileName, outFile)
			// Neutralize the entire line containing the documentstyle directive
			// in the buffer so it won't be output again.
			// Search backwards from bPtr to find the start of the line
			startPtr := bPtr
			for startPtr > 0 && Buffer[startPtr-1] != '\n' && Buffer[startPtr-1] != -1 {
				startPtr--
			}
			// Search forward from bPtr to find the end of the line
			tailPtr := bPtr
			for Buffer[tailPtr] >= ' ' {
				tailPtr++
			}
			// Blank the entire line
			for i := startPtr; i < tailPtr; i++ {
				Buffer[i] = 0x20
			}
			fmt.Fprintf(outFile, "\\begin{document}\n\n")
			fmt.Fprintf(outFile, "\\ifx\\sevenrm\\undefined\n")
			fmt.Fprintf(outFile, "  \\font\\sevenrm=cmr7 scaled \\magstep0\n")
			fmt.Fprintf(outFile, "\\fi\n")
			return true
		} else {
			tailPtr := bPtr
			for Buffer[tailPtr] != '}' && Buffer[tailPtr] >= ' ' {
				tailPtr++
			}
			c1 := byte(Buffer[tailPtr-1])
			c2 := byte(Buffer[tailPtr])
			if Buffer[bPtr] == 0x00 || c1 < '0' || (c1 > '9' && c1 < 'A') || (c1 > 'Z' && c1 < 'a') || c1 > 'z' || c2 != '}' {
				if IsLatexMode {
					fmt.Fprintf(outFile, "\\documentstyle")
					if LatexOption != "" {
						fmt.Fprintf(outFile, "[%s]", LatexOption)
					}
					fmt.Fprintf(outFile, "{%s}\n", LatexStyle)
				}
			} else {
				miniBuffer := ""
				for i := 0; i < 255 && Buffer[bPtr] >= ' '; i++ {
					bPtr++
					miniBuffer += string(byte(Buffer[bPtr]))
					switch i {
					case 0:
						Buffer[bPtr] = '\\'
					case 1:
						Buffer[bPtr] = 'n'
					case 2:
						Buffer[bPtr] = 'u'
					case 3:
						Buffer[bPtr] = 'l'
					case 4:
						Buffer[bPtr] = 'l'
					default:
						Buffer[bPtr] = 0x20
					}
					if miniBuffer[len(miniBuffer)-1] == '}' {
						break
					}
				}
				fmt.Fprintf(outFile, "%s\n", miniBuffer)
			}
		}
	}
	fmt.Fprintf(outFile, "\\begin{document}\n\n")
	fmt.Fprintf(outFile, "\\ifx\\sevenrm\\undefined\n")
	fmt.Fprintf(outFile, "  \\font\\sevenrm=cmr7 scaled \\magstep0\n")
	fmt.Fprintf(outFile, "\\fi\n")
	return false
}

// fprintf_footline
func FprintfFootline(fileName string, outFile *os.File) {
	fmt.Fprintf(outFile, "\\footline={\\rm\\hfill ")
	for i := 0; i < len(fileName); i++ {
		b := fileName[i]
		switch b {
		case '"':
			fmt.Fprintf(outFile, "{\\tt \"}")
		case '#':
			fmt.Fprintf(outFile, "{\\tt\\#}")
		case '$':
			fmt.Fprintf(outFile, "{\\tt\\$}")
		case '%':
			fmt.Fprintf(outFile, "{\\tt\\%%}")
		case '&':
			fmt.Fprintf(outFile, "{\\tt\\&}")
		case '*':
			fmt.Fprintf(outFile, "{\\tt *}")
		case '-':
			fmt.Fprintf(outFile, "{\\tt -}")
		case '/':
			fmt.Fprintf(outFile, "{\\tt /}")
		case '<':
			fmt.Fprintf(outFile, "{\\tt <}")
		case '>':
			fmt.Fprintf(outFile, "{\\tt >}")
		case '\\':
			fmt.Fprintf(outFile, "$\\backslash$")
		case '^':
			fmt.Fprintf(outFile, "$\\hat{\\ }$")
		case '_':
			fmt.Fprintf(outFile, "{\\tt\\_}")
		case '{':
			fmt.Fprintf(outFile, "$\\{$")
		case '|':
			fmt.Fprintf(outFile, "{\\tt |}")
		case '}':
			fmt.Fprintf(outFile, "$\\}$")
		case '~':
			fmt.Fprintf(outFile, "$\\tilde{\\ }$")
		default:
			fmt.Fprintf(outFile, "%c", b)
		}
	}
	fmt.Fprintf(outFile, "\\qquad page \\folio}\n")
}

// input_user_style
func InputUserStyle(outFile *os.File) {
	if IsLatexMode {
		if IsUnicodeTeX {
			// In UnicodeTeX mode, fancyhdr is already set up in the preamble.
			// Only try to load user style file, but do NOT override pagestyle.
			fmt.Fprintf(outFile, "\\newread\\MyStyle\n")
			fmt.Fprintf(outFile, "\\openin\\MyStyle=src2latex.s2t\n")
			fmt.Fprintf(outFile, "\\ifeof\\MyStyle\n")
			fmt.Fprintf(outFile, "  \\closein\\MyStyle\n")
			fmt.Fprintf(outFile, "\\else\n")
			fmt.Fprintf(outFile, "  \\input src2latex.s2t\n")
			fmt.Fprintf(outFile, "  \\closein\\MyStyle\n")
			fmt.Fprintf(outFile, "\\fi\n")
		} else {
			fmt.Fprintf(outFile, "\\newread\\MyStyle\n")
			// assume unix
			fmt.Fprintf(outFile, "\\openin\\MyStyle=src2latex.s2t\n")
			fmt.Fprintf(outFile, "\\ifeof\\MyStyle\n")
			fmt.Fprintf(outFile, "  \\closein\\MyStyle\n")
			if HeaderFile != "" {
				// Custom header already handled in preamble
			} else {
				fmt.Fprintf(outFile, "  \\pagestyle{myheadings}\n")
				fmt.Fprintf(outFile, "  \\markboth{\\rm src2tex-go version 2.23}{\\rm src2tex-go version 2.23}\n")
			}
			fmt.Fprintf(outFile, "\\else\n")
			fmt.Fprintf(outFile, "  \\input src2latex.s2t\n")
			fmt.Fprintf(outFile, "  \\closein\\MyStyle\n")
			fmt.Fprintf(outFile, "\\fi\n")
		}
	} else {
		fmt.Fprintf(outFile, "\\newread\\MyStyle\n")
		fmt.Fprintf(outFile, "\\openin\\MyStyle=src2tex.s2t\n")
		fmt.Fprintf(outFile, "\\ifeof\\MyStyle\n")
		fmt.Fprintf(outFile, "  \\closein\\MyStyle\n")
		if HeaderFile != "" {
			fmt.Fprintf(outFile, "%s", readCustomFile(HeaderFile))
		} else {
			fmt.Fprintf(outFile, "  \\headline={\\hfill\\rm src2tex-go version 2.23}\n")
		}
		fmt.Fprintf(outFile, "\\else\n")
		fmt.Fprintf(outFile, "  \\input src2tex.s2t\n")
		fmt.Fprintf(outFile, "  \\closein\\MyStyle\n")
		fmt.Fprintf(outFile, "\\fi\n")
	}
}

// choose_tt_font tries to determine whether to choose cmtt font over cmr font
// The user explicitly requested that all code portions up to the comment's beginning
// whitespace be fixed-width, to prevent jagged comment alignment.
// We force this mathematically by returning 1 globally here, guaranteeing
// identical em-widths for spaces and thus perfect \hfill alignment!
func ChooseTtFont(bufPtr int) int {
	return 1
}

// fprintf_line_number outputs line number
func FprintfLineNumber(outFile *os.File, lineCounter int) {
	if lineCounter < 9 {
		fmt.Fprintf(outFile, "\n\n{\\tt\\noindent\\phantom{00000}%d:\\ }\n", lineCounter+1)
	} else if lineCounter >= 9 && lineCounter < 99 {
		fmt.Fprintf(outFile, "\n\n{\\tt\\noindent\\phantom{0000}%d:\\ }\n", lineCounter+1)
	} else if lineCounter >= 99 && lineCounter < 999 {
		fmt.Fprintf(outFile, "\n\n{\\tt\\noindent\\phantom{000}%d:\\ }\n", lineCounter+1)
	} else if lineCounter >= 999 && lineCounter < 9999 {
		fmt.Fprintf(outFile, "\n\n{\\tt\\noindent\\phantom{00}%d:\\ }\n", lineCounter+1)
	} else if lineCounter >= 9999 && lineCounter < 99999 {
		fmt.Fprintf(outFile, "\n\n{\\tt\\noindent\\phantom{0}%d:\\ }\n", lineCounter+1)
	} else if lineCounter >= 99999 {
		fmt.Fprintf(outFile, "\n\n{\\tt %d:\\ }\n", lineCounter+1)
	}
}
