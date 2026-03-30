package main

import (
	"fmt"
	"os"
)

var (
	gfStatFlag     int
	gfParenCounter int
	gfBraceCounter int
	gfComFlag      int
	gfCharCounter  int
	gfQtCounter    int
	gfBqtCounter   int
	gfDqtCounter   int
	gfSlashCounter int
	gfRgxCounter   int
	gfWarnFlag     int

	tfMathFlag     int
	tfDsFlag       int
	tfEscFlag      int
	tfRcsFlag      int
	tfBraceCounter int
)

func ResetModFlags() {
	gfStatFlag = 0
	gfParenCounter = 0
	gfBraceCounter = 0
	gfComFlag = 0
	gfCharCounter = 0
	gfQtCounter = 0
	gfBqtCounter = 0
	gfDqtCounter = 0
	gfSlashCounter = 0
	gfRgxCounter = 0
	gfWarnFlag = 0

	tfMathFlag = 0
	tfDsFlag = 0
	tfEscFlag = 0
	tfRcsFlag = 0
	tfBraceCounter = 0
}

func GetCommentFlag(bufPtr int) int {
	if gfStatFlag == 0 {
		gfStatFlag++
		InitLangFlag() // Since set_lang_flag takes bufPtr, but init_lang_flag handles it globally mostly. Wait, in C it is `set_lang_flag(buf_ptr)`.
		SetLangFlag(bufPtr)
	}

	ptr := DecBufPtr(bufPtr)
	c_1 := byte(Buffer[ptr])
	ptr = DecBufPtr(ptr)
	c_2 := byte(Buffer[ptr])
	c0 := byte(Buffer[bufPtr])

	ptr1 := IncBufPtr(bufPtr)
	c1 := byte(Buffer[ptr1])
	ptr2 := IncBufPtr(ptr1)
	c2 := byte(Buffer[ptr2])
	ptr3 := IncBufPtr(ptr2)
	c3 := byte(Buffer[ptr3])
	ptr4 := IncBufPtr(ptr3)
	c4 := byte(Buffer[ptr4])
	ptr5 := IncBufPtr(ptr4)
	c5 := byte(Buffer[ptr5])
	ptr6 := IncBufPtr(ptr5)
	c6 := byte(Buffer[ptr6])

	if BASFlag == 1 {
		if c0 == '"' && gfComFlag == 0 {
			gfDqtCounter++
			gfDqtCounter %= 2
		}
		if gfDqtCounter == 0 {
			if gfComFlag == 1 && c_1 == '\n' {
				gfComFlag = 0
			}
			if gfCharCounter == 0 {
				if gfComFlag == 0 && c_1 <= ' ' && c0 == 0x27 {
					gfComFlag = 1
					if c1 == '{' && c2 == '\\' && c3 == ' ' {
						Buffer[bufPtr] = ' '
					}
				}
				if gfComFlag == 0 && c_1 <= ' ' && (c0 == 'R' || c0 == 'r') && (c1 == 'E' || c1 == 'e') && (c2 == 'M' || c2 == 'm') {
					gfComFlag = 1
				}
			}
		}
		if c0 == '\n' {
			gfCharCounter = 0
		} else {
			if gfCharCounter != 0 {
				gfCharCounter++
			} else if c0 > '9' || (c0 > ' ' && c0 < '0') {
				gfCharCounter++
			}
		}
		return gfComFlag
	}

	if CFlag == 1 || MACFlag == 1 {
		if ((c_2 == '\\' && c_1 == '\\') || c_1 != '\\') && c0 == '"' && (c_1 != 0x27 || c1 != 0x27) && gfComFlag == 0 {
			gfDqtCounter++
			gfDqtCounter %= 2
		}
		if gfDqtCounter == 0 {
			if gfComFlag == 1 && c_2 == '*' && c_1 == '/' {
				gfComFlag = 0
			}
			if gfComFlag == 3 && c_1 == '\n' {
				gfComFlag = 0
			}
			if gfComFlag == 0 && c0 == '/' && c1 == '*' {
				gfComFlag = 1
			}
			if gfComFlag == 0 && c0 == '/' && c1 == '/' {
				gfComFlag = 3
			}
		}
		return gfComFlag % 2
	}

	if CBLFlag == 1 {
		if gfComFlag == 1 && gfCharCounter == 0 && c_1 == '\n' {
			gfComFlag = 0
		}
		if gfComFlag == 0 && gfCharCounter == 0 && (c6 == '*' || c6 == '/') {
			gfComFlag = 1
		}
		if c0 == '\n' {
			gfCharCounter = 0
		} else {
			gfCharCounter++
		}
		return gfComFlag
	}

	if F77Flag == 1 {
		if c0 == 0x27 && gfComFlag == 0 {
			gfQtCounter++
			gfQtCounter %= 2
		}
		if gfQtCounter == 0 {
			if gfComFlag == 1 && (c_1 == '\r' || c_1 == '\n') {
				gfComFlag = 0
			}
			if gfComFlag == 1 && (c_1 == 'C' || c_1 == 'c' || c_1 == '*') && (c0 == '\r' || c0 == '\n') {
				gfComFlag = 0
			}
			if gfCharCounter == 0 {
				if gfComFlag == 0 && (c0 == 'C' || c0 == 'c') {
					gfComFlag = 1
				}
				if gfComFlag == 0 && c0 == '*' {
					gfComFlag = 1
					if c1 == '{' && c2 == '\\' && c3 == ' ' {
						Buffer[bufPtr] = ' '
					}
				}
			}
		}
		if c0 == '\n' {
			gfCharCounter = 0
		} else {
			gfCharCounter++
		}
		return gfComFlag
	}

	if HTMLFlag == 1 {
		if ((c_2 == '\\' && c_1 == '\\') || c_1 != '\\') && c0 == '"' && (c_1 != 0x27 || c1 != 0x27) && gfComFlag == 0 {
			gfDqtCounter++
			gfDqtCounter %= 2
		}
		if gfDqtCounter == 0 {
			if gfComFlag == 1 && c_2 == '-' && c_1 == '>' {
				gfComFlag = 0
			}
			// Bug fix: The original C code had `c2 = '-'`, which was assigning. We fix this to `c2 == '-'`.
			if gfComFlag == 0 && c0 == '<' && c1 == '!' && c2 == '-' {
				gfComFlag = 1
			}
		}
		return gfComFlag
	}

	if JAVAFlag == 1 || JSFlag == 1 || TSFlag == 1 || GoFlag == 1 || RustFlag == 1 || KotlinFlag == 1 || SwiftFlag == 1 {
		if ((c_2 == '\\' && c_1 == '\\') || c_1 != '\\') && c0 == '"' && (c_1 != 0x27 || c1 != 0x27) && gfComFlag == 0 {
			gfDqtCounter++
			gfDqtCounter %= 2
		}
		if gfDqtCounter == 0 {
			if gfComFlag == 1 && c_2 == '*' && c_1 == '/' {
				gfComFlag = 0
			}
			if gfComFlag == 3 && c_1 == '\n' {
				gfComFlag = 0
			}
			if gfComFlag == 0 && c0 == '/' && c1 == '*' {
				gfComFlag = 1
			}
			if gfComFlag == 0 && c0 == '/' && c1 == '/' {
				gfComFlag = 3
			}
		}
		return gfComFlag % 2
	}

	if LISPFlag == 1 {
		if c0 == '"' && gfComFlag == 0 {
			gfDqtCounter++
			gfDqtCounter %= 2
		}
		if gfDqtCounter == 0 {
			if gfComFlag == 1 && c_1 == '\n' {
				gfComFlag = 0
			}
			if gfComFlag == 0 && c0 == ';' {
				gfComFlag = 1
				if c1 == '{' && c2 == '\\' && c3 == ' ' {
					Buffer[bufPtr] = ' '
				}
			}
		}
		return gfComFlag
	}

	if MAKEFlag == 1 {
		if c0 == 0x27 && gfComFlag == 0 && gfBqtCounter == 0 && gfDqtCounter == 0 {
			gfQtCounter++
			gfQtCounter %= 2
		}
		if c0 == 0x60 && gfComFlag == 0 && gfQtCounter == 0 && gfDqtCounter == 0 {
			gfBqtCounter++
			gfBqtCounter %= 2
		}
		if c0 == '"' && gfComFlag == 0 && gfQtCounter == 0 && gfBqtCounter == 0 {
			gfDqtCounter++
			gfDqtCounter %= 2
		}

		if gfQtCounter == 0 && gfBqtCounter == 0 && gfDqtCounter == 0 {
			if gfComFlag == 1 && c_2 != '\\' && (c_1 == '\r' || c_1 == '\n') {
				gfComFlag = 0
			}
			if gfComFlag == 1 && c_1 == '#' && (c0 == '\r' || c0 == '\n') {
				gfComFlag = 0
			}
			if gfComFlag == 0 && c0 == '#' {
				gfComFlag = 1
				if c1 == '{' && c2 == '\\' && c3 == ' ' {
					Buffer[bufPtr] = ' '
				}
			}
		}
		return gfComFlag
	}

	if MPADFlag == 1 {
		if c0 == 0x27 && gfComFlag == 0 && gfBqtCounter == 0 && gfDqtCounter == 0 {
			gfQtCounter++
			gfQtCounter %= 2
		}
		if c0 == 0x60 && gfComFlag == 0 && gfQtCounter == 0 && gfDqtCounter == 0 {
			gfBqtCounter++
			gfBqtCounter %= 2
		}
		if c0 == '"' && gfComFlag == 0 && gfQtCounter == 0 && gfBqtCounter == 0 {
			gfDqtCounter++
			gfDqtCounter %= 2
		}

		if gfQtCounter == 0 && gfBqtCounter == 0 && gfDqtCounter == 0 {
			if gfComFlag == 1 && c_1 != '\\' && c0 == '#' {
				gfComFlag = 0
			}
			if gfComFlag == 0 && (c_1 <= 0x20 || c_1 == ';') && c0 == '#' {
				gfComFlag = 1
			}
		}
		return gfComFlag
	}

	if PASFlag == 1 {
		if c0 == 0x27 && gfComFlag == 0 {
			gfQtCounter++
			gfQtCounter %= 2
		}
		if gfQtCounter == 0 && c_1 != '\\' && c0 == '{' {
			gfBraceCounter++
		}
		if gfQtCounter == 0 && c_1 != '\\' && c0 == '}' {
			gfBraceCounter--
		}
		if gfQtCounter == 0 {
			if gfBraceCounter == 0 && gfComFlag == 1 && c_1 == '}' {
				gfComFlag = 0
			}
			if gfBraceCounter == 1 && gfComFlag == 0 && c0 == '{' {
				gfComFlag = 1
			}
			if gfComFlag == 3 && c_2 == '*' && c_1 == ')' {
				gfComFlag = 0
			}
			if gfComFlag == 0 && c0 == '(' && c1 == '*' {
				gfComFlag = 3
			}
		}
		return gfComFlag % 2
	}

	if PERLFlag == 1 {
		if ((c_2 != '$' && c_1 == 'm' && c0 == '!') || (c_1 != '\\' && c0 == '?')) && gfComFlag == 0 && gfRgxCounter == 0 && gfParenCounter == 0 && gfBraceCounter == 0 && gfQtCounter == 0 && gfBqtCounter == 0 && gfDqtCounter == 0 && gfSlashCounter == 0 {
			gfRgxCounter++
		}
		if c_1 != '\\' && (c0 == '!' || c0 == '?') && gfComFlag == 0 && gfRgxCounter == 1 && gfParenCounter == 0 && gfBraceCounter == 0 && gfQtCounter == 0 && gfBqtCounter == 0 && gfDqtCounter == 0 && gfSlashCounter == 0 {
			gfRgxCounter--
		}
		if c_1 != '\\' && (c0 == '(' || c0 == ')') && gfComFlag == 0 && gfRgxCounter == 0 && gfBraceCounter == 0 && gfQtCounter == 0 && gfDqtCounter == 0 && gfBqtCounter == 0 && gfSlashCounter == 0 {
			gfParenCounter++
			gfParenCounter %= 2
		}
		if c_1 != '\\' && c0 == '{' && ((c_1 >= 'A' && c_1 <= 'Z') || (c_1 >= 'a' && c_1 <= 'z')) && gfComFlag == 0 && gfRgxCounter == 0 && gfParenCounter == 0 && gfBraceCounter == 0 && gfQtCounter == 0 && gfBqtCounter == 0 && gfDqtCounter == 0 && gfSlashCounter == 0 {
			gfBraceCounter++
		}
		if c_1 != '\\' && c0 == '}' && gfComFlag == 0 && gfRgxCounter == 0 && gfParenCounter == 0 && gfBraceCounter == 1 && gfQtCounter == 0 && gfBqtCounter == 0 && gfDqtCounter == 0 && gfSlashCounter == 0 {
			gfBraceCounter--
		}
		if c_1 != '\\' && c0 == 0x27 && gfComFlag == 0 && gfRgxCounter == 0 && gfParenCounter == 0 && gfBraceCounter == 0 && gfBqtCounter == 0 && gfDqtCounter == 0 && gfSlashCounter == 0 {
			gfQtCounter++
			gfQtCounter %= 2
		}
		if ((c_1 == 'q' && c0 == '|') || (c_1 != '\\' && c0 == '|')) && gfComFlag == 0 && gfRgxCounter == 0 && gfParenCounter == 0 && gfBraceCounter == 0 && gfBqtCounter == 0 && gfDqtCounter == 0 && gfSlashCounter == 0 {
			gfQtCounter += 2
			gfQtCounter %= 4
		}
		if c_1 != '\\' && c0 == 0x60 && gfComFlag == 0 && gfRgxCounter == 0 && gfParenCounter == 0 && gfBraceCounter == 0 && gfQtCounter == 0 && gfBqtCounter == 0 && gfDqtCounter == 0 && gfSlashCounter == 0 {
			gfBqtCounter = 1
		}
		if c_1 != '\\' && c0 == 0x60 && gfComFlag == 0 && gfRgxCounter == 0 && gfParenCounter == 0 && gfBraceCounter == 0 && gfQtCounter == 0 && gfBqtCounter == 1 && gfDqtCounter == 0 && gfSlashCounter == 0 {
			gfBqtCounter = 0
		}
		if c_1 != '\\' && c0 == '"' && gfComFlag == 0 && gfRgxCounter == 0 && gfParenCounter == 0 && gfBraceCounter == 0 && gfQtCounter == 0 && gfBqtCounter == 0 && gfSlashCounter == 0 {
			gfDqtCounter++
			gfDqtCounter %= 2
		}
		if c_1 == ' ' && c0 == '/' && gfComFlag == 0 && gfRgxCounter == 0 && gfParenCounter == 0 && gfBraceCounter == 0 && gfQtCounter == 0 && gfBqtCounter == 0 && gfDqtCounter == 0 && gfSlashCounter == 0 {
			gfSlashCounter = 2
		}
		if c_2 != '$' && (c_1 >= 'a' && c_1 <= 'z') && c0 == '/' && gfComFlag == 0 && gfRgxCounter == 0 && gfParenCounter == 0 && gfBraceCounter == 0 && gfQtCounter == 0 && gfBqtCounter == 0 && gfDqtCounter == 0 && gfSlashCounter == 0 {
			gfSlashCounter = 3
		}
		if ((c_1 != '\\' || (c_2 == '\\' && c_1 == '\\')) && c0 == '/') && gfComFlag == 0 && gfRgxCounter == 0 && gfParenCounter == 0 && gfBraceCounter == 0 && gfQtCounter == 0 && gfBqtCounter == 0 && gfDqtCounter == 0 && gfSlashCounter != 0 {
			gfSlashCounter--
		}

		if gfRgxCounter == 0 && gfParenCounter == 0 && gfBraceCounter == 0 && gfQtCounter == 0 && gfBqtCounter == 0 && gfDqtCounter == 0 && gfSlashCounter == 0 {
			if gfComFlag == 1 && (c_1 == '\r' || c_1 == '\n') {
				gfComFlag = 0
			}
			if gfComFlag == 1 && c_1 == '#' && (c0 == '\r' || c0 == '\n') {
				gfComFlag = 0
			}

			// Handle stat_flag++ logic properly:
			// "&& ((stat_flag++ == 1) || (c_1 == '\t') || (c_1 == '\n') || (c_1 == ' ') || (c_1 == ';')) && (c0 == '#')"
			conditionMet := false
			if gfStatFlag == 1 {
				gfStatFlag++
				conditionMet = true
			} else {
				gfStatFlag++
				if c_1 == '\t' || c_1 == '\n' || c_1 == ' ' || c_1 == ';' {
					conditionMet = true
				}
			}

			if gfComFlag == 0 && conditionMet && c0 == '#' {
				gfComFlag = 1
				if c1 == '{' && c2 == '\\' && c3 == ' ' {
					Buffer[bufPtr] = ' '
				}
			}
		}

		if c0 == '\n' {
			gfCharCounter = 0
		} else {
			gfCharCounter++
		}
		return gfComFlag
	}

	if SHFlag == 1 || PythonFlag == 1 || RubyFlag == 1 {
		if c0 == 0x27 && gfComFlag == 0 && gfBqtCounter == 0 && gfDqtCounter == 0 {
			gfQtCounter++
			gfQtCounter %= 2
		}
		if c0 == 0x60 && gfComFlag == 0 && gfQtCounter == 0 && gfDqtCounter == 0 {
			gfBqtCounter++
			gfBqtCounter %= 2
		}
		if c0 == '"' && gfComFlag == 0 && gfQtCounter == 0 && gfBqtCounter == 0 {
			gfDqtCounter++
			gfDqtCounter %= 2
		}
		if gfQtCounter == 0 && gfBqtCounter == 0 && gfDqtCounter == 0 {
			if gfComFlag == 1 && c_2 != '\\' && (c_1 == '\r' || c_1 == '\n') {
				gfComFlag = 0
			}
			if gfComFlag == 1 && c_1 == '#' && (c0 == '\r' || c0 == '\n') {
				gfComFlag = 0
			}
			if gfComFlag == 0 && c_1 != '$' && c0 == '#' {
				gfComFlag = 1
				if c1 == '{' && c2 == '\\' && c3 == ' ' {
					Buffer[bufPtr] = ' '
				}
			}
		}
		return gfComFlag
	}

	if TCLFlag == 1 || MAPFlag == 1 {
		if c0 == 0x27 && gfComFlag == 0 && gfBqtCounter == 0 && gfDqtCounter == 0 {
			gfQtCounter++
			gfQtCounter %= 2
		}
		if TCLFlag == 1 {
			if c0 == 0x60 && gfComFlag == 0 && gfQtCounter == 0 && gfDqtCounter == 0 {
				gfBqtCounter++
				gfBqtCounter %= 2
			}
			if c0 == '"' && gfComFlag == 0 && gfQtCounter == 0 && gfBqtCounter == 0 {
				gfDqtCounter++
				gfDqtCounter %= 2
			}
		}
		if gfQtCounter == 0 && gfBqtCounter == 0 && gfDqtCounter == 0 {
			if gfComFlag == 1 && c_2 != '\\' && (c_1 == '\r' || c_1 == '\n') {
				gfComFlag = 0
			}
			if gfComFlag == 1 && c_1 == '#' && (c0 == '\r' || c0 == '\n') {
				gfComFlag = 0
			}
			if gfComFlag == 0 && c0 == '#' {
				gfComFlag = 1
				if c1 == '{' && c2 == '\\' && c3 == ' ' {
					Buffer[bufPtr] = ' '
				}
			}
		}
		return gfComFlag
	}

	if MATFlag == 1 {
		if c0 == '"' && gfComFlag == 0 {
			gfDqtCounter++
			gfDqtCounter %= 2
		}
		if gfDqtCounter == 0 {
			if gfComFlag == 1 && c_2 == '*' && c_1 == ')' {
				gfComFlag = 0
			}
			if gfComFlag == 0 && c0 == '(' && c1 == '*' {
				gfComFlag = 1
			}
		}
		return gfComFlag
	}

	if MLABFlag == 1 {
		if c0 == 0x27 && gfQtCounter == 1 {
			gfQtCounter = 0
		}
		if (c_1 < '0' || (c_1 > '9' && c_1 < 'A') || (c_1 > 'Z' && c_1 < 'a') || c_1 > 'z') && c_1 != '\\' && c_1 != '.' && c_1 != ')' && c_1 != ']' && c0 == 0x27 && gfComFlag == 0 {
			gfQtCounter = 1
		}
		if c_1 != '\\' && c0 == '"' && gfComFlag == 0 {
			gfDqtCounter++
			gfDqtCounter %= 2
		}
		if gfQtCounter == 0 && gfDqtCounter == 0 {
			if gfComFlag == 1 && c_1 == '\n' {
				gfComFlag = 0
			}
			if gfComFlag == 0 && (c0 == '#' || c0 == '%') {
				gfComFlag = 1
				if c1 == '{' && c2 == '\\' && c3 == ' ' {
					Buffer[bufPtr] = ' '
				}
			}
		}
		return gfComFlag
	}

	if REDFlag == 1 {
		if c0 == '"' && gfComFlag == 0 {
			gfDqtCounter++
			gfDqtCounter %= 2
		}
		if gfDqtCounter == 0 {
			if gfComFlag == 1 && c_1 == '\n' {
				gfComFlag = 0
			}
			if gfComFlag == 0 && c0 == '%' {
				gfComFlag = 1
				if c1 == '{' && c2 == '\\' && c3 == ' ' {
					Buffer[bufPtr] = ' '
				}
			}
			if gfComFlag == 3 && c_1 == ';' {
				gfComFlag = 0
			}
			if gfComFlag == 0 && (c0 == 'C' || c0 == 'c') && (c1 == 'O' || c1 == 'o') && (c2 == 'M' || c2 == 'm') && (c3 == 'M' || c3 == 'm') && (c4 == 'E' || c4 == 'e') && (c5 == 'N' || c5 == 'n') && (c6 == 'T' || c6 == 't') {
				gfComFlag = 3
			}
		}
		return gfComFlag % 2
	}

	return gfComFlag
}

