package registry

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/edgefarm/nats-leafnode-sidecar/pkg/unique"
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
	NetworkUsage map[string][]string `json:"network_usage"`
}

// NewState creates a new state
func NewState(path string) *State {
	state := &State{
		Path: path,
		Current: current{
			NetworkUsage: make(map[string][]string),
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

// UpdateAction is an enum for the update action
type UpdateAction string

const (
	// Add is the action for adding a component to the state
	Add UpdateAction = "add"
	// Remove is the action for removing a component from the state
	Remove UpdateAction = "remove"
)

// Update updates the state.
func (s *State) Update(network string, component string, action UpdateAction) error {
	if _, ok := s.Current.NetworkUsage[network]; !ok {
		fmt.Printf("network %s not found. Creating...", network)
	}
	if action == Add {
		// ignore multiple registrations for components
		s.Current.NetworkUsage[network] = unique.Slice(append(s.Current.NetworkUsage[network], component))
	} else {
		s.Current.NetworkUsage[network] = unique.Slice(remove(s.Current.NetworkUsage[network], component))
	}
	return s.Save()
}

func remove(slice []string, s string) []string {
	for i, v := range slice {
		if v == s {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

// Usage returns the usage count of the network
func (s *State) Usage(network string) (int, error) {
	if _, ok := s.Current.NetworkUsage[network]; !ok {
		return 0, fmt.Errorf("network %s not found", network)
	}
	return len(s.Current.NetworkUsage[network]), nil
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
