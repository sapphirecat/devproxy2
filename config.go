package main

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/pelletier/go-toml"
)

// ConfigRule describes a rule as listed in the TOML configuration.
type ConfigRule struct {
	MatchHost string `toml:"match_host"`
	MatchPort string `toml:"match_port"`
	DebugRule bool   `toml:"debug_rule"`
	SendTo    string `toml:"send_to"`
}

// ConfigListen describes the listening port, in section 'listen'
type ConfigListen struct {
	Address string `toml:"address" default:"127.0.0.1"`
	Port    int    `toml:"port" default:"8111"`
}

// ConfigOutput describes verbosity options, in section 'output'
type ConfigOutput struct {
	Status        bool `toml:"status" default:"true"`
	DebugAllRules bool `toml:"debug_all_rules"`
	DebugGoProxy  bool `toml:"debug_proxy"`
}

// ConfigFile describes the top-level structure of the TOML configuration.
type ConfigFile struct {
	Servers map[string]Server `toml:"servers"`
	Rules   []ConfigRule      `toml:"rules"`
	Listen  ConfigListen      `toml:"listen"`
	Verbose ConfigOutput      `toml:"output"`
}

// FindDefaultConfig returns the first config file path it finds
func FindDefaultConfig() (string, error) {
	filename := "devproxy.toml"
	exe, err := os.Executable()
	path := filepath.Dir(exe)

	if err != nil {
		return "", err
	}

	paths := []string{
		filepath.Join(path, filename),
		filepath.Join(filepath.Dir(path), "etc", filename),
	}

	for _, file := range paths {
		_, err := os.Stat(file)

		if err == nil {
			return file, nil
		}

		if os.IsNotExist(err) {
			continue
		}

		return "", err
	}

	return "", errors.New("no config file found at any default path")
}

// ReadConfig reads a given configuration file and returns the config
func ReadConfig(file string) (ConfigFile, error) {
	var config ConfigFile
	contents, err := ioutil.ReadFile(file)
	if err != nil {
		return config, err
	}
	err = toml.Unmarshal(contents, &config)
	return config, err
}

// RulesetFromConfig creates an internal ruleset from a configuration source.
func RulesetFromConfig(config ConfigFile) Ruleset {
	rules := NewRuleset(len(config.Rules))

	for i, configRule := range config.Rules {
		var err error
		var ok bool
		var hostExp, portExp *regexp.Regexp
		var serverName string
		var server Server

		matchHost := configRule.MatchHost
		matchPort := configRule.MatchPort

		if matchHost == "" {
			log.Printf("Rule %d missing 'match_host' key", i)
			continue
		}

		if hostExp, err = regexp.Compile(matchHost); err != nil {
			log.Printf("Rule %d match_host = '%s': %v", i, matchHost, err)
			continue
		}

		if matchPort != "" {
			if portExp, err = regexp.Compile(matchPort); err != nil {
				log.Printf("Rule %d match_port = '%s': %v", i, matchPort, err)
				continue
			}
		}

		serverName = configRule.SendTo
		if server, ok = config.Servers[serverName]; !ok {
			log.Printf("Rule %d send_to = '%s': no matching [servers.%s]", i, serverName, serverName)
			continue
		}

		rules.Add(Rule{
			MatchHost: hostExp,
			MatchPort: portExp,
			DebugRule: configRule.DebugRule,
			SendTo:    server,
		})
	}

	return rules
}
