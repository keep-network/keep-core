package excfg

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Cfg struct {
	KeyFile         string `json:"keyFile"`
	KeyFilePassword string `json:"keyFilePassword"`
	GethServer      string `json:"gethServer"`
	ContractAddress string `json:"contractAddress"`
	DebugFlag       bool   `json:"debug_on"`
}

func ReadCfg(fn string) (rv Cfg) {

	data, err := ioutil.ReadFile(fn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read configuration file %s, err %s\n", fn, err)
		os.Exit(1)
	}

	err = json.Unmarshal(data, &rv)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read configuration file %s, err %s\n", fn, err)
		os.Exit(1)
	}

	return
}
