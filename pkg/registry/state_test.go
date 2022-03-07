package registry

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestStateNewState(t *testing.T) {
	assert := assert.New(t)

	file, err := ioutil.TempFile("", "state-")
	assert.Nil(err)
	defer os.Remove(file.Name())

	state := NewState(file.Name())
	state.Current.NetworkUsage = map[string]int{
		"foo": 1,
		"bar": 2,
	}

	err = state.SaveState()
	assert.Nil(err)

	newState := NewState(file.Name())
	err = newState.ReadState()
	assert.Nil(err)

	assert.Equal(newState.Current.NetworkUsage, map[string]int{
		"foo": 1,
		"bar": 2,
	})
}

func TestStateNewStateForbiden(t *testing.T) {
	assert := assert.New(t)
	path := fmt.Sprintf("/tmp/%s", uuid.New().String())
	err := os.Mkdir(path, 0555)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(path)
	// No need to check whether `recover()` is nil. Just turn off the panic.
	defer func() { recover() }()

	defer os.RemoveAll(path)
	file := fmt.Sprintf("%s/state.json", path)
	assert.Nil(err)
	defer os.Remove(file)

	// should panic, because the directory is not writable
	NewState(file)

	// Never reaches here if `OtherFunctionThatPanics` panics.
	t.Errorf("did not panic")
}
