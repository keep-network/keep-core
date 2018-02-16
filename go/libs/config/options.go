// Copyright 2018 The Keep Authors.  See LICENSE.md for details.
package config

import (
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"reflect"
	"strconv"
	"github.com/pkg/errors"
)

type Conf struct {
	AppEnv               string   `toml:"app_env"`
	ProjectRoot          string   `toml:"project_root_dir"`
	DownloadDir          string   `toml:"download_dir"`
	GcpSourceDir         string   `toml:"gcp_source_dir"`
	GcpSourceKeyFile     string   `toml:"gcp_source_key_file"`
	GcpSourceProjectId   string   `toml:"gcp_source_project_id"`
	GcpSinkDir           string   `toml:"gcp_sink_dir"`
	GcpSinkKeyFile       string   `toml:"gcp_sink_key_file"`
	GcpSinkProjectId     string   `toml:"gcp_sink_project_id"`
	SourceBucketName     string   `toml:"source_bucket_name"`
	SinkBucketName       string   `toml:"sink_bucket_name"`
	ApiPort              string   `toml:"api_port"`
	LogTimeTrack         bool     `toml:"log_timetrack"`
	LogFullStackTrace    bool     `toml:"log_full_stack_trace"`
	LogDebugInfo         bool     `toml:"log_debug_info"`
	LogDebugInfoForTests bool     `toml:"log_debug_info_for_tests"`
}

var Config Conf

func GetOptions() bool {
	var configFile string

	flag.StringVar(&configFile, "config", "", "Configuration file")
	flag.StringVar(&Config.AppEnv, "onion-env", "development", "Runtime environment. Determines whether to run semver scripts or expect env vars")
	flag.StringVar(&Config.ProjectRoot, "project-root", "/Users/lex/clients/packt/dev/fp-go/2-design-patterns/ch06-onion-arch/04_onion", "Project Root directory (required). Must be absolute path")
	flag.StringVar(&Config.DownloadDir, "download-dir", "/Users/lex/clients/packt/dev/fp-go/2-design-patterns/ch06-onion-arch/04_onion/downloads", "Where files are downloaded")
	flag.StringVar(&Config.GcpSourceDir, "gcp_source_dir", "source-events", "Source log file host names (parent directory names in Google Cloud bucket")
	flag.StringVar(&Config.GcpSourceKeyFile, "gcp-source-key-file", "/Users/lex/clients/packt/dev/fp-go/2-design-patterns/ch06-onion-arch/04_onion/keys/google-cloud-storage/source/onion-source-key.json", "Google Cloud Platform source key file")
	flag.StringVar(&Config.GcpSourceProjectId, "gcp-source-project-id", "onion-xxxx", "Google Cloud Platform source project id")
	flag.StringVar(&Config.GcpSourceDir, "gcp_sink_dir", "sink-events", "Source log file host names (parent directory names in Google Cloud bucket")
	flag.StringVar(&Config.GcpSinkKeyFile, "gcp-sink-key-file", "/Users/lex/clients/packt/dev/fp-go/2-design-patterns/ch06-onion-arch/04_onion/keys/google-cloud-storage/sink/onion-sink-key.json", "Google Cloud Platform sink key file")
	flag.StringVar(&Config.GcpSinkProjectId, "gcp_sink_project_id", "xxxx-999999", "Google Cloud Platform sink project id")
	flag.StringVar(&Config.SourceBucketName, "source_bucket_name", "onion-logs", "Cloud bucket where the log files to be processed")
	flag.StringVar(&Config.SinkBucketName, "sink_bucket_name", "lexttc3-my-backup-bucket", "Cloud bucket where they are to be copied after processing")
	flag.StringVar(&Config.ApiPort, "api-port", "8080", "Port that the API listens on")
	flag.BoolVar(&Config.LogTimeTrack, "log-timetrack", true, "Enable or disable logging of utils/TimeTrack() (For benchmarking/debugging)")
	flag.BoolVar(&Config.LogFullStackTrace, "log-full-stack-trace", false, "Print version information and exit")
	flag.BoolVar(&Config.LogDebugInfo, "log-debug-info", false, "Whether to log debug output to the log (set to true for debug purposes)")
	flag.BoolVar(&Config.LogDebugInfoForTests, "log-debug-info-for-tests", true, "Whether to log debug output to the log when running tests (set to true for debug purposes)")

	flag.Parse()

	if configFile != "" {
		if _, err := toml.DecodeFile(configFile, &Config); err != nil {
			HandlePanic(errors.Wrap(err, "unable to read config file"))
		}
	}
	return true
}

type Datastore interface {}


func UpdateConfigVal(d Datastore, key, val string) (oldValue string) {
	Debug.Printf("key (%s), val (%v)\n", key, val)
	value := reflect.ValueOf(d)
	if value.Kind() != reflect.Ptr {
		panic("not a pointer")
	}
	valElem := value.Elem()
	for i := 0; i < valElem.NumField(); i++ {
		tag := valElem.Type().Field(i).Tag
		field := valElem.Field(i)
		switch tag.Get("toml") {
		case key:
			if fmt.Sprintf("%v", field.Kind()) == "int" {
				oldValue = strconv.FormatInt(field.Int(), 10)
				intVal, err := strconv.Atoi(val)
				if err != nil {
					fmt.Printf("could not parse int, key(%s) val(%s)", key, val)
				} else {
					field.SetInt(int64(intVal))
				}
			} else if fmt.Sprintf("%v", field.Kind()) == "bool" {
				oldValue = strconv.FormatBool(field.Bool())
				b, err := strconv.ParseBool(val)
				if err != nil {
					fmt.Printf("could not parse bool, key(%s) val(%s)", key, val)
				} else {
					field.SetBool(b)
				}
			} else {
				// Currently only supports bool, int and string
				oldValue = field.String()
				field.SetString(val)
			}
		}
	}
	return
}

