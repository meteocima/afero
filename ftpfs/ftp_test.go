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
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFtpConnect(t *testing.T) {
	/*err := os.RemoveAll("/home/test/foo")
	assert.NoError(t, err)
	err = os.RemoveAll("/home/test/bar")
	assert.NoError(t, err)
	*/
	Fs, err := FtpConnect("test", "test", "localhost:21")
	if err != nil {
		t.Fatal(err)
	}
	defer Fs.Disconnect()

	err = Fs.Mkdir("foo", os.FileMode(0000))
	assert.NoError(t, err)

	err = Fs.Chmod("foo", os.FileMode(0700))
	assert.NoError(t, err)

	err = Fs.Mkdir("bar", os.FileMode(0777))
	assert.NoError(t, err)

	info, err := Fs.Stat("bar")
	assert.NoError(t, err)
	if !info.IsDir() {
		assert.Fail(t, "dir not created")
	}

	Fs.MkdirAll("test/dir1/dir2/dir3", os.FileMode(0777))

	file, err := Fs.Create("file1")
	if err != nil {
		t.Error(err)
	}
	defer file.Close()

	file.Write([]byte("hello\t"))
	file.WriteString("world!\n")

	f1, err := Fs.Open("file1")
	if err != nil {
		log.Fatalf("open: %v", err)
	}
	defer f1.Close()

	b := make([]byte, 100)

	_, err = f1.Read(b)
	fmt.Println(string(b))

	// TODO check here if "hello\tworld\n" is in buffer b
}
