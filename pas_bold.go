package main

var pasTable = []string{
	"absolute", "and", "array",
	"begin",
	"case", "const",
	"define", "div", "downto", "do",
	"else", "end", "external",
	"file", "forward", "for", "function",
	"goto",
	"if", "inline", "in",
	"label", "let",
	"mod",
	"nil", "not",
	"of", "or", "out",
	"packed", "procedure", "program",
	"record", "repeat", "return",
	"set", "shl", "shr", "string",
	"then", "to", "type",
	"until",
	"var",
	"while", "with",
	"xor",
}

func CompareWord(bufPtr int, tblStr string) int {
	ptr := DecBufPtr(bufPtr)
	c := byte(Buffer[ptr])

	if (c >= '0' && c <= '9') || (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') || c == '_' {
		return -1
	}

	ptr = bufPtr
	flag := 0
	for i := 0; i < len(tblStr); i++ {
		c = byte(Buffer[ptr])
		if c >= 'A' && c <= 'Z' {
			c += 0x20 // to lower
		}
		flag = int(c) - int(tblStr[i])
		if flag != 0 {
			break
		}
		ptr = IncBufPtr(ptr)
	}

	c = byte(Buffer[ptr])
	if (c >= '0' && c <= '9') || (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') || c == '_' {
		return -1
	}

	return flag
}

func GetBfFlag(bufPtr int) int {
	i := 0
	for _, word := range pasTable {
		if CompareWord(bufPtr, word) == 0 {
			i = len(word)
			break
		}
	}
	return i
}
