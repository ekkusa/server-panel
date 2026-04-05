package main

import (
	"io"
)

// readDockerStream reads a Docker multiplexed log stream (8-byte header per frame)
// and writes the raw payload bytes to w.
func readDockerStream(r io.Reader, w io.Writer) {
	header := make([]byte, 8)
	buf := make([]byte, 4096)
	for {
		if _, err := io.ReadFull(r, header); err != nil {
			return
		}
		size := int(header[4])<<24 | int(header[5])<<16 | int(header[6])<<8 | int(header[7])
		if size == 0 {
			continue
		}
		if size > len(buf) {
			buf = make([]byte, size)
		}
		if _, err := io.ReadFull(r, buf[:size]); err != nil {
			return
		}
		w.Write(buf[:size])
	}
}
