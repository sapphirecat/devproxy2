package main

// Copyright (c) 2013-2020, Sapphire Cat <https://github.com/sapphirecat>.  All
// rights reserved.  See the accompanying LICENSE file for license terms.

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
)

// Args collects all command-line arguments.
type Args struct {
	ConfigFile   string
	DebugSelf    bool
	DebugGoProxy bool
}

// DefineFlags sets up parsing of command-line flags into an Args instance.
func DefineFlags(args *Args) {
	flag.StringVar(&args.ConfigFile, "config", "", "Configuration file to use")
	flag.BoolVar(&args.DebugSelf, "verbose", false, "Enables logging of devproxy rule results")
	flag.BoolVar(&args.DebugGoProxy, "debug", false, "Enables excessive logging in goproxy")
}

// resolveConfig does all argument parsing and loading
func resolveConfig() (*ConfigFile, error) {
	var args Args
	var configPath string
	var config ConfigFile
	var err error

	// parse args
	DefineFlags(&args)
	flag.Parse()

	// resolve config
	if args.ConfigFile != "" {
		configPath = args.ConfigFile
	} else {
		configPath, err = FindDefaultConfig()
		if err != nil {
			return nil, err
		}
	}

	// parse config
	config, err = ReadConfig(configPath)
	if err != nil {
		return nil, err
	}

	// combine and return results
	err = MergeConfig(&config, args)
	return &config, err
}

// DoMain runs the main logic, given the parsed Args.
func DoMain(config ConfigFile) error {
	ruleSet := RulesetFromConfig(config)
	proxy := NewServer(ruleSet, config.Verbose.DebugGoProxy)

	listen := config.Listen
	listenAddr := fmt.Sprintf("%s:%d", listen.Address, listen.Port)
	if config.Verbose.Status {
		log.Println("listening on", listenAddr, "with", ruleSet.Length(), "active rules")
	}
	return http.ListenAndServe(listenAddr, proxy)
}

// MergeConfig overrides configuration file options from CLI arguments
func MergeConfig(config *ConfigFile, args Args) error {
	if config == nil {
		return errors.New("no config given to merge args into")
	}

	c := *config

	// -verbose and -debug
	if args.DebugSelf {
		c.Verbose.DebugAllRules = true
	}
	if args.DebugGoProxy {
		c.Verbose.DebugGoProxy = true
	}

	return nil
}

// main parses command-line arguments and runs the entire server.
func main() {
	config, err := resolveConfig()
	if err != nil {
		log.Fatal("Cannot read configuration:", err)
	}

	log.Fatal(DoMain(*config))
}
