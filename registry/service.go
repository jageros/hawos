package registry

import (
	"encoding/json"
)

func marshal(si *ServiceInstance) (string, error) {
	data, err := json.Marshal(si)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func unmarshal(data []byte) (si *ServiceInstance, err error) {
	err = json.Unmarshal(data, &si)
	return
}
