package registry

import (
	"bytes"
	"encoding/json"
	"fmt"

	jsonpatch "github.com/evanphx/json-patch"
)

const (
	ngsHost       = "tls://connect.ngs.global:7422"
	defaultConfig = `{
	"http": 8222,
	"leafnodes": {
		"remotes": []
	},
	"pid_file": "/var/run/nats/nats.pid"
}`
)

func jsonPrettyPrint(in string) string {
	var out bytes.Buffer
	err := json.Indent(&out, []byte(in), "", "\t")
	if err != nil {
		return in
	}
	return out.String()
}

func (r *Registry) addCredentials(network string, creds string) error {
	raw, err := decodeRawJSON(r.configFileContent)
	if err != nil {
		return err
	}

	// check if remote already exists
	remoteFound := false
	for _, r := range raw["leafnodes"].(map[string]interface{})["remotes"].([]interface{}) {
		remote := r.(map[string]interface{})
		if ok := remote["account"] == network; ok {
			remoteFound = true
			break
		}
	}

	modified, err := json.Marshal(raw)
	if err != nil {
		return err
	}
	if !remoteFound {
		patchJSON := []byte(fmt.Sprintf(`[{"op": "add", "path": "/leafnodes/remotes/-", "value": {"url": "%s", "credentials": "%s/%s.creds"}}]`, ngsHost, r.credsFilesPath, network))
		patch, err := jsonpatch.DecodePatch(patchJSON)
		if err != nil {
			return err
		}

		modified, err = patch.Apply([]byte(r.configFileContent))
		if err != nil {
			return err
		}
	}

	r.configFileContent = string(modified)
	return nil
}

func (r *Registry) removeCredentials(account string) error {
	raw, err := decodeRawJSON(r.configFileContent)
	if err != nil {
		return err
	}
	accountIndex := 0
	for k := range raw["accounts"].(map[string]interface{}) {
		if k == account {
			delete(raw["accounts"].(map[string]interface{}), k)
			break
		}
		accountIndex++
	}

	//  check if remote is already existing
	remoteFound := false
	for _, remote := range raw["leafnodes"].(map[string]interface{})["remotes"].([]interface{}) {
		if remote.(map[string]interface{})["account"] == account {
			remoteFound = true
		}
	}
	if !remoteFound {
		return fmt.Errorf("remote for account %s not found", account)
	}
	newRemotes := []interface{}{}
	for _, remote := range raw["leafnodes"].(map[string]interface{})["remotes"].([]interface{}) {
		if remote.(map[string]interface{})["account"] != account {
			newRemotes = append(newRemotes, remote.(map[string]interface{}))
			break
		}
	}
	if len(newRemotes) == len(raw["leafnodes"].(map[string]interface{})["remotes"].([]interface{})) {
		return fmt.Errorf("account %s not found", account)
	}

	raw["leafnodes"].(map[string]interface{})["remotes"] = newRemotes
	config, err := json.Marshal(raw)
	if err != nil {
		return err
	}
	r.configFileContent = string(config)
	return nil
}

// Dump prints the registry configuration as pretty formatted JSON
func (r *Registry) Dump() {
	fmt.Println(jsonPrettyPrint(r.configFileContent))
}

func decodeRawJSON(jsonStr string) (map[string]interface{}, error) {
	jsonMap := make(map[string]interface{})
	err := json.Unmarshal([]byte(jsonStr), &jsonMap)
	if err != nil {
		return nil, err
	}
	return jsonMap, nil
}
