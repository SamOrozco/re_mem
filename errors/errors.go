package errors

import (
	"fmt"
	"log"
)

func DirNotExist(loc string) {
	log.Fatal(fmt.Sprintf("directory %s does not exist", loc))
}

func InitDirError(loc string, err error) {
	log.Fatal(fmt.Errorf("unable to create dir %s, error %s", loc, err))
}
