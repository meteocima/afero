package ftpfs

import (
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"github.com/jlaffaye/ftp"
	"github.com/parro-it/afero"
)

// Fs is a afero.Fs implementation that uses functions provided by the sftp package.
//
// For details in any method, check the documentation of the sftp package
// (github.com/pkg/sftp).
type Fs struct {
	client *ftp.ServerConn
}

func FtpConnect(user, password, host string) (*Fs, error) {
	client, err := ftp.Dial(host, ftp.DialWithTimeout(25*time.Second))
	if err != nil {
		return nil, err
	}
	err = client.Login(user, password)
	if err != nil {
		return nil, err
	}
	err = client.ChangeDir("~")
	if err != nil {
		return nil, err
	}

	dir, err := client.CurrentDir()
	if err != nil {
		return nil, err
	}
	fmt.Println(dir)
	return &Fs{client: client}, nil
}

func New(client *ftp.ServerConn) afero.Fs {
	return &Fs{client: client}
}

func (s Fs) Disconnect() {
	if err := s.client.Logout(); err != nil {
		log.Printf("ERROR: Logout from remote server: %v\n", err)
	}
	if err := s.client.Quit(); err != nil {
		log.Printf("ERROR: Quit from remote server: %v\n", err)
	}
}

func (s Fs) Name() string { return "ftpfs" }

func (s *Fs) Create(name string) (afero.File, error) {
	return FileCreate(s, name)
}

func (s Fs) Mkdir(name string, perm os.FileMode) error {
	return s.client.MakeDir(name)
}

func (s Fs) MkdirAll(path string, perm os.FileMode) error {

	/*// Fast path: if we can tell whether path is a directory or file, stop with success or error.
	dir, err := s.Stat(path)
	if err == nil {
		if dir.IsDir() {
			return nil
		}
		return err
	}

	// Slow path: make sure parent exists and then call Mkdir for path.
	i := len(path)
	for i > 0 && os.IsPathSeparator(path[i-1]) { // Skip trailing path separator.
		i--
	}

	j := i
	for j > 0 && !os.IsPathSeparator(path[j-1]) { // Scan backward over element.
		j--
	}

	if j > 1 {
		// Create parent
		err = s.MkdirAll(path[0:j-1], perm)
		if err != nil {
			return err
		}
	}

	// Parent now exists; invoke Mkdir and use its result.
	err = s.Mkdir(path, perm)
	if err != nil {
		// Handle arguments like "foo/." by
		// double-checking that directory doesn't exist.
		dir, err1 := s.Lstat(path)
		if err1 == nil && dir.IsDir() {
			return nil
		}
		return err
	}
	return nil*/
	return nil
}

func (s *Fs) Open(name string) (afero.File, error) {
	return FileOpen(s, name)
}

func (s Fs) OpenFile(name string, flag int, perm os.FileMode) (afero.File, error) {
	return nil, nil
}

func (s Fs) Remove(name string) error {
	return nil // s.client.Remove(name)
}

func (s Fs) RemoveAll(path string) error {
	// TODO have a look at os.RemoveAll
	// https://github.com/golang/go/blob/master/src/os/path.go#L66
	return nil
}

func (s Fs) Rename(oldname, newname string) error {
	return s.client.Rename(oldname, newname)
}

type FileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modtime time.Time
	isdir   bool
	sys     interface{}
}

// base name of the file
func (info FileInfo) Name() string {
	return info.name
}

// length in bytes for regular files; system-dependent for others
func (info FileInfo) Size() int64 {
	return info.size
}

// file mode bits
func (info FileInfo) Mode() os.FileMode {
	return info.mode
}

// modification time
func (info FileInfo) ModTime() time.Time {
	return info.modtime
}

// abbreviation for Mode().IsDir()
func (info FileInfo) IsDir() bool {
	return info.isdir
}

// underlying data source (can return nil)
func (info FileInfo) Sys() interface{} {
	return info.sys
}

func (s Fs) Stat(name string) (os.FileInfo, error) {

	dir, basename := path.Split(name)

	entries, err := s.client.List(dir)
	if err != nil {
		return nil, err
	}
	if len(entries) == 0 {
		return nil, os.ErrNotExist
	}

	for _, e := range entries {
		if e.Name != basename {
			continue
		}

		return &FileInfo{
			name:    e.Name,
			size:    int64(e.Size),
			isdir:   e.Type == ftp.EntryTypeFolder,
			sys:     e,
			modtime: e.Time,
		}, nil //s.client.Stat(name)
	}

	return nil, os.ErrNotExist
}

func (s Fs) Lstat(p string) (os.FileInfo, error) {
	return nil, nil //s.client.Lstat(p)
}

func (s Fs) Chmod(name string, mode os.FileMode) error {
	return nil //s.client.Chmod(name, mode)
}

func (s Fs) Chtimes(name string, atime time.Time, mtime time.Time) error {
	return nil //s.client.Chtimes(name, atime, mtime)
}

func (s Fs) Link(name, targetDir string) error {
	panic("not implemented")
}
