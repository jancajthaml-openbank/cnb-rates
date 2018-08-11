// Copyright (c) 2016-2018, Jan Cajthaml <jan.cajthaml@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package utils

import (
	"bytes"
	"io"
	"os"
	"reflect"
	"syscall"
	"unsafe"
)

var underlyingBufferSize int

func init() {
	underlyingBufferSize = 2 * os.Getpagesize()
}

// ReadFileFully reads whole file given absolute path
func ReadFileFully(absPath string) ([]byte, bool) {
	f, err := os.OpenFile(absPath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, false
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return nil, false
	}

	buf := make([]byte, fi.Size())
	_, err = f.Read(buf)
	if err != nil && err != io.EOF {
		return nil, false
	}

	return buf, true
}

// Exists returns true if absolute path exists
func Exists(absPath string) bool {
	_, err := os.Stat(absPath)
	return !os.IsNotExist(err)
}

// UpdateFile rewrite file with data given absolute path to a file if that file
// exist
func UpdateFile(absPath string, data io.Reader) bool {
	f, err := os.OpenFile(absPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return false
	}
	defer f.Close()

	if _, err := io.Copy(f, data); err != nil {
		return false
	}

	return true
}

func nameFromDirent(de *syscall.Dirent) []byte {
	ml := int(uint64(de.Reclen) - uint64(unsafe.Offsetof(syscall.Dirent{}.Name)))

	var name []byte
	sh := (*reflect.SliceHeader)(unsafe.Pointer(&name))
	sh.Cap = ml
	sh.Len = ml
	sh.Data = uintptr(unsafe.Pointer(&de.Name[0]))

	if index := bytes.IndexByte(name, 0); index >= 0 {
		sh.Cap = index
		sh.Len = index
	}

	return name
}

// ListDirectory returns slice of item names in given absolute path
func ListDirectory(absPath string) []string {
	v := make([]string, 0)

	dh, err := os.Open(absPath)
	if err != nil {
		return nil
	}

	fd := int(dh.Fd())

	scratchBuffer := make([]byte, underlyingBufferSize)

	var de *syscall.Dirent

	for {
		n, err := syscall.ReadDirent(fd, scratchBuffer)
		if err != nil {
			_ = dh.Close()
			return nil
		}
		if n <= 0 {
			break
		}
		buf := scratchBuffer[:n]
		for len(buf) > 0 {
			de = (*syscall.Dirent)(unsafe.Pointer(&buf[0]))
			buf = buf[de.Reclen:]

			if de.Ino == 0 {
				continue
			}

			nameSlice := nameFromDirent(de)
			namlen := len(nameSlice)
			if (namlen == 0) || (namlen == 1 && nameSlice[0] == '.') || (namlen == 2 && nameSlice[0] == '.' && nameSlice[1] == '.') {
				continue
			}
			v = append(v, string(nameSlice))
		}
	}

	if err = dh.Close(); err != nil {
		return nil
	}

	return v
}
