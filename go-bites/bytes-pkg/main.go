package main

import (
	"bytes"
	"fmt"
	"os"
)

func main() {
	buffer := bytes.NewBuffer([]byte("Hello, World!"))
	buffer.Write([]byte(" How are you?"))
	buffer.WriteByte('c')
	buffer.WriteString(" I am fine.")
	buffer.WriteRune('😊')

	println(buffer.String())

	rdbuffer := make([]byte, 13)
	buffer.Read(rdbuffer)

	println(string(rdbuffer))

	var buffer2 bytes.Buffer
	// here we can see that buffer2 is empty, so it has no bytes in it
	println("bytes in buffer ", buffer2.Len())
	println("capacity of buffer ", buffer2.Cap())
	// If we try writing to zero capacity buffer, it will internally grow the buffer
	// to accomodate the bytes, here it grows by 64 bytes
	buffer2.Write([]byte(" How are you?"))

	println("bytes in buffer ", buffer2.Len())
	println("capacity of buffer ", buffer2.Cap())
	println("unused bytes in buffer ", buffer2.Available())

	var buffer3 bytes.Buffer
	file, err := os.Open("hello.txt")
	if err != nil {
		println("Error opening file: ", err)
	}
	n, err := buffer3.ReadFrom(file)
	if err != nil {
		println("Error reading from buffer: ", err)
	}
	fmt.Printf("bytes read from buffer %d, text read from file %s", n, buffer3.String())

	buffer3.WriteTo(os.Stdout)

}
