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
	cmdFlags.Int(fmt.Sprintf("%v%v", prefix, "jobStoreConfig.jacheSize"), defaultConfig.JobStoreConfig.CacheSize, "Maximum informer cache size as number of items. Caches are used as an optimization to lessen the load on AWS Services.")
	cmdFlags.Int(fmt.Sprintf("%v%v", prefix, "jobStoreConfig.parallelizm"), defaultConfig.JobStoreConfig.Parallelizm, "")
	cmdFlags.Int(fmt.Sprintf("%v%v", prefix, "jobStoreConfig.batchChunkSize"), defaultConfig.JobStoreConfig.BatchChunkSize, "Determines the size of each batch sent to GetJobDetails api.")
	cmdFlags.Int(fmt.Sprintf("%v%v", prefix, "defCacheSize"), defaultConfig.JobDefCacheSize, "Maximum job definition cache size as number of items. Caches are used as an optimization to lessen the load on AWS Services.")
	cmdFlags.Int64(fmt.Sprintf("%v%v", prefix, "getRateLimiter.rate"), defaultConfig.GetRateLimiter.Rate, "Allowed rate of calls per second.")
	cmdFlags.Int(fmt.Sprintf("%v%v", prefix, "getRateLimiter.burst"), defaultConfig.GetRateLimiter.Burst, "Allowed burst rate of calls.")
	cmdFlags.Int64(fmt.Sprintf("%v%v", prefix, "defaultRateLimiter.rate"), defaultConfig.DefaultRateLimiter.Rate, "Allowed rate of calls per second.")
	cmdFlags.Int(fmt.Sprintf("%v%v", prefix, "defaultRateLimiter.burst"), defaultConfig.DefaultRateLimiter.Burst, "Allowed burst rate of calls.")
	cmdFlags.Int64(fmt.Sprintf("%v%v", prefix, "maxArrayJobSize"), defaultConfig.MaxArrayJobSize, "Maximum size of array job.")
	cmdFlags.Int32(fmt.Sprintf("%v%v", prefix, "minRetries"), defaultConfig.MinRetries, "Minimum number of retries")
	cmdFlags.Int32(fmt.Sprintf("%v%v", prefix, "maxRetries"), defaultConfig.MaxRetries, "Maximum number of retries")
	cmdFlags.Int(fmt.Sprintf("%v%v", prefix, "maxErrLength"), defaultConfig.MaxErrorStringLength, "Determines the maximum length of the error string returned for the array.")
	cmdFlags.String(fmt.Sprintf("%v%v", prefix, "roleAnnotationKey"), defaultConfig.RoleAnnotationKey, "Map key to use to lookup role from task annotations.")
	cmdFlags.String(fmt.Sprintf("%v%v", prefix, "resyncPeriod"), defaultConfig.ResyncPeriod.String(), "Defines the duration for syncing job details from AWS Batch.")
	cmdFlags.Int(fmt.Sprintf("%v%v", prefix, "outputAssembler.workers"), defaultConfig.OutputAssembler.Workers, "Number of concurrent workers to start processing the queue.")
	cmdFlags.Int(fmt.Sprintf("%v%v", prefix, "outputAssembler.maxRetries"), defaultConfig.OutputAssembler.MaxRetries, "Maximum number of retries per item.")
	cmdFlags.Int(fmt.Sprintf("%v%v", prefix, "outputAssembler.maxItems"), defaultConfig.OutputAssembler.IndexCacheMaxItems, "Maximum number of entries to keep in the index.")
	cmdFlags.Int(fmt.Sprintf("%v%v", prefix, "errorAssembler.workers"), defaultConfig.ErrorAssembler.Workers, "Number of concurrent workers to start processing the queue.")
	cmdFlags.Int(fmt.Sprintf("%v%v", prefix, "errorAssembler.maxRetries"), defaultConfig.ErrorAssembler.MaxRetries, "Maximum number of retries per item.")
	cmdFlags.Int(fmt.Sprintf("%v%v", prefix, "errorAssembler.maxItems"), defaultConfig.ErrorAssembler.IndexCacheMaxItems, "Maximum number of entries to keep in the index.")
	return cmdFlags
}
