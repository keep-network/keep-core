package electrum

import (
	"fmt"

	"github.com/spf13/pflag"
)

// ProtocolVarFlag is a custom flag to handle `Protocol` type, that is not supported
// by `pflag.FlagSet`.
func ProtocolVarFlag(f *pflag.FlagSet, p *Protocol, name string, value Protocol, usage string) {
	ProtocolVarPFlag(f, p, name, "", value, usage)
}

// ProtocolVarPFlag is a custom flag to handle `Protocol` type, that is not supported
// by `pflag.FlagSet`.
func ProtocolVarPFlag(f *pflag.FlagSet, p *Protocol, name string, short string, value Protocol, usage string) {
	f.VarP(newProtocolValue(value, p), name, short, usage)
}

type protocolValue Protocol

func newProtocolValue(val Protocol, p *Protocol) *protocolValue {
	*p = val
	return (*protocolValue)(p)
}

func (p *protocolValue) Set(s string) error {
	v, ok := ParseProtocol(s)
	if !ok {
		return fmt.Errorf("failed to parse string [%s] as protocol", s)
	}

	*p = protocolValue(v)

	return nil
}

func (p *protocolValue) Type() string {
	return "protocol"
}

func (p *protocolValue) String() string { return (*Protocol)(p).String() }
