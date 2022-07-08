package main

// Copyright (c) 2013-2020, Sapphire Cat <https://github.com/sapphirecat>.  All
// rights reserved.  See the accompanying LICENSE file for license terms.

import (
	"log"
	"net"
	"regexp"
	"strconv"
	"strings"
)

const (
	// RuleForHTTP indicates an interception of a plain-text HTTP connection.
	RuleForHTTP = iota
	// RuleForTLS indicates an interception of a TLS-secured HTTPS connection.
	RuleForTLS
	// NoOp represents "no modification" for a request, e.g. for the HTTPS
	// branch of an HTTP-only action.
	NoOp = ""
)

// Rule represents a single host[:port] interception and action to take.
type Rule struct {
	MatchHost *regexp.Regexp
	MatchPort *regexp.Regexp
	DebugRule bool
	SendTo    Server
}

// Server represents a destination to forward intercepted traffic to.
type Server struct {
	Address   string `toml:"address"`
	HTTPPort  int    `toml:"http_port" default:"80"`
	HTTPSPort int    `toml:"https_port" default:"443"`
}

// Mode is either RuleForHttp (plain text) or RuleForTls (TLS)
type Mode int

// Ruleset contains the list of interception Rules.
type Ruleset struct {
	Verbose bool
	items   []Rule
}

// Add adds a Rule to a Ruleset.
func (r *Ruleset) Add(a Rule) {
	r.items = append(r.items, a)
}

// Length returns the number of Rules that were added to the Ruleset.
func (r *Ruleset) Length() int {
	return len(r.items)
}

// NewRuleset creates a new Ruleset.  The capacity can be pre-set.
func NewRuleset(capacity int) Ruleset {
	return Ruleset{
		items: make([]Rule, 0, capacity),
	}
}

func getTarget(rules Ruleset, hostname string, mode Mode) string {
	host, port, err := net.SplitHostPort(hostname)

	if err != nil {
		// Parse the error string, and hope it's never changed/localized....
		if strings.Contains(err.Error(), "missing port") {
			host = hostname
			port = "80"
		} else {
			log.Fatalf("Bad hostname: %s: %s", err.Error(), hostname)
		}
	}

	for i := range rules.items {
		// take a pointer to prevent range from copying the rule struct
		rule := &rules.items[i]

		if !rule.MatchHost.MatchString(host) {
			if rule.DebugRule {
				text := rule.MatchHost.String()
				log.Println("!match: host", host, "with", text)
			}
			continue
		}
		if rule.MatchPort != nil && !rule.MatchPort.MatchString(port) {
			if rule.DebugRule {
				text := rule.MatchPort.String()
				log.Println("!match: port", port, "with", text)
			}
			continue
		}

		destination := rule.SendTo.Address
		var destPort int
		if mode == RuleForHTTP {
			destPort = rule.SendTo.HTTPPort
		} else {
			destPort = rule.SendTo.HTTPSPort
		}

		destHostPort := net.JoinHostPort(destination, strconv.Itoa(destPort))
		if rule.DebugRule {
			srcHostPort := net.JoinHostPort(host, port)
			log.Printf("match for %s -> %s\n", srcHostPort, destHostPort)
		}

		return destHostPort
	}

	// no rules specified an operation
	return NoOp
}
