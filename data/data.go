package data

import "encoding/json"

type JsonMap map[string]interface{}
type Record map[string]string

func parseJsonBytesToMap(data []byte) (JsonMap, error) {
	var dataMap map[string]interface{}
	err := json.Unmarshal(data, &dataMap)
	if err != nil {
		return nil, err
	}
	return dataMap, nil
}

func ParseToJsonMap(obj interface{}) (JsonMap, error) {
	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	return parseJsonBytesToMap(jsonBytes)
}
