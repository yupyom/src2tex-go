package main

import (
	"io"
	"log"
	"os"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

var statFlagGetData = 0
var bufPtr1 = 0
var bufPtr2 = BufferSize / 2
var readerWrapped io.Reader

// IncBufPtr increases buffer index circularly.
func IncBufPtr(idx int) int {
	if idx < 0 || idx >= BufferSize {
		log.Fatal("\nError: buffering error !!\n")
	}
	idx++
	if idx == BufferSize {
		idx = 0
	}
	return idx
}

// DecBufPtr decreases buffer index circularly.
func DecBufPtr(idx int) int {
	if idx < 0 || idx >= BufferSize {
		log.Fatal("\nError: buffering error !!\n")
	}
	idx--
	if idx < 0 {
		idx = BufferSize - 1
	}
	return idx
}

// Fgetc2Buffer reads a character into the circular buffer and advances pointers.
// In case of first call, it fills half of the buffer.
// It returns the index pointing to the current character being processed.
func Fgetc2Buffer(file *os.File) int {
	if statFlagGetData == 0 {
		statFlagGetData++

		// Initialize wrapped reader if needed based on InputEncoding
		if InputEncoding == "sjis" {
			readerWrapped = transform.NewReader(file, japanese.ShiftJIS.NewDecoder())
		} else if InputEncoding == "euc" {
			readerWrapped = transform.NewReader(file, japanese.EUCJP.NewDecoder())
		} else {
			readerWrapped = file
		}

		// Fill the first half of the buffer with input data
		for i := 0; i < BufferSize/2; i++ {
			b := make([]byte, 1)
			n, err := readerWrapped.Read(b)
			if err != nil || n == 0 {
				Buffer[i] = -1 // EOF
			} else {
				Buffer[i] = int(b[0])
			}
		}
		// Clear the second half (history)
		for i := BufferSize / 2; i < BufferSize; i++ {
			Buffer[i] = 0x00
		}
		bufPtr1 = 0
		bufPtr2 = BufferSize / 2
		return bufPtr1
	} else {
		// Read a datum from the file to Buffer[bufPtr2]
		b := make([]byte, 1)
		n, err := readerWrapped.Read(b)
		if err != nil || n == 0 {
			Buffer[bufPtr2] = -1 // EOF
		} else {
			Buffer[bufPtr2] = int(b[0])
		}
		bufPtr2 = IncBufPtr(bufPtr2)
		bufPtr1 = IncBufPtr(bufPtr1)
		return bufPtr1
	}
}
