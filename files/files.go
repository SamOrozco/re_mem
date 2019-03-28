package files

import (
	"github.com/SamOrozco/re_mem/data"
	"io/ioutil"
	"os"
	"strings"
)

func Exists(location string) bool {
	if _, err := os.Open(location); os.IsNotExist(err) {
		return false
	}

	return true
}

func CreateDirIfNotExists(loc string) error {
	if Exists(loc) {
		return nil
	}
	return os.Mkdir(loc, os.ModeDir|os.ModePerm)
}

func FileSep() string {
	return string(os.PathSeparator)
}

func WriteNewData(newFileLocation, data string) error {
	return ioutil.WriteFile(newFileLocation, []byte(data), os.ModePerm)
}

func WriteData(fileLocation, data string) error {
	if Exists(fileLocation) {
		return AppendData(fileLocation, data)
	} else {
		return ioutil.WriteFile(fileLocation, []byte(data), os.ModePerm)
	}
}

func AppendData(newFileLocation, data string) error {
	file, err := os.OpenFile(newFileLocation, os.O_RDWR|os.O_APPEND, 0660)
	if err != nil {
		return err
	}
	_, err = file.Write([]byte("\n" + data))
	if err != nil {
		return err
	}
	return nil
}

func ReadDataFromFile(loc string) ([]byte, error) {
	return ioutil.ReadFile(loc)
}

func ReadJsonMapFromFile(loc string) (data.JsonMap, error) {
	dat, err := ioutil.ReadFile(loc)
	if err != nil {
		return nil, err
	}
	return data.ParseJsonBytesToMap(dat)
}

func ReadLinesFromFile(location string) ([]string, error) {
	data, err := ReadDataFromFile(location)
	if err != nil {
		return nil, err
	}
	return strings.Split(string(data), "\n"), nil
}

func DeleteFile(loc string) error {
	return os.Remove(loc)
}

func DeleteDir(loc string) error {
	return os.RemoveAll(loc)
}
