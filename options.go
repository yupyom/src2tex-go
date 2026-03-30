package main

import (
	"fmt"
	"os"
)

// str_cmp compares a string in Buffer starting at iptr with a Go string.
// Returns 0 if identical.
func str_cmp(iptr int, cptr string) int {
	for i := 0; i < len(cptr); i++ {
		if byte(Buffer[iptr]) != cptr[i] {
			return int(byte(Buffer[iptr])) - int(cptr[i])
		}
		if i+1 == len(cptr) {
			break
		}
		iptr = IncBufPtr(iptr)
	}
	return 0
}

// asc2int parses an integer from Buffer starting at iptr.
func asc2int(iptr int) int {
	value := 0
	for byte(Buffer[iptr]) >= '0' && byte(Buffer[iptr]) <= '9' {
		value = 10*value + int(byte(Buffer[iptr])-'0')
		iptr = IncBufPtr(iptr)
	}
	return value
}

// next_word advances iptr1 to the next word or until it hits iptr2.
func next_word(iptr1, iptr2 int) int {
	if iptr1 == iptr2 {
		return iptr1
	}
	for {
		iptr1 = IncBufPtr(iptr1)
		if iptr1 == iptr2 {
			break
		}
		prev := DecBufPtr(iptr1)

		c1 := byte(Buffer[prev])
		c2 := byte(Buffer[iptr1])

		c1_not_alpha := (c1 < 'a' || c1 > 'z') && (c1 < '0' || c1 > '9')
		c2_is_alpha := (c2 >= 'a' && c2 <= 'z') || (c2 >= '0' && c2 <= '9')

		if c1_not_alpha && c2_is_alpha {
			break
		}
	}
	return iptr1
}

// parse_options parses src2tex inline options like \src2tex{htab 4}
func parse_options(ptr *FlagChar) int {
	brace_counter := 0
	error_flag := 0

	iptr1 := ptr.BufferIdx
	for byte(Buffer[iptr1]) != '\\' && byte(Buffer[iptr1]) != '\n' && Buffer[iptr1] != -1 {
		c := byte(Buffer[iptr1])
		if (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') {
			error_flag++
		}
		iptr1 = IncBufPtr(iptr1)
	}

	if str_cmp(iptr1, "\\src2tex{") != 0 {
		return 0
	}
	if error_flag != 0 {
		fmt.Fprintf(os.Stderr, "\nError: syntax error of src2tex escape sequence\n")
		fmt.Fprintf(os.Stderr, "       \\src2tex{...} is not written properly\n")
		fmt.Fprintf(os.Stderr, "       junks are found before \\src2tex sequence\n")
		os.Exit(1)
	}

	if ptr.Flag == 1 {
		iptr2 := DecBufPtr(iptr1)
		if byte(Buffer[iptr2]) != '{' {
			fmt.Fprintf(os.Stderr, "\nError: syntax error of src2tex escape sequence\n")
			fmt.Fprintf(os.Stderr, "       \\src2tex{...} is not written properly\n")
			fmt.Fprintf(os.Stderr, "       missing a left brace { \n")
			os.Exit(1)
		}
	}

	iptr2 := iptr1
	for i := 0; i < 256; i++ {
		c := byte(Buffer[iptr2])
		if c == '{' {
			brace_counter++
		}
		if c == '}' {
			brace_counter--
		}
		if c == '}' && brace_counter == 0 {
			break
		}
		if c == '\n' || Buffer[iptr2] == -1 {
			break
		}
		iptr2 = IncBufPtr(iptr2)
	}

	if byte(Buffer[iptr2]) != '}' || brace_counter != 0 {
		fmt.Fprintf(os.Stderr, "\nError: syntax error of src2tex escape sequence\n")
		fmt.Fprintf(os.Stderr, "       \\src2tex{...} is not written properly\n")
		fmt.Fprintf(os.Stderr, "       missing right brace } \n")
		os.Exit(1)
	}

	iptr3 := iptr2
	for byte(Buffer[iptr3]) != '\n' && Buffer[iptr3] != -1 {
		c := byte(Buffer[iptr3])
		if (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') {
			fmt.Fprintf(os.Stderr, "\nError: syntax error of src2tex escape sequence\n")
			fmt.Fprintf(os.Stderr, "       \\src2tex{...} is not written properly\n")
			fmt.Fprintf(os.Stderr, "       junks are found after \\src2tex sequence\n")
			os.Exit(1)
		}
		iptr3 = IncBufPtr(iptr3)
	}

	// start parsing
	for i := 1; i <= 8; i++ {
		iptr1 = IncBufPtr(iptr1)
	}
	iptr1 = next_word(iptr1, iptr2)

	for iptr1 != iptr2 {
		if str_cmp(iptr1, "htab") == 0 {
			iptr1 = next_word(iptr1, iptr2)
			if byte(Buffer[iptr1]) >= '0' && byte(Buffer[iptr1]) <= '9' {
				HtabSize = asc2int(iptr1)
			}
		}
		if str_cmp(iptr1, "vtab") == 0 {
			iptr1 = next_word(iptr1, iptr2)
			if byte(Buffer[iptr1]) >= '0' && byte(Buffer[iptr1]) <= '9' {
				VtabSize = asc2int(iptr1)
			}
		}
		if str_cmp(iptr1, "textfont") == 0 {
			iptr1 = next_word(iptr1, iptr2)
			if str_cmp(iptr1, "bf") == 0 {
				TextModeFont = Bold
			}
			if str_cmp(iptr1, "it") == 0 {
				TextModeFont = Italic
			}
			if str_cmp(iptr1, "rm") == 0 {
				TextModeFont = Roman
			}
			if str_cmp(iptr1, "sc") == 0 {
				TextModeFont = SmallCaps
			}
			if str_cmp(iptr1, "sl") == 0 {
				TextModeFont = Slant
			}
			if str_cmp(iptr1, "tt") == 0 {
				TextModeFont = Typewriter
			}
		}
		if str_cmp(iptr1, "texfont") == 0 {
			iptr1 = next_word(iptr1, iptr2)
			if str_cmp(iptr1, "bf") == 0 {
				TeXModeFont = Bold
			}
			if str_cmp(iptr1, "it") == 0 {
				TeXModeFont = Italic
			}
			if str_cmp(iptr1, "rm") == 0 {
				TeXModeFont = Roman
			}
			if str_cmp(iptr1, "sc") == 0 {
				TeXModeFont = SmallCaps
			}
			if str_cmp(iptr1, "sl") == 0 {
				TeXModeFont = Slant
			}
			if str_cmp(iptr1, "tt") == 0 {
				TeXModeFont = Typewriter
			}
		}
		iptr1 = next_word(iptr1, iptr2)
	}

	return 1
}
