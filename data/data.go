package data

import json "github.com/json-iterator/go"

type JsonMap map[string]interface{}
type Record map[string]string

func (data JsonMap) String() string {
	byte, err := json.Marshal(&data)
	if err != nil {
		panic(err)
	}
	return string(byte)
}

func ParseJsonBytesToMap(data []byte) (JsonMap, error) {
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
	return ParseJsonBytesToMap(jsonBytes)
}
