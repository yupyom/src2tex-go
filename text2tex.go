package main

import (
	"fmt"
	"os"
	"strings"
)

func Text2Tex(cptr []string, fptr []*os.File) {

	char_counter := 0
	LineCounter = 0 // use global LineCounter
	page_len := 1
	qt_counter := 0
	dqt_counter := 0
	space_counter := 0
	PREVFlag := 0
	prev_char := 0x20
	BFFlag := 0
	skip_amount := 0
	STATFlag := 0
	NLFlag := 0
	RMFlag := 0
	TTFlag := 0
	DOCFlag := 0
	WARNFlag1 := 0
	WARNFlag2 := 0

	var i int
	var cptr1, cptr2, cptr3, cptr4, cptr5, cptr6, cptr7 int
	var ptr *FlagChar

	for {
		ptr = GetFlagChar(fptr[0])
		if ptr.Character == -1 {
			break
		}
		// DEBUG
		// fmt.Printf("DEBUG: char='%c', flag=%d, idx=%d\n", ptr.Character, ptr.Flag, ptr.BufferIdx)

		if STATFlag == 0 {
			STATFlag++

			TTFlag = ChooseTtFont(ptr.BufferIdx)

			if IsLatexMode {
				fmt.Fprintf(os.Stderr, "src2tex-go: translating ...")
				docStyleSkip := FprintfDocumentstyle(cptr[0], ptr.BufferIdx, fptr[1])
				if docStyleSkip {
					// Skip the rest of the current line (documentstyle line was neutralized)
					for ptr.Character != -1 && ptr.Character != '\n' {
						ptr = GetFlagChar(fptr[0])
					}
				}
			} else {
				fmt.Fprintf(os.Stderr, "src2tex: translating ...")
				FprintfFootline(cptr[0], fptr[1])
				if PageLenMax >= 0 {
					fmt.Fprintf(fptr[1], "\n")
					fmt.Fprintf(fptr[1], "\\baselineskip=0pt\n")
				}
			}

			fmt.Fprintf(fptr[1], "\n")
			InputUserStyle(fptr[1])

			fmt.Fprintf(fptr[1], "\n")
			MergeNttAscii(fptr[1])
			fmt.Fprintf(fptr[1], "\n")
			fmt.Fprintf(fptr[1], "%s\n", TextModeFont)
			fmt.Fprintf(fptr[1], "\n")
			if PageLenMax < 0 {
				fmt.Fprintf(fptr[1], "\\noindent\n")
			} else {
				if IsLatexMode {
					fmt.Fprintf(fptr[1], "\\hfill")
					FprintfLineNumber(fptr[1], int(LineCounter))
				} else {
					fmt.Fprintf(fptr[1], "\\hfill\n\n\\item{\\tt %d:\\ }\n",
						LineCounter+1)
				}
			}
			if IsLatexMode {

				if (char_counter == 0) && (DOCFlag == 0) {
					if BASFlag != 0 {
						if SearchLine(ptr.BufferIdx, "'{\\null}") != 0 || SearchLine(ptr.BufferIdx, "1'{\\null}") != 0 || SearchLine(ptr.BufferIdx, "10'{\\null}") != 0 || SearchLine(ptr.BufferIdx, "100'{\\null}") != 0 || SearchLine(ptr.BufferIdx, "1000'{\\null}") != 0 || SearchLine(ptr.BufferIdx, "rem{\\null}") != 0 || SearchLine(ptr.BufferIdx, "1rem{\\null}") != 0 || SearchLine(ptr.BufferIdx, "10rem{\\null}") != 0 || SearchLine(ptr.BufferIdx, "100rem{\\null}") != 0 || SearchLine(ptr.BufferIdx, "1000rem{\\null}") != 0 {
							NLFlag = 1
						}
					}
					if CFlag != 0 || GoFlag != 0 || JSFlag != 0 || TSFlag != 0 || RustFlag != 0 || KotlinFlag != 0 || SwiftFlag != 0 {
						if SearchLine(ptr.BufferIdx, "/*{\\null}*/") != 0 || SearchLine(ptr.BufferIdx, "//{\\null}") != 0 {
							NLFlag = 1
						}
					}
					if CBLFlag != 0 {
						if SearchLine(ptr.BufferIdx, "*{\\null}") != 0 || SearchLine(ptr.BufferIdx, "/{\\null}") != 0 {
							NLFlag = 1
						}
					}
					if F77Flag != 0 {
						if SearchLine(ptr.BufferIdx, "c{\\null}") != 0 || SearchLine(ptr.BufferIdx, "*{\\null}") != 0 {
							NLFlag = 1
						}
					}
					if LISPFlag != 0 {
						if SearchLine(ptr.BufferIdx, ";{\\null}") != 0 {
							NLFlag = 1
						}
						if SearchLine(ptr.BufferIdx, ";;{\\null}") != 0 {
							NLFlag = 1
						}
						if SearchLine(ptr.BufferIdx, ";;;{\\null}") != 0 {
							NLFlag = 1
						}
						if SearchLine(ptr.BufferIdx, ";;;;{\\null}") != 0 {
							NLFlag = 1
						}
						if SearchLine(ptr.BufferIdx, ";;;;;{\\null}") != 0 {
							NLFlag = 1
						}
					}
					if MAKEFlag != 0 {
						if SearchLine(ptr.BufferIdx, "#{\\null}") != 0 {
							NLFlag = 1
						}
					}
					if PASFlag != 0 {
						if SearchLine(ptr.BufferIdx, "{{\\null}}") != 0 || SearchLine(ptr.BufferIdx, "(*{\\null}*)") != 0 {
							NLFlag = 1
						}
					}
					if PERLFlag != 0 {
						if SearchLine(ptr.BufferIdx, "#{\\null}") != 0 {
							NLFlag = 1
						}
					}
					if SHFlag != 0 || PythonFlag != 0 || RubyFlag != 0 {
						if SearchLine(ptr.BufferIdx, "#{\\null}") != 0 {
							NLFlag = 1
						}
					}
					if TCLFlag != 0 {
						if SearchLine(ptr.BufferIdx, "#{\\null}") != 0 {
							NLFlag = 1
						}
					}
					if MAPFlag != 0 {
						if SearchLine(ptr.BufferIdx, "#{\\null}") != 0 {
							NLFlag = 1
						}
					}
					if MATFlag != 0 {
						if SearchLine(ptr.BufferIdx, "(*{\\null}*)") != 0 {
							NLFlag = 1
						}
					}
					if MLABFlag != 0 {
						if SearchLine(ptr.BufferIdx, "#{\\null}") != 0 || SearchLine(ptr.BufferIdx, "%{\\null}") != 0 {
							NLFlag = 1
						}
					}
					if REDFlag != 0 {
						if SearchLine(ptr.BufferIdx, "%{\\null}") != 0 || SearchLine(ptr.BufferIdx, "COMMENT{\\null};") != 0 || SearchLine(ptr.BufferIdx, "comment{\\null};") != 0 {
							NLFlag = 1
						}
						if NLFlag != 0 {
							DOCFlag++
							for ptr.Character != -1 && ptr.Character != '\n' {
								ptr = GetFlagChar(fptr[0])
							}
							continue
						}
					}
					if char_counter == 0 && LineCounter <= 3 && DOCFlag == 0 {
						if CBLFlag != 0 {
							if SearchLine(ptr.BufferIdx, "*{\\null}") != 0 || SearchLine(ptr.BufferIdx, "/{\\null}") != 0 {
								NLFlag = 1
							}
							if SearchLine(ptr.BufferIdx, "000001*{\\null}") != 0 || SearchLine(ptr.BufferIdx, "000001/{\\null}") != 0 {
								NLFlag = 1
							}
							if SearchLine(ptr.BufferIdx, "000002*{\\null}") != 0 || SearchLine(ptr.BufferIdx, "000002/{\\null}") != 0 {
								NLFlag = 1
							}
							if SearchLine(ptr.BufferIdx, "000003*{\\null}") != 0 || SearchLine(ptr.BufferIdx, "000003/{\\null}") != 0 {
								NLFlag = 1
							}
							if SearchLine(ptr.BufferIdx, "000004*{\\null}") != 0 || SearchLine(ptr.BufferIdx, "000004/{\\null}") != 0 {
								NLFlag = 1
							}
							if SearchLine(ptr.BufferIdx, "000010*{\\null}") != 0 || SearchLine(ptr.BufferIdx, "000010/{\\null}") != 0 {
								NLFlag = 1
							}
							if SearchLine(ptr.BufferIdx, "000020*{\\null}") != 0 || SearchLine(ptr.BufferIdx, "000020/{\\null}") != 0 {
								NLFlag = 1
							}
							if SearchLine(ptr.BufferIdx, "000030*{\\null}") != 0 || SearchLine(ptr.BufferIdx, "000030/{\\null}") != 0 {
								NLFlag = 1
							}
							if SearchLine(ptr.BufferIdx, "000040*{\\null}") != 0 || SearchLine(ptr.BufferIdx, "000040/{\\null}") != 0 {
								NLFlag = 1
							}
							if SearchLine(ptr.BufferIdx, "000100*{\\null}") != 0 || SearchLine(ptr.BufferIdx, "000100/{\\null}") != 0 {
								NLFlag = 1
							}
							if SearchLine(ptr.BufferIdx, "000200*{\\null}") != 0 || SearchLine(ptr.BufferIdx, "000200/{\\null}") != 0 {
								NLFlag = 1
							}
							if SearchLine(ptr.BufferIdx, "000300*{\\null}") != 0 || SearchLine(ptr.BufferIdx, "000300/{\\null}") != 0 {
								NLFlag = 1
							}
							if SearchLine(ptr.BufferIdx, "000400*{\\null}") != 0 || SearchLine(ptr.BufferIdx, "000400/{\\null}") != 0 {
								NLFlag = 1
							}
						}
						if NLFlag != 0 {
							DOCFlag++
							for ptr.Character != -1 && ptr.Character != '\n' {
								ptr = GetFlagChar(fptr[0])
							}
							continue
						}
					}
				}
			}
		}

		if (ptr.Flag == 1) && (CBLFlag != 0) && (char_counter == 0) {
						cptr1 = ptr.BufferIdx
						cptr2 = IncBufPtr(cptr1)
						cptr3 = IncBufPtr(cptr2)
						cptr4 = IncBufPtr(cptr3)
						cptr5 = IncBufPtr(cptr4)
						cptr6 = IncBufPtr(cptr5)
						cptr7 = IncBufPtr(cptr6)
						if Buffer[cptr7] == '/' {
							fmt.Fprintf(fptr[1], "\\vfill\\eject\n\n\\noindent\n")
						}
					}
					if (PREVFlag > 1) && (ptr.Flag == 0) {
						fmt.Fprintf(os.Stderr,
							"\nError: unexpected end of TeX-mode in %s\n", cptr[0])
						fmt.Fprintf(os.Stderr,
							"       illegal transition TeX-mode -> Text-mode\n")
						os.Exit(1)
					}

					if (PREVFlag == 0) && (ptr.Flag > 0) {

						fmt.Fprintf(os.Stderr, ".")

						if ((CFlag == 0) && (F77Flag == 0) && (MAKEFlag == 0) && (PASFlag == 0) && (PERLFlag == 0) && (SHFlag == 0) && (TCLFlag == 0) && (MAPFlag == 0) && (MATFlag == 0) && (MLABFlag == 0) && (PythonFlag == 0) && (RubyFlag == 0) && (GoFlag == 0) && (JSFlag == 0) && (TSFlag == 0) && (RustFlag == 0) && (KotlinFlag == 0) && (SwiftFlag == 0)) || (char_counter >= HtabSize) {
							// Both full-line and inline comments use TeXModeFont (roman/mincho)
							fmt.Fprintf(fptr[1], "%s", TeXModeFont)
						}

						if (CFlag != 0 || GoFlag != 0 || JSFlag != 0 || TSFlag != 0 || RustFlag != 0 || KotlinFlag != 0 || SwiftFlag != 0) && (char_counter < HtabSize) {
							TTFlag = 0
							cptr1 = ptr.BufferIdx
							cptr2 = IncBufPtr(cptr1)
							cptr3 = IncBufPtr(cptr2)
							cptr4 = IncBufPtr(cptr3)
							cptr5 = IncBufPtr(cptr4)
							cptr6 = IncBufPtr(cptr5)
							for i = 0; i < 1024; i++ {
								if (Buffer[cptr1] == '{') && (Buffer[cptr2] == '\\') {
									RMFlag = 1
									TTFlag = 0
									break
								}
								if (IsLatexMode && Buffer[cptr1] == '\\' && Buffer[cptr2] == '[') || (!IsLatexMode && Buffer[cptr1] == '$' && Buffer[cptr2] == '$') {
									RMFlag = 1
									TTFlag = 0
									break
								}
								if (IsLatexMode && Buffer[cptr1] == '\\' && Buffer[cptr2] == '(') || (!IsLatexMode && Buffer[cptr1] != '\\' && Buffer[cptr2] == '$') {
									RMFlag = 1
									TTFlag = 0
									break
								}
								if ((Buffer[cptr1] == '*') && (Buffer[cptr2] == '/')) && ((Buffer[cptr3] == '\n') || (Buffer[cptr3] == '\r')) && ((Buffer[cptr4] == '/') && (Buffer[cptr5] == '*')) {
									RMFlag = 0
									TTFlag = 2
									break
								}
								if ((Buffer[cptr1] == '*') && (Buffer[cptr2] == '/')) && ((Buffer[cptr3] == '\n') || (Buffer[cptr3] == '\r')) && ((Buffer[cptr4] == '\n') || (Buffer[cptr4] == '\r')) && ((Buffer[cptr5] == '/') && (Buffer[cptr6] == '*')) {
									RMFlag = 0
									TTFlag = 2
									break
								}
								if ((Buffer[cptr1] == '*') && (Buffer[cptr2] == '/')) && ((Buffer[cptr3] == '\n') || (Buffer[cptr3] == '\r')) {
									RMFlag = 0
									TTFlag = 1
									break
								}
								if ((Buffer[cptr1] == '\n') || (Buffer[cptr1] == '\r')) && (Buffer[cptr2] == '/') && (Buffer[cptr3] == '*') {
									RMFlag = 0
									TTFlag = 1
									break
								}
								if (Buffer[cptr1] == '\n') || (Buffer[cptr1] == '\r') {
									TTFlag = 1
								}
								cptr1 = IncBufPtr(cptr1)
								cptr2 = IncBufPtr(cptr2)
								cptr3 = IncBufPtr(cptr3)
								cptr4 = IncBufPtr(cptr4)
								cptr5 = IncBufPtr(cptr5)
								cptr6 = IncBufPtr(cptr6)
							}

							if (TTFlag > 1) && (WARNFlag1 <= 5) {
								WARNFlag1++

								if WARNFlag1 > 5 {
									if WARNFlag2 <= 2 {
										fmt.Fprintf(os.Stderr, "\n")
									}
									fmt.Fprintf(os.Stderr,
										"Warning: It is better to use TeX-mode\n")
									fmt.Fprintf(os.Stderr,
										"         when you write long comment in C.\n")
								}
							}
							// Full-line comment: always use TeXModeFont (comment/mincho font)
							fmt.Fprintf(fptr[1], "%s", TeXModeFont)
						}

						if (F77Flag != 0) && (char_counter < HtabSize) {
							TTFlag = 0
							cptr1 = ptr.BufferIdx
							cptr2 = IncBufPtr(cptr1)
							cptr3 = IncBufPtr(cptr2)
							for i = 0; i < 1024; i++ {
								if (Buffer[cptr1] == '{') && (Buffer[cptr2] == '\\') {
									RMFlag = 1
									TTFlag = 0
									break
								}
								if (IsLatexMode && Buffer[cptr1] == '\\' && Buffer[cptr2] == '[') || (!IsLatexMode && Buffer[cptr1] == '$' && Buffer[cptr2] == '$') {
									RMFlag = 1
									TTFlag = 0
									break
								}
								if (IsLatexMode && Buffer[cptr1] == '\\' && Buffer[cptr2] == '(') || (!IsLatexMode && Buffer[cptr1] != '\\' && Buffer[cptr2] == '$') {
									RMFlag = 1
									TTFlag = 0
									break
								}
								if ((Buffer[cptr1] != '\n') && (Buffer[cptr1] != '\r')) && ((Buffer[cptr2] == '*') || (Buffer[cptr2] == 'C') || (Buffer[cptr2] == 'c')) && ((Buffer[cptr3] == '\n') || (Buffer[cptr3] == '\r')) {
									RMFlag = 0
									TTFlag = 2
									break
								}
								if (Buffer[cptr1] == '\n') || (Buffer[cptr1] == '\r') {
									TTFlag = 1
								}
								cptr1 = IncBufPtr(cptr1)
								cptr2 = IncBufPtr(cptr2)
								cptr3 = IncBufPtr(cptr3)
							}

							if (TTFlag > 1) && (WARNFlag1 <= 5) {
								WARNFlag1++

								if WARNFlag1 > 5 {
									if WARNFlag2 <= 2 {
										fmt.Fprintf(os.Stderr, "\n")
									}
									fmt.Fprintf(os.Stderr,
										"Warning: It is better to use TeX-mode\n")
									fmt.Fprintf(os.Stderr,
										"         when you write long comment in FORTRAN.\n")
								}
							}
							// Full-line comment: always use TeXModeFont (comment/mincho font)
							fmt.Fprintf(fptr[1], "%s", TeXModeFont)
						}

						if (MAKEFlag != 0) && (char_counter < HtabSize) {
							TTFlag = 0
							cptr1 = ptr.BufferIdx
							cptr2 = IncBufPtr(cptr1)
							cptr3 = IncBufPtr(cptr2)
							for i = 0; i < 1024; i++ {
								if (Buffer[cptr1] == '{') && (Buffer[cptr2] == '\\') {
									RMFlag = 1
									TTFlag = 0
									break
								}
								if (IsLatexMode && Buffer[cptr1] == '\\' && Buffer[cptr2] == '[') || (!IsLatexMode && Buffer[cptr1] == '$' && Buffer[cptr2] == '$') {
									RMFlag = 1
									TTFlag = 0
									break
								}
								if (IsLatexMode && ((Buffer[cptr1] == '\\') && (Buffer[cptr2] == '('))) || (!IsLatexMode && (((Buffer[cptr1] == '\t') || (Buffer[cptr1] == ' ')) && (Buffer[cptr2] == '$') && (Buffer[cptr3] != '('))) {
									RMFlag = 1
									TTFlag = 0
									break
								}
								if ((Buffer[cptr1] != '\n') && (Buffer[cptr1] != '\r')) && (Buffer[cptr2] == '#') && ((Buffer[cptr3] == '\n') || (Buffer[cptr3] == '\r')) {
									RMFlag = 0
									TTFlag = 2
									break
								}
								if (Buffer[cptr1] == '\n') || (Buffer[cptr1] == '\r') {
									TTFlag = 1
								}
								cptr1 = IncBufPtr(cptr1)
								cptr2 = IncBufPtr(cptr2)
								cptr3 = IncBufPtr(cptr3)
							}

							if (TTFlag > 1) && (WARNFlag1 <= 5) {
								WARNFlag1++

								if WARNFlag1 > 5 {
									if WARNFlag2 <= 2 {
										fmt.Fprintf(os.Stderr, "\n")
									}
									fmt.Fprintf(os.Stderr,
										"Warning: It is better to use TeX-mode\n")
									fmt.Fprintf(os.Stderr,
										"         when you write long comment in MAKE.\n")
								}
							}
							// Full-line comment: always use TeXModeFont (comment/mincho font)
							fmt.Fprintf(fptr[1], "%s", TeXModeFont)
						}

						if (PASFlag != 0) && (char_counter < HtabSize) {
							TTFlag = 0
							cptr1 = ptr.BufferIdx
							cptr2 = IncBufPtr(cptr1)
							cptr3 = IncBufPtr(cptr2)
							cptr4 = IncBufPtr(cptr3)
							cptr5 = IncBufPtr(cptr4)
							cptr6 = IncBufPtr(cptr5)
							for i = 0; i < 1024; i++ {
								if (Buffer[cptr1] == '{') && (Buffer[cptr2] == '\\') {
									RMFlag = 1
									TTFlag = 0
									break
								}
								if (IsLatexMode && Buffer[cptr1] == '\\' && Buffer[cptr2] == '[') || (!IsLatexMode && Buffer[cptr1] == '$' && Buffer[cptr2] == '$') {
									RMFlag = 1
									TTFlag = 0
									break
								}
								if (IsLatexMode && Buffer[cptr1] == '\\' && Buffer[cptr2] == '(') || (!IsLatexMode && Buffer[cptr1] != '\\' && Buffer[cptr2] == '$') {
									RMFlag = 1
									TTFlag = 0
									break
								}
								if (Buffer[cptr1] == '}') && ((Buffer[cptr2] == '\n') || (Buffer[cptr2] == '\r')) && (Buffer[cptr3] == '{') {
									RMFlag = 0
									TTFlag = 2
									break
								}
								if ((Buffer[cptr1] == '*') && (Buffer[cptr2] == ')')) && ((Buffer[cptr3] == '\n') || (Buffer[cptr3] == '\r')) && ((Buffer[cptr4] == '(') && (Buffer[cptr5] == '*')) {
									RMFlag = 0
									TTFlag = 2
									break
								}
								if (Buffer[cptr1] == '}') && ((Buffer[cptr2] == '\n') || (Buffer[cptr2] == '\r')) && ((Buffer[cptr3] == '\n') || (Buffer[cptr3] == '\r')) && (Buffer[cptr4] == '{') {
									RMFlag = 0
									TTFlag = 2
									break
								}
								if ((Buffer[cptr1] == '*') && (Buffer[cptr2] == ')')) && ((Buffer[cptr3] == '\n') || (Buffer[cptr3] == '\r')) && ((Buffer[cptr4] == '\n') || (Buffer[cptr4] == '\r')) && ((Buffer[cptr5] == '(') && (Buffer[cptr6] == '*')) {
									RMFlag = 0
									TTFlag = 2
									break
								}
								if (Buffer[cptr1] == '}') && ((Buffer[cptr2] == '\n') || (Buffer[cptr2] == '\r')) {
									RMFlag = 0
									TTFlag = 1
									break
								}
								if ((Buffer[cptr1] == '*') && (Buffer[cptr2] == ')')) && ((Buffer[cptr3] == '\n') || (Buffer[cptr3] == '\r')) {
									RMFlag = 0
									TTFlag = 1
									break
								}
								if ((Buffer[cptr1] == '\n') || (Buffer[cptr1] == '\r')) && (Buffer[cptr2] == '}') {
									RMFlag = 0
									TTFlag = 1
									break
								}
								if ((Buffer[cptr1] == '\n') || (Buffer[cptr1] == '\r')) && (Buffer[cptr2] == '(') && (Buffer[cptr3] == '*') {
									RMFlag = 0
									TTFlag = 1
									break
								}
								if (Buffer[cptr1] == '\n') || (Buffer[cptr1] == '\r') {
									TTFlag = 1
								}
								cptr1 = IncBufPtr(cptr1)
								cptr2 = IncBufPtr(cptr2)
								cptr3 = IncBufPtr(cptr3)
								cptr4 = IncBufPtr(cptr4)
								cptr5 = IncBufPtr(cptr5)
								cptr6 = IncBufPtr(cptr6)
							}

							if (TTFlag > 1) && (WARNFlag1 <= 5) {
								WARNFlag1++

								if WARNFlag1 > 5 {
									if WARNFlag2 <= 2 {
										fmt.Fprintf(os.Stderr, "\n")
									}
									fmt.Fprintf(os.Stderr,
										"Warning: It is better to use TeX-mode\n")
									fmt.Fprintf(os.Stderr,
										"         when you write long comment in PASCAL.\n")
								}
							}
							// Full-line comment: always use TeXModeFont (comment/mincho font)
							fmt.Fprintf(fptr[1], "%s", TeXModeFont)
						}

						if (PERLFlag != 0) && (char_counter < HtabSize) {
							TTFlag = 0
							cptr1 = ptr.BufferIdx
							cptr2 = IncBufPtr(cptr1)
							cptr3 = IncBufPtr(cptr2)
							for i = 0; i < 1024; i++ {
								if (Buffer[cptr1] == '{') && (Buffer[cptr2] == '\\') {
									RMFlag = 1
									TTFlag = 0
									break
								}
								if (IsLatexMode && Buffer[cptr1] == '\\' && Buffer[cptr2] == '[') || (!IsLatexMode && Buffer[cptr1] == '$' && Buffer[cptr2] == '$') {
									RMFlag = 1
									TTFlag = 0
									break
								}
								if (IsLatexMode && ((Buffer[cptr1] == '\\') && (Buffer[cptr2] == '('))) || (!IsLatexMode && (((Buffer[cptr1] == '\t') || (Buffer[cptr1] == ' ')) && (Buffer[cptr2] == '$') && (Buffer[cptr3] == '\\'))) {
									RMFlag = 1
									TTFlag = 0
									break
								}
								if ((Buffer[cptr1] != '\n') && (Buffer[cptr1] != '\r')) && (Buffer[cptr2] == '#') && ((Buffer[cptr3] == '\n') || (Buffer[cptr3] == '\r')) {
									RMFlag = 0
									TTFlag = 2
									break
								}
								if (Buffer[cptr1] == '\n') || (Buffer[cptr1] == '\r') {
									TTFlag = 1
								}
								cptr1 = IncBufPtr(cptr1)
								cptr2 = IncBufPtr(cptr2)
								cptr3 = IncBufPtr(cptr3)
							}

							if (TTFlag > 1) && (WARNFlag1 <= 5) {
								WARNFlag1++

								if WARNFlag1 > 5 {
									if WARNFlag2 <= 2 {
										fmt.Fprintf(os.Stderr, "\n")
									}
									fmt.Fprintf(os.Stderr,
										"Warning: It is better to use TeX-mode\n")
									fmt.Fprintf(os.Stderr,
										"         when you write long comment in PERL.\n")
								}
							}
							// Full-line comment: always use TeXModeFont (comment/mincho font)
							fmt.Fprintf(fptr[1], "%s", TeXModeFont)
						}

						if (SHFlag != 0 || PythonFlag != 0 || RubyFlag != 0) && (char_counter < HtabSize) {
							TTFlag = 0
							cptr1 = ptr.BufferIdx
							cptr2 = IncBufPtr(cptr1)
							cptr3 = IncBufPtr(cptr2)
							for i = 0; i < 1024; i++ {
								if (Buffer[cptr1] == '{') && (Buffer[cptr2] == '\\') {
									RMFlag = 1
									TTFlag = 0
									break
								}
								if (IsLatexMode && Buffer[cptr1] == '\\' && Buffer[cptr2] == '[') || (!IsLatexMode && Buffer[cptr1] == '$' && Buffer[cptr2] == '$') {
									RMFlag = 1
									TTFlag = 0
									break
								}
								if (IsLatexMode && ((Buffer[cptr1] == '\\') && (Buffer[cptr2] == '('))) || (!IsLatexMode && (((Buffer[cptr1] == '\t') || (Buffer[cptr1] == ' ')) && (Buffer[cptr2] == '$'))) {
									RMFlag = 1
									TTFlag = 0
									break
								}
								if ((Buffer[cptr1] != '\n') && (Buffer[cptr1] != '\r')) && (Buffer[cptr2] == '#') && ((Buffer[cptr3] == '\n') || (Buffer[cptr3] == '\r')) {
									RMFlag = 0
									TTFlag = 2
									break
								}
								if (Buffer[cptr1] == '\n') || (Buffer[cptr1] == '\r') {
									TTFlag = 1
								}
								cptr1 = IncBufPtr(cptr1)
								cptr2 = IncBufPtr(cptr2)
								cptr3 = IncBufPtr(cptr3)
							}

							if (TTFlag > 1) && (WARNFlag1 <= 5) {
								WARNFlag1++

								if WARNFlag1 > 5 {
									if WARNFlag2 <= 2 {
										fmt.Fprintf(os.Stderr, "\n")
									}
									fmt.Fprintf(os.Stderr,
										"Warning: It is better to use TeX-mode\n")
									fmt.Fprintf(os.Stderr,
										"         when you write long comment in SHELL.\n")
								}
							}
							// Full-line comment: always use TeXModeFont (comment/mincho font)
							fmt.Fprintf(fptr[1], "%s", TeXModeFont)
						}

						if (TCLFlag != 0) && (char_counter < HtabSize) {
							TTFlag = 0
							cptr1 = ptr.BufferIdx
							cptr2 = IncBufPtr(cptr1)
							cptr3 = IncBufPtr(cptr2)
							for i = 0; i < 1024; i++ {
								if (Buffer[cptr1] == '{') && (Buffer[cptr2] == '\\') {
									RMFlag = 1
									TTFlag = 0
									break
								}
								if (IsLatexMode && Buffer[cptr1] == '\\' && Buffer[cptr2] == '[') || (!IsLatexMode && Buffer[cptr1] == '$' && Buffer[cptr2] == '$') {
									RMFlag = 1
									TTFlag = 0
									break
								}
								if (IsLatexMode && ((Buffer[cptr1] == '\\') && (Buffer[cptr2] == '('))) || (!IsLatexMode && (((Buffer[cptr1] == '\t') || (Buffer[cptr1] == ' ')) && (Buffer[cptr2] == '$'))) {
									RMFlag = 1
									TTFlag = 0
									break
								}
								if ((Buffer[cptr1] != '\n') && (Buffer[cptr1] != '\r')) && (Buffer[cptr2] == '#') && ((Buffer[cptr3] == '\n') || (Buffer[cptr3] == '\r')) {
									RMFlag = 0
									TTFlag = 2
									break
								}
								if (Buffer[cptr1] == '\n') || (Buffer[cptr1] == '\r') {
									TTFlag = 1
								}
								cptr1 = IncBufPtr(cptr1)
								cptr2 = IncBufPtr(cptr2)
								cptr3 = IncBufPtr(cptr3)
							}

							if (TTFlag > 1) && (WARNFlag1 <= 5) {
								WARNFlag1++

								if WARNFlag1 > 5 {
									if WARNFlag2 <= 2 {
										fmt.Fprintf(os.Stderr, "\n")
									}
									fmt.Fprintf(os.Stderr,
										"Warning: It is better to use TeX-mode\n")
									fmt.Fprintf(os.Stderr,
										"         when you write long comment in TCL/TK.\n")
								}
							}
							// Full-line comment: always use TeXModeFont (comment/mincho font)
							fmt.Fprintf(fptr[1], "%s", TeXModeFont)
						}

						if (MAPFlag != 0) && (char_counter < HtabSize) {
							TTFlag = 0
							cptr1 = ptr.BufferIdx
							cptr2 = IncBufPtr(cptr1)
							cptr3 = IncBufPtr(cptr2)
							for i = 0; i < 1024; i++ {
								if (Buffer[cptr1] == '{') && (Buffer[cptr2] == '\\') {
									RMFlag = 1
									TTFlag = 0
									break
								}
								if (IsLatexMode && Buffer[cptr1] == '\\' && Buffer[cptr2] == '[') || (!IsLatexMode && Buffer[cptr1] == '$' && Buffer[cptr2] == '$') {
									RMFlag = 1
									TTFlag = 0
									break
								}
								if (IsLatexMode && Buffer[cptr1] == '\\' && Buffer[cptr2] == '(') || (!IsLatexMode && Buffer[cptr1] != '\\' && Buffer[cptr2] == '$') {
									RMFlag = 1
									TTFlag = 0
									break
								}
								if ((Buffer[cptr1] != '\n') && (Buffer[cptr1] != '\r')) && (Buffer[cptr2] == '#') && ((Buffer[cptr3] == '\n') || (Buffer[cptr3] == '\r')) {
									RMFlag = 0
									TTFlag = 2
									break
								}
								if (Buffer[cptr1] == '\n') || (Buffer[cptr1] == '\r') {
									TTFlag = 1
								}
								cptr1 = IncBufPtr(cptr1)
								cptr2 = IncBufPtr(cptr2)
								cptr3 = IncBufPtr(cptr3)
							}

							if (TTFlag > 1) && (WARNFlag1 <= 5) {
								WARNFlag1++

								if WARNFlag1 > 5 {
									if WARNFlag2 <= 2 {
										fmt.Fprintf(os.Stderr, "\n")
									}
									fmt.Fprintf(os.Stderr,
										"Warning: It is better to use TeX-mode\n")
									fmt.Fprintf(os.Stderr,
										"         when you write long comment in MAPLE.\n")
								}
							}
							// Full-line comment: always use TeXModeFont (comment/mincho font)
							fmt.Fprintf(fptr[1], "%s", TeXModeFont)
						}

						if (MATFlag != 0) && (char_counter < HtabSize) {
							TTFlag = 0
							cptr1 = ptr.BufferIdx
							cptr2 = IncBufPtr(cptr1)
							cptr3 = IncBufPtr(cptr2)
							cptr4 = IncBufPtr(cptr3)
							cptr5 = IncBufPtr(cptr4)
							cptr6 = IncBufPtr(cptr5)
							for i = 0; i < 1024; i++ {
								if (Buffer[cptr1] == '{') && (Buffer[cptr2] == '\\') {
									RMFlag = 1
									TTFlag = 0
									break
								}
								if (IsLatexMode && Buffer[cptr1] == '\\' && Buffer[cptr2] == '[') || (!IsLatexMode && Buffer[cptr1] == '$' && Buffer[cptr2] == '$') {
									RMFlag = 1
									TTFlag = 0
									break
								}
								if (IsLatexMode && Buffer[cptr1] == '\\' && Buffer[cptr2] == '(') || (!IsLatexMode && Buffer[cptr1] != '\\' && Buffer[cptr2] == '$') {
									RMFlag = 1
									TTFlag = 0
									break
								}
								if ((Buffer[cptr1] == '*') && (Buffer[cptr2] == ')')) && ((Buffer[cptr3] == '\n') || (Buffer[cptr3] == '\r')) && ((Buffer[cptr4] == '(') && (Buffer[cptr5] == '*')) {
									RMFlag = 0
									TTFlag = 2
									break
								}
								if ((Buffer[cptr1] == '*') && (Buffer[cptr2] == ')')) && ((Buffer[cptr3] == '\n') || (Buffer[cptr3] == '\r')) && ((Buffer[cptr4] == '\n') || (Buffer[cptr4] == '\r')) && ((Buffer[cptr5] == '(') && (Buffer[cptr6] == '*')) {
									RMFlag = 0
									TTFlag = 2
									break
								}
								if ((Buffer[cptr1] == '*') && (Buffer[cptr2] == ')')) && ((Buffer[cptr3] == '\n') || (Buffer[cptr3] == '\r')) {
									RMFlag = 0
									TTFlag = 1
									break
								}
								if ((Buffer[cptr1] == '\n') || (Buffer[cptr1] == '\r')) && (Buffer[cptr2] == '(') && (Buffer[cptr3] == '*') {
									RMFlag = 0
									TTFlag = 1
									break
								}
								if (Buffer[cptr1] == '\n') || (Buffer[cptr1] == '\r') {
									TTFlag = 1
								}
								cptr1 = IncBufPtr(cptr1)
								cptr2 = IncBufPtr(cptr2)
								cptr3 = IncBufPtr(cptr3)
								cptr4 = IncBufPtr(cptr4)
								cptr5 = IncBufPtr(cptr5)
								cptr6 = IncBufPtr(cptr6)
							}

							if (TTFlag > 1) && (WARNFlag1 <= 5) {
								WARNFlag1++

								if WARNFlag1 > 5 {
									if WARNFlag2 <= 2 {
										fmt.Fprintf(os.Stderr, "\n")
									}
									fmt.Fprintf(os.Stderr,
										"Warning: It is better to use TeX-mode\n")
									fmt.Fprintf(os.Stderr,
										"         when you write long comment in")
									fmt.Fprintf(os.Stderr,
										" MATHEMATICA.\n")
								}
							}
							// Full-line comment: always use TeXModeFont (comment/mincho font)
							fmt.Fprintf(fptr[1], "%s", TeXModeFont)
						}

						if (MLABFlag != 0) && (char_counter < HtabSize) {
							TTFlag = 0
							cptr1 = ptr.BufferIdx
							cptr2 = IncBufPtr(cptr1)
							cptr3 = IncBufPtr(cptr2)
							for i = 0; i < 1024; i++ {
								if (Buffer[cptr1] == '{') && (Buffer[cptr2] == '\\') {
									RMFlag = 1
									TTFlag = 0
									break
								}
								if (IsLatexMode && Buffer[cptr1] == '\\' && Buffer[cptr2] == '[') || (!IsLatexMode && Buffer[cptr1] == '$' && Buffer[cptr2] == '$') {
									RMFlag = 1
									TTFlag = 0
									break
								}
								if (IsLatexMode && Buffer[cptr1] == '\\' && Buffer[cptr2] == '(') || (!IsLatexMode && Buffer[cptr1] != '\\' && Buffer[cptr2] == '$') {
									RMFlag = 1
									TTFlag = 0
									break
								}
								if (Buffer[cptr1] == '\n') && ((Buffer[cptr2] == '#') || (Buffer[cptr2] == '%')) && ((Buffer[cptr3] == '\t') || (Buffer[cptr3] == ' ')) {
									RMFlag = 0
									TTFlag = 2
									break
								}
								if ((Buffer[cptr1] != '\n') && (Buffer[cptr1] != '\r')) && ((Buffer[cptr2] == '#') || (Buffer[cptr2] == '%')) && ((Buffer[cptr3] == '\n') || (Buffer[cptr3] == '\r')) {
									RMFlag = 0
									TTFlag = 2
									break
								}
								if (Buffer[cptr1] == '\n') || (Buffer[cptr1] == '\r') {
									TTFlag = 1
								}
								cptr1 = IncBufPtr(cptr1)
								cptr2 = IncBufPtr(cptr2)
								cptr3 = IncBufPtr(cptr3)
							}

							if (TTFlag > 1) && (WARNFlag1 <= 5) {
								WARNFlag1++

								if WARNFlag1 > 5 {
									if WARNFlag2 <= 2 {
										fmt.Fprintf(os.Stderr, "\n")
									}
									fmt.Fprintf(os.Stderr,
										"Warning: It is better to use TeX-mode\n")
									fmt.Fprintf(os.Stderr,
										"         when you write long comment in")
									fmt.Fprintf(os.Stderr,
										" MATLAB.\n")
								}
							}
							// Full-line comment: always use TeXModeFont (comment/mincho font)
							fmt.Fprintf(fptr[1], "%s", TeXModeFont)
						}
					}
					_ = RMFlag // RMFlag is assigned in buffer scanning loops but no longer read for font selection

					if (PREVFlag > 0) && (ptr.Flag == 0) {
						if char_counter == 0 {
							fmt.Fprintf(fptr[1], "%s", TextModeFont)
						} else {
							fmt.Fprintf(fptr[1], "\n%s", TextModeFont)
						}
					}
					if (PREVFlag > 1) && (ptr.Flag == 1) {
						// TeX-mode to comment-mode transition: always use TeXModeFont (roman/mincho)
						fmt.Fprintf(fptr[1], "%s", TeXModeFont)
					}

					if ptr.Flag == 0 {
						if (prev_char < 'A') || ((prev_char > 'Z') && (prev_char < 'a')) || (prev_char > 'z') {
							if ((PASFlag != 0) || (REDFlag != 0)) && (qt_counter == 0) && (BFFlag == 0) {
								BFFlag = GetBfFlag(ptr.BufferIdx)
							}
						}
					}

					if (ptr.Flag >= 1) && (char_counter == 0) {
						if parse_options(ptr) != 0 {
							// After options parsed at line start (full-line comment),
							// use TeXModeFont to honor texfont= directives
							fmt.Fprintf(fptr[1], "%s", TeXModeFont)
						}
					}

					if (ptr.Flag >= 1) && (char_counter == 0) {
						cptr1 = ptr.BufferIdx
						for Buffer[cptr1] != '\\' && Buffer[cptr1] != '\n' && Buffer[cptr1] != -1 {
							cptr1 = IncBufPtr(cptr1)
						}
						if str_cmp(cptr1, "\\src2tex{") == 0 {
							fmt.Fprintf(fptr[1], "%c ", 0x25)

						}
					}
					PREVFlag = ptr.Flag
					prev_char = ptr.Character

					if ptr.Flag <= 1 {
						switch ptr.Character {
						case 0x00:
							fmt.Fprintf(fptr[1],
								"{\\sevenrm N\\kern-.15em\\lower.5ex\\hbox{U}}")
							break
						case 0x01:
							fmt.Fprintf(fptr[1],
								"{\\sevenrm S\\kern-.15em\\lower.5ex\\hbox{H}}")
							break
						case 0x02:
							fmt.Fprintf(fptr[1],
								"{\\sevenrm S\\kern-.15em\\lower.5ex\\hbox{X}}")
							break
						case 0x03:
							fmt.Fprintf(fptr[1],
								"{\\sevenrm E\\kern-.15em\\lower.5ex\\hbox{X}}")
							break
						case 0x04:
							fmt.Fprintf(fptr[1],
								"{\\sevenrm E\\kern-.15em\\lower.5ex\\hbox{T}}")
							break
						case 0x05:
							fmt.Fprintf(fptr[1],
								"{\\sevenrm E\\kern-.15em\\lower.5ex\\hbox{Q}}")
							break
						case 0x06:
							fmt.Fprintf(fptr[1],
								"{\\sevenrm A\\kern-.15em\\lower.5ex\\hbox{K}}")
							break
						case 0x07:
							fmt.Fprintf(fptr[1],
								"{\\sevenrm B\\kern-.15em\\lower.5ex\\hbox{L}}")
							break
						case '\b':
							fmt.Fprintf(fptr[1],
								"{\\sevenrm B\\kern-.15em\\lower.5ex\\hbox{S}}")
							break
						case '\t':
							skip_amount = HtabSize - (char_counter % HtabSize)
							if UseUnifiedFont {
								fmt.Fprintf(fptr[1], "{%s%s}", TextModeFont, strings.Repeat("\\ ", skip_amount))
							} else {
								fmt.Fprintf(fptr[1], "{%s\\kern%d\\charwd}",
									TextModeFont, skip_amount)
							}
							char_counter += skip_amount - 1
							break
						case '\n':
							if (PageLenMax > 0) && (page_len >= PageLenMax) {
								fmt.Fprintf(fptr[1], "\n\n\\vfill\\eject\n\n")
								page_len = 0
							}
							LineCounter++
							page_len++
							if char_counter == 0 {
								cptr1 = IncBufPtr(ptr.BufferIdx)
								if (PageLenMax < 0) || (Buffer[cptr1] == -1) {
									fmt.Fprintf(fptr[1], "\\hfill\n\n\\noindent\n")
								} else {
									if IsLatexMode {
										fmt.Fprintf(fptr[1], "\\hfill")
										FprintfLineNumber(fptr[1], int(LineCounter))
									} else {
										fmt.Fprintf(fptr[1], "\\hfill\n\n\\item{\\tt %d:\\ }\n",
											LineCounter+1)
									}
								}
							} else {
								cptr1 = IncBufPtr(ptr.BufferIdx)
								if (PageLenMax < 0) || (Buffer[cptr1] == -1) {
									fmt.Fprintf(fptr[1], "\n\n\\noindent\n")
								} else {
									if IsLatexMode {
										FprintfLineNumber(fptr[1], int(LineCounter))
									} else {
										fmt.Fprintf(fptr[1], "\n\n\\item{\\tt %d:\\ }\n",
											LineCounter+1)
									}
								}
							}
							if ptr.Flag <= 1 {
								fmt.Fprintf(fptr[1], "\\mbox{}")
							}
							char_counter = -1
							break
						case '\v':
							fmt.Fprintf(fptr[1], "{\\vskip%dex\\relax }", VtabSize)
							break
						case '\f':

							fmt.Fprintf(fptr[1],
								"\\vfill\\eject\n\n\\noindent\n")
							break
						case '\r':
							if true {
								fmt.Fprintf(fptr[1],
									"{\\sevenrm C\\kern-.15em\\lower.5ex\\hbox{R}}")
							} else {
								fptr[1].Write([]byte{byte(ptr.Character)})
							}
							break
						case 0x0e:
							fmt.Fprintf(fptr[1],
								"{\\sevenrm S\\kern-.15em\\lower.5ex\\hbox{O}}")
							break
						case 0x0f:
							fmt.Fprintf(fptr[1],
								"{\\sevenrm S\\kern-.15em\\lower.5ex\\hbox{I}}")
							break
						case 0x10:
							fmt.Fprintf(fptr[1],
								"{\\sevenrm D\\kern-.15em\\lower.5ex\\hbox{L}}")
							break
						case 0x11:
							fmt.Fprintf(fptr[1],
								"{\\sevenrm D\\kern-.15em\\lower.5ex\\hbox{1}}")
							break
						case 0x12:
							fmt.Fprintf(fptr[1],
								"{\\sevenrm D\\kern-.15em\\lower.5ex\\hbox{2}}")
							break
						case 0x13:
							fmt.Fprintf(fptr[1],
								"{\\sevenrm D\\kern-.15em\\lower.5ex\\hbox{3}}")
							break
						case 0x14:
							fmt.Fprintf(fptr[1],
								"{\\sevenrm D\\kern-.15em\\lower.5ex\\hbox{4}}")
							break
						case 0x15:
							fmt.Fprintf(fptr[1],
								"{\\sevenrm N\\kern-.15em\\lower.5ex\\hbox{K}}")
							break
						case 0x16:
							fmt.Fprintf(fptr[1],
								"{\\sevenrm S\\kern-.15em\\lower.5ex\\hbox{Y}}")
							break
						case 0x17:
							fmt.Fprintf(fptr[1],
								"{\\sevenrm E\\kern-.15em\\lower.5ex\\hbox{B}}")
							break
						case 0x18:
							fmt.Fprintf(fptr[1],
								"{\\sevenrm C\\kern-.15em\\lower.5ex\\hbox{N}}")
							break
						case 0x19:
							fmt.Fprintf(fptr[1],
								"{\\sevenrm E\\kern-.15em\\lower.5ex\\hbox{M}}")
							break
						case 0x1a:
							fmt.Fprintf(fptr[1],
								"{\\sevenrm S\\kern-.15em\\lower.5ex\\hbox{B}}")
							break
						case 0x1b:
							fmt.Fprintf(fptr[1],
								"{\\sevenrm E\\kern-.15em\\lower.5ex\\hbox{C}}")
							break
						case 0x1c:
							fmt.Fprintf(fptr[1],
								"{\\sevenrm F\\kern-.15em\\lower.5ex\\hbox{S}}")
							break
						case 0x1d:
							fmt.Fprintf(fptr[1],
								"{\\sevenrm G\\kern-.15em\\lower.5ex\\hbox{S}}")
							break
						case 0x1e:
							fmt.Fprintf(fptr[1],
								"{\\sevenrm R\\kern-.15em\\lower.5ex\\hbox{S}}")
							break
						case 0x1f:
							fmt.Fprintf(fptr[1],
								"{\\sevenrm U\\kern-.15em\\lower.5ex\\hbox{S}}")
							break
						case ' ':
							space_counter++
							cptr1 = IncBufPtr(ptr.BufferIdx)
							if (Buffer[cptr1] != ' ') || (Buffer[cptr1] == -1) {
								if UseUnifiedFont {
									if ptr.Flag == 0 {
										fmt.Fprintf(fptr[1], "{%s%s}", TextModeFont, strings.Repeat("\\ ", space_counter))
									} else {
										fmt.Fprintf(fptr[1], "%s", strings.Repeat("\\ ", space_counter))
									}
								} else {
									if ptr.Flag == 0 {
										fmt.Fprintf(fptr[1], "{%s\\kern%d\\charwd}", TextModeFont, space_counter)
									} else {
										fmt.Fprintf(fptr[1], "\\kern%d\\charwd ", space_counter)
									}
								}
								space_counter = 0
							}
							break
						case '"':
							fmt.Fprintf(fptr[1], "{\\tt \"}")
							break
						case '#':
							fmt.Fprintf(fptr[1], "{\\tt\\#}")
							break
						case '$':
							fmt.Fprintf(fptr[1], "{\\tt\\$}")
							break
						case '%':
							fmt.Fprintf(fptr[1], "{\\tt\\%c}", byte(ptr.Character))
							break
						case '&':
							fmt.Fprintf(fptr[1], "{\\tt\\&}")
							break
						case '*':
							fmt.Fprintf(fptr[1], "{\\tt *}")
							break
						case '-':
							fmt.Fprintf(fptr[1], "{\\tt -}")
							break
						case '/':
							fmt.Fprintf(fptr[1], "{\\tt /}")
							break
						case '<':
							fmt.Fprintf(fptr[1], "{\\tt <}")
							break
						case '>':
							fmt.Fprintf(fptr[1], "{\\tt >}")
							break
						case '\\':
							fmt.Fprintf(fptr[1], "{\\tt\\char92}")
							break
						case '^':
							fmt.Fprintf(fptr[1], "{\\tt\\char'136}")
							break
						case '_':
							if UseUnifiedFont {
								fmt.Fprintf(fptr[1], "{\\tt\\_}")
							} else {
								fmt.Fprintf(fptr[1], "{\\tt\\_\\kern.141em}")
							}
							break
						case '{':
							fmt.Fprintf(fptr[1], "{\\tt\\char'173}")
							break
						case '|':
							fmt.Fprintf(fptr[1], "{\\tt |}")
							break
						case '}':
							fmt.Fprintf(fptr[1], "{\\tt\\char'175}")
							break
						case '~':
							fmt.Fprintf(fptr[1], "{\\tt\\char'176}")
							break
						case 0x7f:
							fmt.Fprintf(fptr[1],
								"{\\sevenrm D\\kern-.15em\\lower.5ex\\hbox{T}}")
							break
						default:

							if (ptr.Flag == 0) && (BFFlag != 0) {
								BFFlag--

								if strings.Contains(TextModeFont, "\\tt") {
									if UseUnifiedFont {
										// Unified font has a real Bold weight — use \textbf instead of
										// the legacy double-strike kerning trick used for cmtt.
										fmt.Fprintf(fptr[1], "{\\textbf{%c}}", byte(ptr.Character))
									} else {
										fptr[1].Write([]byte{byte(ptr.Character)})
										switch ptr.Character {
										case 'e':
											fmt.Fprintf(fptr[1], "\\kern-.445em ")
											fptr[1].Write([]byte{byte(ptr.Character)})
											fmt.Fprintf(fptr[1], "\\kern-.055em ")
											break
										case 'n':
											fmt.Fprintf(fptr[1], "\\kern-.46em ")
											fptr[1].Write([]byte{byte(ptr.Character)})
											fmt.Fprintf(fptr[1], "\\kern-.04em ")
											break
										case 't':
											fmt.Fprintf(fptr[1], "\\kern-.445em ")
											fptr[1].Write([]byte{byte(ptr.Character)})
											fmt.Fprintf(fptr[1], "\\kern-.055em ")
											break
										case 'u':
											fmt.Fprintf(fptr[1], "\\kern-.46em ")
											fptr[1].Write([]byte{byte(ptr.Character)})
											fmt.Fprintf(fptr[1], "\\kern-.04em ")
											break
										default:
											fmt.Fprintf(fptr[1], "\\kern-.455em ")
											fptr[1].Write([]byte{byte(ptr.Character)})
											fmt.Fprintf(fptr[1], "\\kern-.045em ")
										}
									}
								} else {
									if strings.Contains(TextModeFont, "\\bf") {
										fmt.Fprintf(fptr[1], "{\\rm\\mc ")
										fptr[1].Write([]byte{byte(ptr.Character)})
										fmt.Fprintf(fptr[1], "}")
									} else {
										fmt.Fprintf(fptr[1], "{\\bf\\gt ")
										fptr[1].Write([]byte{byte(ptr.Character)})
										fmt.Fprintf(fptr[1], "}")
									}
								}
							} else {
								fptr[1].Write([]byte{byte(ptr.Character)})
							}
						}
						// Correct char_counter for UTF-8 multi-byte characters.
						// Buffer stores raw bytes; CJK chars (3 bytes UTF-8) have display width 2.
						if ptr.Character >= 0x80 && ptr.Character <= 0xBF {
							// UTF-8 continuation byte: no width contribution
						} else if ptr.Character >= 0xE0 && ptr.Character <= 0xF7 {
							// Leading byte of 3-byte (CJK) or 4-byte sequence: full-width = 2
							char_counter += 2
						} else {
							char_counter++
						}
						if (char_counter >= 100) && (WARNFlag2 <= 2) {
							WARNFlag2++

							if WARNFlag2 > 2 {
								if WARNFlag1 <= 5 {
									fmt.Fprintf(os.Stderr, "\n")
								}
								fmt.Fprintf(os.Stderr,
									"Warning: source file contains very long lines;\n")
								fmt.Fprintf(os.Stderr,
									"         their tails are sometimes truncated\n")
							}
						}
						if ((CFlag != 0) || (F77Flag != 0) || (PASFlag != 0)) && (ptr.Character == 0x27) {
							qt_counter++
							qt_counter %= 2
						}
						if ((BASFlag != 0) || (CFlag != 0)) && (ptr.Character == '"') {
							dqt_counter++
							dqt_counter %= 2
						}
						continue
					}

					if ptr.Flag == 2 {
						switch ptr.Character {

						case '\n':
							fmt.Fprintf(fptr[1], "\n")
							LineCounter++
							page_len++
							char_counter = -1
							cptr1 = IncBufPtr(ptr.BufferIdx)
							cptr2 = IncBufPtr(cptr1)
							if ((BASFlag != 0) || (CBLFlag != 0)) && ((Buffer[cptr1] <= ' ') || ((Buffer[cptr1] >= '0') && (Buffer[cptr1] <= '9'))) {

								for Buffer[cptr1] != -1 {
									cptr1 = Fgetc2Buffer(fptr[0])
									cptr2 = IncBufPtr(cptr1)
									if (Buffer[cptr1] <= ' ') && (((Buffer[cptr2] > ' ') && (Buffer[cptr2] < '0')) || (Buffer[cptr2] > '9')) {
										if BASFlag != 0 {
											Buffer[cptr1] = '\n'
										}
										break
									}
								}
							}
							break

						case '#':
							if ((MAKEFlag == 0) && (PERLFlag == 0) && (SHFlag == 0) && (TCLFlag == 0) && (MAPFlag == 0) && (MLABFlag == 0) && (PythonFlag == 0) && (RubyFlag == 0)) || (char_counter != 0) {
								fptr[1].Write([]byte{byte(ptr.Character)})
							}
							break

						case '%':
							if ((REDFlag == 0) && (MLABFlag == 0)) || (char_counter != 0) {
								fptr[1].Write([]byte{byte(ptr.Character)})
							}
							break

						case 0x27:
							if (BASFlag == 0) || (char_counter != 0) {
								fptr[1].Write([]byte{byte(ptr.Character)})
							}
							break

						case '*':
							if ((F77Flag == 0) && (CBLFlag == 0)) || (char_counter != 0) {
								fmt.Fprintf(fptr[1], "*")
							}
							break

						case '/':
							if ((CFlag != 0) || (JAVAFlag != 0)) && (char_counter == 0) {
								cptr1 = IncBufPtr(ptr.BufferIdx)
								if Buffer[cptr1] == '/' {
									Buffer[cptr1] = ' '
								}
							} else {
								fmt.Fprintf(fptr[1], "/")
							}
							break

						case ';':
							if (LISPFlag != 0) && (char_counter == 0) {
								cptr1 = IncBufPtr(ptr.BufferIdx)
								for Buffer[cptr1] == ';' {
									Buffer[cptr1] = ' '
									cptr1 = IncBufPtr(cptr1)
								}
							} else {
								fmt.Fprintf(fptr[1], ";")
							}
							break

						case 'C':
							if (F77Flag == 0) || (char_counter != 0) {
								fmt.Fprintf(fptr[1], "C")
							}
							break

						case 'R':
							cptr1 = IncBufPtr(ptr.BufferIdx)
							cptr2 = IncBufPtr(cptr1)
							if (BASFlag != 0) && (char_counter == 0) && ((Buffer[cptr1] == 'E') || (Buffer[cptr1] == 'e')) && ((Buffer[cptr2] == 'M') || (Buffer[cptr2] == 'm')) {
								Buffer[cptr1] = ' '
								Buffer[cptr2] = ' '
							} else {
								fmt.Fprintf(fptr[1], "R")
							}
							break

						case 'c':
							if (F77Flag == 0) || (char_counter != 0) {
								fmt.Fprintf(fptr[1], "c")
							}
							break

						case 'r':
							cptr1 = IncBufPtr(ptr.BufferIdx)
							cptr2 = IncBufPtr(cptr1)
							if (BASFlag != 0) && (char_counter == 0) && ((Buffer[cptr1] == 'E') || (Buffer[cptr1] == 'e')) && ((Buffer[cptr2] == 'M') || (Buffer[cptr2] == 'm')) {
								Buffer[cptr1] = ' '
								Buffer[cptr2] = ' '
							} else {
								fmt.Fprintf(fptr[1], "r")
							}
							break
						default:
							fptr[1].Write([]byte{byte(ptr.Character)})
						}
						char_counter++
						continue
					}
	} // THIS CLOSES THE FOR LOOP

	if IsLatexMode {
		fmt.Fprintf(fptr[1], "\n\n")
		if len(TextModeFont) >= 3 && TextModeFont[0] == '\\' && TextModeFont[1] == 'r' && TextModeFont[2] == 'm' {
			fmt.Fprintf(fptr[1], "\\rm\n\n")
		} else {
			fmt.Fprintf(fptr[1], "\\rm\\mc\n\n")
		}
		fmt.Fprintf(fptr[1], "\\end{document}\n")
	} else {
		fmt.Fprintf(fptr[1], "\n\n")
		if len(TextModeFont) >= 3 && TextModeFont[0] == '\\' && TextModeFont[1] == 'r' && TextModeFont[2] == 'm' {
			fmt.Fprintf(fptr[1], "\\rm\n\n")
		} else {
			fmt.Fprintf(fptr[1], "\\rm\\mc\n\n")
		}
		fmt.Fprintf(fptr[1], "\\bye\n")
	}
	fmt.Fprintf(os.Stderr, "... done\n")
}
