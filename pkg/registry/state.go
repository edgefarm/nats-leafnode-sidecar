package registry

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// State is a struct that holds the state of the registry
type State struct {
	// Path to the state file
	Path string
	// Current state
	Current current
}

// State is the current state of the registry
type current struct {
	// NetworkUsage counts the number of components that are using the network.
	// This is needed for a proper clean up.
	NetworkUsage map[string]int `json:"network_usage"`
}

// NewState creates a new state
func NewState(path string) *State {
	state := &State{
		Path: path,
		Current: current{
			NetworkUsage: map[string]int{},
		},
	}

	if err := state.ReadState(); err != nil {
		fmt.Println("State file not found, creating empty state")
		err = state.createEmptyState()
		if err != nil {
			panic(err)
		}
	}
	return state
}

// SaveState saves the state to the file
func (s *State) SaveState() error {
	str, err := json.Marshal(s.Current)
	if err != nil {
		return err
	}

	f, err := os.Create(s.Path)
	if err != nil {
		return err
	}
	_, err = f.Write(str)
	if err != nil {
		return err
	}
	return nil
}

// ReadState reads the state from the file
func (s *State) ReadState() error {
	data, err := ioutil.ReadFile(s.Path)
	if err != nil {
		fmt.Print(err)
	}
	err = json.Unmarshal(data, &s.Current)
	if err != nil {
		return err
	}
	return nil
}

// UpdateState updates the state
func (s *State) UpdateState(network string, increment int) error {
	s.Current.NetworkUsage[network] += increment
	return nil
}

func (s *State) createEmptyState() error {
	return s.SaveState()
}
