package archiveManager

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func LoadJson(jsonFile string, data interface{}) (err error) {
	var fp *os.File
	var fileData []byte
	if _, err = os.Stat(jsonFile); err != nil {
		return
	}
	if fp, err = os.OpenFile(jsonFile, os.O_RDONLY, 0755); err != nil {
		return
	}
	defer fp.Close()
	if fileData, err = ioutil.ReadAll(fp); err != nil {
		return
	}
	if err = json.Unmarshal(fileData, data); err != nil {
		return
	}
	return nil
}

func SaveJson(jsonFile string, data interface{}) (err error) {
	var fp *os.File
	var fileData []byte
	if fileData, err = json.Marshal(data); err != nil {
		return
	}
	if fp, err = os.OpenFile(jsonFile, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0755); err != nil {
		return
	}
	defer fp.Close()
	if _, err = fp.Write(fileData); err != nil {
		return
	}
	return nil
}

func GetFileHash(data []byte) string {
	return fmt.Sprintf("%x", md5.Sum(data))
}
