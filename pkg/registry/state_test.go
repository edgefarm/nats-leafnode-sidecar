package registry

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStateNewState(t *testing.T) {
	assert := assert.New(t)

	file, err := ioutil.TempFile("", "state-")
	assert.Nil(err)
	defer os.Remove(file.Name())

	state := NewState(file.Name())
	state.Current.NetworkUsage = map[string][]string{
		"foo": {"a"},
		"bar": {"a", "b"},
	}

	err = state.Save()
	assert.Nil(err)

	newState := NewState(file.Name())
	err = newState.Read()
	assert.Nil(err)

	assert.Equal(newState.Current.NetworkUsage, map[string][]string{
		"foo": {"a"},
		"bar": {"a", "b"},
	})
}

// Removed test, this won't run properly using github actions
// func TestStateNewStateForbiden(t *testing.T) {
// 	assert := assert.New(t)
// 	path := fmt.Sprintf("/tmp/%s", uuid.New().String())
// 	err := os.Mkdir(path, 0555)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println(path)
// 	// No need to check whether `recover()` is nil. Just turn off the panic.
// 	defer func() { recover() }()

// 	defer os.RemoveAll(path)
// 	file := fmt.Sprintf("%s/state.json", path)
// 	assert.Nil(err)
// 	defer os.Remove(file)

// 	// should panic, because the directory is not writable
// 	NewState(file)

// 	// Never reaches here if `OtherFunctionThatPanics` panics.
// 	t.Errorf("did not panic")
// }

func TestStateStateIncrement(t *testing.T) {
	assert := assert.New(t)
	file, err := ioutil.TempFile("", "state-")
	assert.Nil(err)
	defer os.Remove(file.Name())

	state := NewState(file.Name())
	state.Current.NetworkUsage = map[string][]string{
		"foo": {"a"},
		"bar": {"a", "b"},
	}

	// add participant to foo
	err = state.Update("foo", "z", Add)
	assert.Nil(err)
	foo, err := state.Usage("foo")
	assert.Nil(err)
	assert.Equal(foo, 2)

	// again add the same participant to foo. Should be ignored.
	err = state.Update("foo", "z", Add)
	assert.Nil(err)
	foo, err = state.Usage("foo")
	assert.Nil(err)
	assert.Equal(foo, 2)

	// not existent network
	foobar, err := state.Usage("foobar")
	assert.NotNil(err)
	assert.Equal(foobar, 0)

	// remove participant from foo
	err = state.Update("foo", "z", Remove)
	assert.Nil(err)
	foo, err = state.Usage("foo")
	assert.Nil(err)
	assert.Equal(foo, 1)

	// foo cannot be deleted. one remaining participant
	canDelete, err := state.CanDelete("foo")
	assert.Nil(err)
	assert.False(canDelete)

	// try to delete foo
	err = state.Delete("foo")
	assert.NotNil(err)

	// remove participant from foo
	err = state.Update("foo", "a", Remove)
	assert.Nil(err)

	// foo now can be deleted
	canDelete, err = state.CanDelete("foo")
	assert.Nil(err)
	assert.True(canDelete)

	// delete foo
	err = state.Delete("foo")
	assert.Nil(err)

	// foo is not in the state anymore
	err = state.Update("foo", "y", Remove)
	assert.NotNil(err)
	_, err = state.Usage("foo")
	assert.NotNil(err)

	// foo is not existent anymore
	_, err = state.CanDelete("foo")
	assert.NotNil(err)
}
