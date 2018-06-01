package runtime

import (
	//. "utils"
	"reflect"
)

// AppOption is the function parameter type for adding a runtime argument
type AppOption func(*options) error

// KeyVal is a basic option type
type KeyVal map[string]string

// OsArgs is the list of runtime arguements
type OsArgs []string

// OsArguments is the list of arguments passed to the CLI at runtime
func OsArguments(args []string) AppOption {
	return func(o *options) error {
		o.osArgs = args
		return nil
	}
}

// Argument is an argument passed programmatically
func Argument(kv KeyVal) AppOption {
	return func(o *options) error {
		var key string
		keys := reflect.ValueOf(kv).MapKeys()
		for i := 0; i < len(keys); i++ {
			key = keys[i].String()
			o.arguments[key] = kv[key]
		}
		return nil
	}
}

// Command specifies the command to run
func Command(name string) AppOption {
	return func(o *options) error {
		o.command = name
		return nil
	}
}
