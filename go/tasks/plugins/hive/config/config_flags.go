// Code generated by go generate; DO NOT EDIT.
// This file was generated by robots.

package config

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
	cmdFlags.String(fmt.Sprintf("%v%v", prefix, "endpoint"), defaultConfig.Endpoint.String(), "Endpoint for qubole to use")
	cmdFlags.String(fmt.Sprintf("%v%v", prefix, "commandApiPath"), defaultConfig.CommandAPIPath.String(), "API Path where commands can be launched on Qubole. Should be a valid url.")
	cmdFlags.String(fmt.Sprintf("%v%v", prefix, "analyzeLinkPath"), defaultConfig.AnalyzeLinkPath.String(), "URL path where queries can be visualized on qubole website. Should be a valid url.")
	cmdFlags.String(fmt.Sprintf("%v%v", prefix, "quboleTokenKey"), defaultConfig.TokenKey, "Name of the key where to find Qubole token in the secret manager.")
	cmdFlags.String(fmt.Sprintf("%v%v", prefix, "defaultClusterLabel"), defaultConfig.DefaultClusterLabel, "The default cluster label. This will be used if label is not specified on the hive job.")
	cmdFlags.Int(fmt.Sprintf("%v%v", prefix, "readRateLimiter.qps"), defaultConfig.ReadRateLimiter.QPS, "Defines the max rate of calls per second.")
	cmdFlags.Int(fmt.Sprintf("%v%v", prefix, "readRateLimiter.burst"), defaultConfig.ReadRateLimiter.Burst, "Defines the maximum burst size.")
	cmdFlags.Int(fmt.Sprintf("%v%v", prefix, "writeRateLimiter.qps"), defaultConfig.WriteRateLimiter.QPS, "Defines the max rate of calls per second.")
	cmdFlags.Int(fmt.Sprintf("%v%v", prefix, "writeRateLimiter.burst"), defaultConfig.WriteRateLimiter.Burst, "Defines the maximum burst size.")
	cmdFlags.Int(fmt.Sprintf("%v%v", prefix, "caching.size"), defaultConfig.Caching.Size, "Defines the maximum number of items to cache.")
	cmdFlags.String(fmt.Sprintf("%v%v", prefix, "caching.resyncInterval"), defaultConfig.Caching.ResyncInterval.String(), "Defines the sync interval.")
	cmdFlags.Int(fmt.Sprintf("%v%v", prefix, "caching.workers"), defaultConfig.Caching.Workers, "Defines the number of workers to start up to process items.")
	return cmdFlags
}
