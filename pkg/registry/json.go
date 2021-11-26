package registry

import (
	"bytes"
	"encoding/json"
	"fmt"

	jsonpatch "github.com/evanphx/json-patch"
)

const (
	ngsHost = "tls://connect.ngs.global:7422"
)

var config = `{
	"pid_file": "/var/run/nats.pid",
	"http": 8222,
	"server_name": "edge",
	"leafnodes": {
		"remotes": []
	},
	"accounts": {
	}
}`

func jsonPrettyPrint(in string) string {
	var out bytes.Buffer
	err := json.Indent(&out, []byte(in), "", "\t")
	if err != nil {
		return in
	}
	return out.String()
}

func (r *Registry) addCredentials(account string, user string, password string, creds string) error {
	patchJSON := []byte(fmt.Sprintf(`[{"op": "add", "path": "/leafnodes/remotes/-", "value": {"url": "%s", "credentials": "%s.creds", "account": "%s"}}]`, ngsHost, user, account))
	patch, err := jsonpatch.DecodePatch(patchJSON)
	if err != nil {
		return err
	}

	modified, err := patch.Apply([]byte(r.configFileContent))
	if err != nil {
		return err
	}
	raw, err := decodeRawJSON(r.configFileContent)
	if err != nil {
		return err
	}

	//  check if accounts is already existing
	accountFound := false
	for k := range raw["accounts"].(map[string]interface{}) {
		if k == account {
			accountFound = true
			break
		}
	}
	if accountFound {
		return fmt.Errorf("account %s already exists", account)
	}

	accounts := raw["accounts"].(map[string]interface{})
	accounts[account] = map[string]interface{}{
		"users": map[string]interface{}{
			"user":     user,
			"password": password,
		},
	}
	accountsString, err := json.Marshal(accounts)
	if err != nil {
		return err
	}
	patchJSON = []byte(fmt.Sprintf(`[{"op": "replace", "path": "/accounts", "value": %s}]`, accountsString))
	patch, err = jsonpatch.DecodePatch(patchJSON)
	if err != nil {
		return err
	}

	modified, err = patch.Apply(modified)
	if err != nil {
		return err
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
