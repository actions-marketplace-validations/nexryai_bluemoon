package core

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	DockerImage         string   `yaml:"dockerImage"`
	RamLimit            string   `yaml:"ramLimit"`
	BindDir             []string `yaml:"bindDir"`
	AllowWriteToBindDir bool     `yaml:"allowWriteToBindDir"`
	Exec                string   `yaml:"exec"`
	DebugMode           bool     `yaml:"debugMode"`
}

type BuildConfig struct {
	PackageName string `yaml:"packageName"`
	SrcRepo     string `yaml:"srcRepo"`
	SrcDir      string `yaml:"srcDir"`
}

func LoadConfig(configFilePath string) Config {
	// ファイルの読み込み
	data, err := ioutil.ReadFile(configFilePath)
	ExitOnError(err, "An error occurred while loading the configuration file. Are the configuration file paths and permissions correct?")

	// ファイルの内容を構造体にマッピング
	var config Config
	err = yaml.Unmarshal(data, &config)
	ExitOnError(err, "The configuration file was loaded successfully, but the mapping failed.")

	return config
}

func LoadBuildConfig(configFilePath string) BuildConfig {
	// ファイルの読み込み
	data, err := ioutil.ReadFile(configFilePath)
	ExitOnError(err, "An error occurred while loading the build config file. Are the configuration file paths and permissions correct?")

	// ファイルの内容を構造体にマッピング
	var config BuildConfig
	err = yaml.Unmarshal(data, &config)
	ExitOnError(err, "The configuration file was loaded successfully, but the mapping failed.")

	return config
}
