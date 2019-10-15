// Copyright © 2015 Jerry Jacobs <jerry.jacobs@xor-gate.org>.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ftpfs

import (
	"os"

	"github.com/jlaffaye/ftp"
)

type fileMode int

const (
	modeClosed fileMode = iota
	modeOpenRead
	modeOpenWrite
)

type File struct {
	writePipe    *os.File
	readResponse *ftp.Response
	mode         fileMode
	name         string
	conn         *ftp.ServerConn
	fs           *Fs
}

func FileOpen(fs *Fs, name string) (*File, error) {
	return &File{fs: fs, conn: fs.client, name: name}, nil
}

func FileCreate(fs *Fs, name string) (*File, error) {
	return &File{fs: fs, conn: fs.client, name: name}, nil
}

func (f *File) Close() error {
	if f.writePipe != nil {
		return f.writePipe.Close()
	}

	if f.readResponse != nil {
		return f.readResponse.Close()
	}

	return nil
}

func (f *File) Name() string {
	return f.name
}

func (f *File) Stat() (os.FileInfo, error) {
	return f.fs.Stat(f.name)
}

func (f *File) Sync() error {
	return nil
}

func (f *File) Truncate(size int64) error {
	return nil
}

func (f *File) Read(b []byte) (n int, err error) {
	if f.readResponse == nil {
		f.readResponse, err = f.conn.Retr(f.name)
		if err != nil {
			return 0, nil
		}

	}

	return f.readResponse.Read(b)
}

// TODO
func (f *File) ReadAt(b []byte, off int64) (n int, err error) {
	return 0, nil
}

func (f *File) Readdir(count int) (res []os.FileInfo, err error) {
	entries, err := f.conn.List(f.name)
	if err != nil {
		return nil, err
	}

	res = []os.FileInfo{}
	for _, e := range entries {
		res = append(res, FileInfo{
			name:    e.Name,
			size:    int64(e.Size),
			isdir:   e.Type == ftp.EntryTypeFolder,
			sys:     e,
			modtime: e.Time,
		})
	}
	return res, nil
}

func (f *File) Readdirnames(n int) (names []string, err error) {
	entries, err := f.conn.List(f.name)
	if err != nil {
		return nil, err
	}

	res := []string{}
	for _, e := range entries {
		res = append(res, e.Name)
	}
	return res, nil
}

func (f *File) Seek(offset int64, whence int) (int64, error) {
	return 0, nil
}

func (f *File) Write(b []byte) (n int, err error) {
	if f.writePipe == nil {
		r, w, err := os.Pipe()
		if err != nil {
			return 0, nil
		}
		go func() {
			f.conn.StorFrom(f.name, r, 0)
			f.writePipe = nil
		}()

		f.writePipe = w

	}
	return f.writePipe.Write(b)
}

// TODO
func (f *File) WriteAt(b []byte, off int64) (n int, err error) {
	return 0, nil
}

func (f *File) WriteString(s string) (ret int, err error) {
	return f.Write([]byte(s))
}
