// Copyright (c) 2015 Uber Technologies, Inc.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"reflect"
	"unsafe"
)

// embeddedThriftGZip is the contents of a gzipped binary that can be used instead of the
// thrift binary on the system. This is used to tie thrift-gen to a specific version of
// the thrift compiler when generating code.
// This variable can be set by overlaying the thrift-gen folder with an additional file
// that sets this variable in an init function.
var embeddedThriftGZip string

func execCmd(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// readOnlyByteSlice is an unsafe function that returns a byte slice
// that represents the read-only data for a string.
func readOnlyByteSlice(s *string) []byte {
	sx := (*reflect.StringHeader)(unsafe.Pointer(s))
	b := make([]byte, 0, 0)
	bx := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	bx.Data = sx.Data
	bx.Len = len(*s)
	bx.Cap = bx.Len
	return b
}

func execThrift(args ...string) error {
	// If we do not have an embedded Thrift binary for this OS, use the system binary.
	if len(embeddedThriftGZip) == 0 {
		return execCmd("thrift", args...)
	}

	// Create a temporary file to write out the uncompressed binary.
	f, err := ioutil.TempFile("", "thrift")
	if err != nil {
		return fmt.Errorf("failed to get temp file: %v", err)
	}

	// Decompress the file and write it out.
	reader, err := gzip.NewReader(bytes.NewReader(readOnlyByteSlice(&embeddedThriftGZip)))
	if err != nil {
		return fmt.Errorf("gzip.NewReader failed: %v", err)
	}
	if _, err := io.Copy(f, reader); err != nil {
		return fmt.Errorf("Write out binary failed: %v", err)
	}
	f.Close()

	// Temp file is not created with +x, so Chmod it.
	if err := os.Chmod(f.Name(), 0700); err != nil {
		return fmt.Errorf("failed to add execute permissions to temp file: %v", err)
	}

	return execCmd(f.Name(), args...)
}
