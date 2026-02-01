package common

import (
	"strings"

	"sigs.k8s.io/yaml"
)

type Values map[string]interface{}

func (v Values) YAML() (string, error) {
	b, err := yaml.Marshal(v)
	return string(b), err
}

func (v Values) AsMap() map[string]interface{} {
	if len(v) == 0 {
		return map[string]interface{}{}
	}
	return v
}

func ReadValues(data []byte) (vals Values, err error) {
	err = yaml.Unmarshal(data, &vals)
	if len(vals) == 0 {
		vals = Values{}
	}
	return vals, err
}

func Pick(values Values, path string) interface{} {
	var temp = values
	var parts = strings.Split(path, ".")

	for index, part := range parts {
		if value, ok := temp[part].(map[string]interface{}); ok {
			temp = value
		} else {
			if index == (len(parts) - 1) {
				return temp[part]
			} else {
				return nil
			}
		}
	}

	return temp
}
