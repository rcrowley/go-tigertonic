package tigertonic

import (
	"encoding/json"
	"io/ioutil"
)

// Configure reads the given configuration file and unmarshals the JSON found
// into the given configuration structure, which may be any Go type.
func Configure(pathname string, i interface{}) error {
	if "" == pathname {
		return nil
	}
	buf, err := ioutil.ReadFile(pathname)
	if nil != err {
		return err
	}
	err = json.Unmarshal(buf, i)
	if nil != err {
		return err
	}
	return nil
}
