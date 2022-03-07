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

	if err := state.Read(); err != nil {
		fmt.Println("State file not found, creating empty state")
		err = state.createEmpty()
		if err != nil {
			panic(err)
		}
	}
	return state
}

// Save saves the state to the file
func (s *State) Save() error {
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

// Read reads the state from the file
func (s *State) Read() error {
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

const (
	// RegisterParticipant is the increment for a new participant
	RegisterParticipant int = 1
	// UnregisterParticipant is the decrement for a removed participant
	UnregisterParticipant int = -1
)

// Update updates the state. Takes either RegisterParticipant or UnregisterParticipant
func (s *State) Update(network string, increment int) error {
	if _, ok := s.Current.NetworkUsage[network]; !ok {
		return fmt.Errorf("network %s not found", network)
	}
	s.Current.NetworkUsage[network] += increment
	return s.Save()
}

// Usage returns the usage count of the network
func (s *State) Usage(network string) (int, error) {
	if _, ok := s.Current.NetworkUsage[network]; !ok {
		return 0, fmt.Errorf("network %s not found", network)
	}
	return s.Current.NetworkUsage[network], nil
}

func (s *State) createEmpty() error {
	return s.Save()
}

// CanDelete checks whether the network can be deleted
func (s *State) CanDelete(network string) (bool, error) {
	usage, err := s.Usage(network)
	if err != nil {
		return false, err
	}
	return usage <= 0, nil
}

// Delete deletes the network from the state
func (s *State) Delete(network string) error {
	if _, ok := s.Current.NetworkUsage[network]; !ok {
		return fmt.Errorf("network %s not found", network)
	}
	usage, err := s.Usage(network)
	if err != nil {
		return err
	}
	if usage > 0 {
		return fmt.Errorf("network %s is still in use", network)
	}
	delete(s.Current.NetworkUsage, network)
	return s.Save()
}
