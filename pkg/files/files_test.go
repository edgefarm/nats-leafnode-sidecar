/*
Copyright © 2021 Ci4Rail GmbH <engineering@ci4rail.com>
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package files

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// GetFiles returns a list of files in the given directory.
func TestGetFiles(t *testing.T) {
	assert := assert.New(t)
	assert.Nil(nil)
	dir, err := ioutil.TempDir("/tmp", "")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(dir)
	defer os.RemoveAll(dir)

	for _, file := range []string{"file1", "file2", "file3"} {
		_, err := os.Create(fmt.Sprintf("%s/%s", dir, file))
		if err != nil {
			fmt.Println(err)
			return
		}
		defer os.Remove(fmt.Sprintf("%s/%s", dir, file))
	}

	err = os.Symlink(fmt.Sprintf("%s/file1", dir), fmt.Sprintf("%s/link", dir))
	assert.Nil(err)
	defer os.Remove(fmt.Sprintf("%s/link", dir))

	files, err := GetFiles(dir)
	assert.Nil(err)
	assert.Equal(4, len(files))
	assert.Contains(files, fmt.Sprintf("%s/file1", dir))
	assert.Contains(files, fmt.Sprintf("%s/file2", dir))
	assert.Contains(files, fmt.Sprintf("%s/file3", dir))
	assert.Contains(files, fmt.Sprintf("%s/link", dir))
}
