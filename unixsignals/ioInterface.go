package main

import (
	"bufio"
	"fmt"
	"io"
)

type S1 struct {
	F1 int
	F2 string
}

type S2 struct {
	F1   S1
	text []byte
}

// note here, the original p []byte won't be updated because everyhting in go is pass by value
//
//	case *[]byte:
//	We scan to string and convert so we get a copy of the data.
//	If we scanned to bytes, the slice would point at the buffer.
//	*v = []byte(s.convertString(verb))
//
// the original slice headerPtr, len and cap is not updated, instead the value is changed to a new obj.
// Change the function to accept a pointer if you want to modify the original slice.
func (s *S1) Read(p []byte) (n int, err error) {
	fmt.Println("Give me you name: ")
	fmt.Scanln(&p)
	s.F2 = string(p)
	return len(p), nil
}

func (s *S1) Write(p []byte) (n int, err error) {
	if s.F1 < 0 {
		return -1, nil
	}

	for i := 0; i < s.F1; i++ {
		fmt.Printf("%s ", p)
	}
	fmt.Println()
	return s.F1, nil
}

func (s S2) eof() bool {
	return len(s.text) == 0
}

// readByte function reads a single byte from the buffer at a time.
// This functions assumes that eof() check has been done.
func (s *S2) readByte() byte {
	temp := s.text[0]
	s.text = s.text[1:]
	return temp
}

func (s *S2) Read(p []byte) (n int, err error) {
	if s.eof() {
		err = io.EOF
		return
	}

	l := len(p)
	if l > 0 {
		for n < l {
			p[n] = s.readByte()
			n++
			if s.eof() {
				s.text = s.text[0:0]
				break
			}
		}
	}
	return
}

func main2() {
	s1var := S1{
		F1: 4,
		F2: "Hello",
	}
	fmt.Println(s1var)
	buf := make([]byte, 2)
	_, err := s1var.Read(buf)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Read:", s1var.F2)
	_, _ = s1var.Write([]byte("Hello Worlds"))

	s2var := S2{
		F1:   s1var,
		text: []byte("Hello World"),
	}
	// Read s2var.text
	r := bufio.NewReader(&s2var)

	for {
		n, err := r.Read(buf)
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("*", err)
			break
		}
		fmt.Println("**", n, "[", string(buf[:n]), "]")
	}
}
