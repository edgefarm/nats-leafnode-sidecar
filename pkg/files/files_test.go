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