func GetTexFlag(bufPtr int) int {
	a := make([]byte, 4)
	b := make([]byte, 10)

	ptr := bufPtr
	a[0] = byte(Buffer[ptr])
	ptr = DecBufPtr(ptr)
	a[1] = byte(Buffer[ptr])
	ptr = DecBufPtr(ptr)
	a[2] = byte(Buffer[ptr])
	ptr = DecBufPtr(ptr)
	a[3] = byte(Buffer[ptr])

	if IsLatexMode {
		if tfRcsFlag != 0 {
			if a[2] != '\\' && a[2] != '$' && a[1] == '$' && a[0] != '$' {
				tfRcsFlag--
			}
		}
		if tfMathFlag != 0 {
			if tfDsFlag == 0 {
				if a[1] == '$' && a[0] != '$' {
					if a[2] != '\\' && a[2] != '$' {
						tfMathFlag--
					}
				}
				if a[2] == '\\' && a[1] == ')' {
					tfMathFlag = 0
				}
			} else {
				if a[2] == '\\' && a[1] == ']' {
					tfMathFlag = 0
					tfDsFlag = 0
				}
			}
		}
	} else {
		if tfRcsFlag != 0 {
			if a[2] != '\\' && a[2] != '$' && a[1] == '$' && a[0] != '$' {
				tfRcsFlag--
			}
		}
		if tfMathFlag != 0 {
			if a[1] == '$' && a[0] != '$' {
				if tfDsFlag == 0 {
					if a[2] != '\\' && a[2] != '$' {
						tfMathFlag--
					}
				} else {
					if a[3] != '\\' && a[3] != '$' && a[2] == '$' {
						tfMathFlag--
					}
					if tfMathFlag == 0 {
						tfDsFlag = 0
					}
				}
			}
		}
	}

	if tfEscFlag != 0 {
		if a[1] != '\\' && a[0] == '{' {
			tfBraceCounter++
		}
		if a[1] != '\\' && a[0] == '}' {
			tfBraceCounter--
		}
		if tfEscFlag == 2 {
			if tfBraceCounter == 0 {
				if a[0] != '}' {
					fmt.Fprintf(os.Stderr, "\nError: brace counter error !!\n")
					os.Exit(1)
				}
				tfEscFlag--
				Buffer[bufPtr] = ' '
			}
		} else {
			tfEscFlag--
		}
	}

	ptr = bufPtr
	b[0] = byte(Buffer[ptr])
	ptr = IncBufPtr(ptr)
	b[1] = byte(Buffer[ptr])
	ptr = IncBufPtr(ptr)
	b[2] = byte(Buffer[ptr])
	ptr = IncBufPtr(ptr)
	b[3] = byte(Buffer[ptr])
	ptr = IncBufPtr(ptr)
	b[4] = byte(Buffer[ptr])
	ptr = IncBufPtr(ptr)
	b[5] = byte(Buffer[ptr])
	ptr = IncBufPtr(ptr)
	b[6] = byte(Buffer[ptr])
	ptr = IncBufPtr(ptr)
	b[7] = byte(Buffer[ptr])
	ptr = IncBufPtr(ptr)
	b[8] = byte(Buffer[ptr])
	ptr = IncBufPtr(ptr)
	b[9] = byte(Buffer[ptr])

	if IsLatexMode {
		if tfMathFlag == 0 {
			if a[1] != '\\' && a[1] != '$' && b[0] == '$' {
				if b[1] != '$' {
					if b[1] == 'A' && b[2] == 'u' && b[3] == 't' && b[4] == 'h' && b[5] == 'o' && b[6] == 'r' && (b[7] == ':' || b[7] == '$') {
						tfRcsFlag = 2
					}
					if b[1] == 'D' && b[2] == 'a' && b[3] == 't' && b[4] == 'e' && (b[5] == ':' || b[5] == '$') {
						tfRcsFlag = 2
					}
					if b[1] == 'H' && b[2] == 'e' && b[3] == 'a' && b[4] == 'd' && b[5] == 'e' && b[6] == 'r' && (b[7] == ':' || b[7] == '$') {
						tfRcsFlag = 2
					}
					if b[1] == 'I' && b[2] == 'd' && (b[3] == ':' || b[3] == '$') {
						tfRcsFlag = 2
					}
					if b[1] == 'L' && b[2] == 'o' && b[3] == 'c' && b[4] == 'k' && b[5] == 'e' && b[6] == 'r' && (b[7] == ':' || b[7] == '$') {
						tfRcsFlag = 2
					}
					if b[1] == 'L' && b[2] == 'o' && b[3] == 'g' && (b[4] == ':' || b[4] == '$') {
						tfRcsFlag = 2
					}
					if b[1] == 'R' && b[2] == 'e' && b[3] == 'v' && b[4] == 'i' && b[5] == 's' && b[6] == 'i' && b[7] == 'o' && b[8] == 'n' && (b[9] == ':' || b[9] == '$') {
						tfRcsFlag = 2
					}
					if b[1] == 'R' && b[2] == 'C' && b[3] == 'S' && b[4] == 'f' && b[5] == 'i' && b[6] == 'l' && b[7] == 'e' && (b[8] == ':' || b[8] == '$') {
						tfRcsFlag = 2
					}
					if b[1] == 'S' && b[2] == 'o' && b[3] == 'n' && b[4] == 'y' && b[5] == 'I' && b[6] == 'd' && (b[7] == ':' || b[7] == '$') {
						tfRcsFlag = 2
					}
					if b[1] == 'S' && b[2] == 'o' && b[3] == 'u' && b[4] == 'r' && b[5] == 'c' && b[6] == 'e' && (b[7] == ':' || b[7] == '$') {
						tfRcsFlag = 2
					}
					if b[1] == 'S' && b[2] == 't' && b[3] == 'a' && b[4] == 't' && b[5] == 'e' && (b[6] == ':' || b[6] == '$') {
						tfRcsFlag = 2
					}
				}
				if BASFlag == 0 && PERLFlag == 0 {
					if b[1] != '$' {
						tfMathFlag = 2
					}
				} else {
					if BASFlag != 0 && b[1] != '$' && (a[1] < '0' || a[1] > '9') && (a[1] < 'A' || a[1] > 'Z') && (a[1] < 'a' || a[1] > 'z') {
						tfMathFlag = 2
					}
					if PERLFlag != 0 && b[1] != '$' && (b[1] < 'A' || b[1] > 'Z') && (b[1] < 'a' || b[1] > 'z') {
						tfMathFlag = 2
					}
				}
			}
		}

		if tfMathFlag != 0 && (MAKEFlag != 0 || PERLFlag != 0 || SHFlag != 0 || TCLFlag != 0) {
			ptr = bufPtr
			texCounter := 0
			var i int
			for i = 0; i < FormulaLenMax; i++ {
				c0 := byte(Buffer[ptr])
				ptr1 := IncBufPtr(ptr)
				c1 := byte(Buffer[ptr1])
				ptr2 := IncBufPtr(ptr1)
				c2 := byte(Buffer[ptr2])
				if c0 == '\n' && c1 != '#' {
					break
				}
				if c0 != '\\' && c1 == '$' && c2 <= ' ' {
					break
				}
				if c0 == '_' || c0 == '^' || c0 == '\\' {
					texCounter++
				}
				ptr = ptr1
			}
			if i >= FormulaLenMax {
				tfMathFlag = 0
			}

			// Needs to explicitly re-fetch after loop depending on how loop broke
			// Actually we can just do what C did: it broke but c0, c1, c2 correspond to the last accessed
			c0 := byte(Buffer[ptr])
			ptr1 := IncBufPtr(ptr)
			c1 := byte(Buffer[ptr1])
			ptr2 := IncBufPtr(ptr1)
			c2 := byte(Buffer[ptr2])

			if c0 == '\n' && c1 != '#' {
				tfMathFlag = 0
			}
			if c0 != '\\' && c1 == '$' && c2 <= ' ' && texCounter == 0 {
				tfMathFlag = 0
			}
		}

		if b[1] <= ' ' && b[0] == '\\' && b[1] == '(' && ((b[2] >= '0' && b[2] <= '9') || (b[2] >= 'a' && b[2] <= 'z') || (b[2] == '\\' && b[3] >= 'a' && b[3] <= 'z') || b[2] == '{') {
			tfMathFlag = 2
		}
		if b[1] <= ' ' && b[0] == '\\' && b[1] == '[' && ((b[2] >= '0' && b[2] <= '9') || (b[2] >= 'a' && b[2] <= 'z') || (b[2] == '\\' && b[3] >= 'a' && b[3] <= 'z') || b[2] == '{') {
			tfMathFlag = 2
			tfDsFlag = 1
		}

	} else {
		if tfMathFlag == 0 {
			if a[1] != '\\' && a[1] != '$' && b[0] == '$' {
				if b[1] != '$' {
					if b[1] == 'A' && b[2] == 'u' && b[3] == 't' && b[4] == 'h' && b[5] == 'o' && b[6] == 'r' && (b[7] == ':' || b[7] == '$') {
						tfRcsFlag = 2
					}
					if b[1] == 'D' && b[2] == 'a' && b[3] == 't' && b[4] == 'e' && (b[5] == ':' || b[5] == '$') {
						tfRcsFlag = 2
					}
					if b[1] == 'H' && b[2] == 'e' && b[3] == 'a' && b[4] == 'd' && b[5] == 'e' && b[6] == 'r' && (b[7] == ':' || b[7] == '$') {
						tfRcsFlag = 2
					}
					if b[1] == 'I' && b[2] == 'd' && (b[3] == ':' || b[3] == '$') {
						tfRcsFlag = 2
					}
					if b[1] == 'L' && b[2] == 'o' && b[3] == 'c' && b[4] == 'k' && b[5] == 'e' && b[6] == 'r' && (b[7] == ':' || b[7] == '$') {
						tfRcsFlag = 2
					}
					if b[1] == 'L' && b[2] == 'o' && b[3] == 'g' && (b[4] == ':' || b[4] == '$') {
						tfRcsFlag = 2
					}
					if b[1] == 'R' && b[2] == 'e' && b[3] == 'v' && b[4] == 'i' && b[5] == 's' && b[6] == 'i' && b[7] == 'o' && b[8] == 'n' && (b[9] == ':' || b[9] == '$') {
						tfRcsFlag = 2
					}
					if b[1] == 'R' && b[2] == 'C' && b[3] == 'S' && b[4] == 'f' && b[5] == 'i' && b[6] == 'l' && b[7] == 'e' && (b[8] == ':' || b[8] == '$') {
						tfRcsFlag = 2
					}
					if b[1] == 'S' && b[2] == 'o' && b[3] == 'n' && b[4] == 'y' && b[5] == 'I' && b[6] == 'd' && (b[7] == ':' || b[7] == '$') {
						tfRcsFlag = 2
					}
					if b[1] == 'S' && b[2] == 'o' && b[3] == 'u' && b[4] == 'r' && b[5] == 'c' && b[6] == 'e' && (b[7] == ':' || b[7] == '$') {
						tfRcsFlag = 2
					}
					if b[1] == 'S' && b[2] == 't' && b[3] == 'a' && b[4] == 't' && b[5] == 'e' && (b[6] == ':' || b[6] == '$') {
						tfRcsFlag = 2
					}
				}
				if BASFlag == 0 && PERLFlag == 0 {
					if b[1] != '$' {
						tfMathFlag = 2
					}
				} else {
					if BASFlag != 0 && b[1] != '$' && (a[1] < '0' || a[1] > '9') && (a[1] < 'A' || a[1] > 'Z') && (a[1] < 'a' || a[1] > 'z') {
						tfMathFlag = 2
					}
					if PERLFlag != 0 && b[1] != '$' && (b[1] < 'A' || b[1] > 'Z') && (b[1] < 'a' || b[1] > 'z') {
						tfMathFlag = 2
					}
				}
			}

			if tfMathFlag != 0 && (MAKEFlag != 0 || PERLFlag != 0 || SHFlag != 0 || TCLFlag != 0) {
				ptr = bufPtr
				texCounter := 0
				var i int
				for i = 0; i < FormulaLenMax; i++ {
					c0 := byte(Buffer[ptr])
					ptr1 := IncBufPtr(ptr)
					c1 := byte(Buffer[ptr1])
					ptr2 := IncBufPtr(ptr1)
					c2 := byte(Buffer[ptr2])
					if c0 == '\n' && c1 != '#' {
						break
					}
					if c0 != '\\' && c1 == '$' && c2 <= ' ' {
						break
					}
					if c0 == '_' || c0 == '^' || c0 == '\\' {
						texCounter++
					}
					ptr = ptr1
				}
				if i >= FormulaLenMax {
					tfMathFlag = 0
				}

				c0 := byte(Buffer[ptr])
				ptr1 := IncBufPtr(ptr)
				c1 := byte(Buffer[ptr1])
				ptr2 := IncBufPtr(ptr1)
				c2 := byte(Buffer[ptr2])

				if c0 == '\n' && c1 != '#' {
					tfMathFlag = 0
				}
				if c0 != '\\' && c1 == '$' && c2 <= ' ' && texCounter == 0 {
					tfMathFlag = 0
				}
			}

			if b[1] == '$' && b[2] != '$' {
				tfMathFlag = 2
				tfDsFlag = 1
			}
		}
	}

	if tfMathFlag == 0 && tfEscFlag == 0 {
		if a[1] != '\\' && b[0] == '{' {
			if b[1] == '\\' || (b[1] == '{' && b[2] == '\\') {
				tfBraceCounter = 1
				tfEscFlag = 2
				Buffer[bufPtr] = ' '
			}
			if b[1] == '\\' && b[2] == ' ' {
				ptr1 := IncBufPtr(bufPtr)
				Buffer[ptr1] = '{'
				ptr2 := IncBufPtr(ptr1)
				Buffer[ptr2] = '}'
			}
			if b[1] == '{' && b[2] == '\\' && b[3] == ' ' {
				ptr1 := IncBufPtr(bufPtr)
				ptr2 := IncBufPtr(ptr1)
				Buffer[ptr2] = '{'
				ptr3 := IncBufPtr(ptr2)
				Buffer[ptr3] = '}'
			}
		}
	}

	texFlag := 0
	if tfRcsFlag != 0 {
		if tfEscFlag != 0 {
			texFlag = 1
		}
	} else {
		if tfMathFlag != 0 || tfEscFlag != 0 {
			texFlag = 1
		}
	}
	return texFlag
}
