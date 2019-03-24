package data

import json "github.com/json-iterator/go"

type JsonMap map[string]interface{}
type Record map[string]string

type LookupMap map[string]bool

func (mp LookupMap) exists(val string) bool {
	_, ok := mp[val]
	return ok
}

func (mp LookupMap) remove(val string) {
	delete(mp, val)
}

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

func StringsToLookupMap(vals []string) LookupMap {
	result := make(LookupMap, len(vals))
	for _, val := range vals {
		result[val] = true
	}
	return result
}
