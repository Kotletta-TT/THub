package main

import (
	"bufio"
	"bytes"
	"fmt"
)

const (
	new        byte = 0xc8
	resize     byte = 0xc9
	disconnect byte = 0xca
	err        byte = 0xcb
	input      byte = 0xcc
	output     byte = 0xcd
)

var data []byte = []byte{0xc8, 0x00, 0x00, 0x01, 'h'}
var twoData []byte = []byte{0xc8, 0x00, 0x00, 0x01, 'h', 0xc8, 0x00, 0x00, 0x01, 'e'}
var inputData []byte = []byte{0xc8, 0x00, 0x00, 0x0a, 'e', 'c', 'h', 'o', ' ', 'h', 'e', 'l', 'l', 'o'}

func main() {
	// parse()
	scan()
}

func parse() {
	r := bytes.NewReader(inputData)
	b := bufio.NewReader(r)
	f, err := Parse(b)
	if err != nil {
		fmt.Println(err)
		return
	}
	// f2, err := Parse(b)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	fmt.Printf("f.Type: %d\nf.Length: %d\nf.TermId: %d\nf.Data: %v\n", f.Type, f.Length, f.TermId, f.Data[:f.Length])
	// fmt.Printf("f.Type: %d\nf.Length: %d\nf.TermId: %d\nf.Data: %v\n", f2.Type, f2.Length, f2.TermId, f2.Data[:f.Length])
}

type Frame struct {
	Type   byte
	Length uint16
	TermId byte
	Data   []byte
}

func Parse(br *bufio.Reader) (Frame, error) {
	var err error
	frame := Frame{Data: make([]byte, 65536)}
	frameLen := make([]byte, 2)
	frame.Type, err = br.ReadByte()
	if err != nil {
		return frame, fmt.Errorf("Parse ReadByte err: %s", err)
	}
	if frame.Type < 0xc8 || frame.Type > 0xcd {
		return frame, fmt.Errorf("Parse Type err: Unknown Type package")
	}
	frame.TermId, err = br.ReadByte()
	if err != nil {
		return frame, fmt.Errorf("Parse Read TermId err: %s", err)
	}
	n, err := br.Read(frameLen)
	if err != nil {
		return frame, fmt.Errorf("Parse Read Len err: %s", err)
	}
	if n != 2 {
		return frame, fmt.Errorf("Parse Read Len err: Len ==%d", n)
	}
	frame.Length = uint16(frameLen[0])<<8 | uint16(frameLen[1])

	n, err = br.Read(frame.Data[:frame.Length])
	if err != nil {
		return frame, fmt.Errorf("Parse Read Data err: %s", err)
	}
	if n != int(frame.Length) {
		return frame, fmt.Errorf("Parse Read Data err: reading %d len %d", n, frame.Length)
	}
	return frame, nil
}

func ScanFrames(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.IndexByte(data, 0xc7); i >= 0 {
		// We have a full newline-terminated line.
		return i + 1, data[0:i], nil
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), data, nil
	}
	// Request more data.
	return 0, nil, nil
}

var twoDatawithHeader []byte = []byte{0xc7, 0xc8, 0x00, 0x00, 0x01, 'h', 0xc7, 0xc8, 0x00, 0x00, 0x01, 'e'}
var twoDatawithHeaderwithTrash []byte = []byte{0xc7, 0xc8, 0x00, 0x00, 0x01, 'h', 0xc7, 0xc8, 0x00, 0x00, 0x01, 'e', 0x00, 0x01, 0x02, 0x03}

func scan() {
	rb := bytes.NewReader(twoDatawithHeader)
	s := bufio.NewScanner(rb)
	s.Split(ScanFrames)
	for s.Scan() {
		f := Frame{}
		b := s.Bytes()
		if len(b) < 4 {
			continue
		}
		f.Type = b[0]
		f.TermId = b[1]
		f.Length = uint16(b[2])<<8 | uint16(b[3])
		if len(b) > 4 {
			f.Data = b[4 : 4+f.Length]
		}
		fmt.Printf("f.Type: %d\nf.Length: %d\nf.TermId: %d\nf.Data: %v\n", f.Type, f.Length, f.TermId, f.Data[:f.Length])
	}
}
