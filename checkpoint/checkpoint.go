package checkpoint

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"sync"
)

var lock sync.Mutex

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func Save(cname string, data interface{}) {
	lock.Lock()
	defer lock.Unlock()

	f, err := os.Create(cname)
	check(err)
	defer f.Close()

	dataBytes, err := json.MarshalIndent(data, "", "\t")
	check(err)
	dataReader := bytes.NewReader(dataBytes)
	_, err = io.Copy(f, dataReader)
	check(err)
}

func Load(cname string, data interface{}) {
	lock.Lock()
	defer lock.Unlock()

	f, err := os.Open(cname)
	check(err)
	defer f.Close()

	dec := json.NewDecoder(f)
	err = dec.Decode(data)
	check(err)
}
