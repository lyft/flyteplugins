// Code generated by go generate; DO NOT EDIT.
// This file was generated by robots.

package catalog

import (
	"encoding/json"
	"reflect"

	"fmt"

	"github.com/spf13/pflag"
)

// If v is a pointer, it will get its element value or the zero value of the element type.
// If v is not a pointer, it will return it as is.
func (Config) elemValueOrNil(v interface{}) interface{} {
	if t := reflect.TypeOf(v); t.Kind() == reflect.Ptr {
		if reflect.ValueOf(v).IsNil() {
			return reflect.Zero(t.Elem()).Interface()
		} else {
			return reflect.ValueOf(v).Interface()
		}
	} else if v == nil {
		return reflect.Zero(t).Interface()
	}

	return v
}

func (Config) mustMarshalJSON(v json.Marshaler) string {
	raw, err := v.MarshalJSON()
	if err != nil {
		panic(err)
	}

	return string(raw)
}

// GetPFlagSet will return strongly types pflags for all fields in Config and its nested types. The format of the
// flags is json-name.json-sub-name... etc.
func (cfg Config) GetPFlagSet(prefix string) *pflag.FlagSet {
	cmdFlags := pflag.NewFlagSet("Config", pflag.ExitOnError)
	cmdFlags.Int(fmt.Sprintf("%v%v", prefix, "ReaderWorkqueueConfig.workers"), defaultConfig.ReaderWorkqueueConfig.Workers, "Number of concurrent workers to start processing the queue.")
	cmdFlags.Int(fmt.Sprintf("%v%v", prefix, "ReaderWorkqueueConfig.maxRetries"), defaultConfig.ReaderWorkqueueConfig.MaxRetries, "Maximum number of retries per item.")
	cmdFlags.Int(fmt.Sprintf("%v%v", prefix, "ReaderWorkqueueConfig.maxItems"), defaultConfig.ReaderWorkqueueConfig.IndexCacheMaxItems, "Maximum number of entries to keep in the index.")
	cmdFlags.Int(fmt.Sprintf("%v%v", prefix, "WriterWorkqueueConfig.workers"), defaultConfig.WriterWorkqueueConfig.Workers, "Number of concurrent workers to start processing the queue.")
	cmdFlags.Int(fmt.Sprintf("%v%v", prefix, "WriterWorkqueueConfig.maxRetries"), defaultConfig.WriterWorkqueueConfig.MaxRetries, "Maximum number of retries per item.")
	cmdFlags.Int(fmt.Sprintf("%v%v", prefix, "WriterWorkqueueConfig.maxItems"), defaultConfig.WriterWorkqueueConfig.IndexCacheMaxItems, "Maximum number of entries to keep in the index.")
	return cmdFlags
}
