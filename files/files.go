package files

import (
	"io/ioutil"
	"os"
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

func WriteData(newFileLocation, data string) error {
	return ioutil.WriteFile(newFileLocation, []byte(data), os.ModePerm)
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
