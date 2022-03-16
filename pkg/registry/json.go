package registry

import (
	"bytes"
	"encoding/json"
)

type Remote struct {
	Url         string `json:"url"`
	Credentials string `json:"credentials"`
}

type Leafnodes struct {
	Remotes []Remote `json:"remotes"`
}

type NatsConfig struct {
	PidFile   string    `json:"pid_file"`
	Http      int       `json:"http"`
	Leafnodes Leafnodes `json:"leafnodes"`
}

func NewJson(path string) *NatsConfig {
	// Load config file if it exists
	var config NatsConfig
	str, err := readFile(path)
	if err == nil {
		err = json.Unmarshal([]byte(str), &config)
		if err != nil {
			panic(err)
		}
	} else {
		config = NatsConfig{
			PidFile: "/var/run/nats/nats.pid",
			Http:    8222,
			Leafnodes: Leafnodes{
				Remotes: []Remote{},
			},
		}
	}
	return &config
}

func (c *NatsConfig) Dump() string {
	b, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(b)
}

func jsonPrettyPrint(in string) string {
	var out bytes.Buffer
	err := json.Indent(&out, []byte(in), "", "\t")
	if err != nil {
		return in
	}
	return out.String()
}
