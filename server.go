package main

// Copyright (c) 2013-2020, Sapphire Cat <https://github.com/sapphirecat>.  All
// rights reserved.  See the accompanying LICENSE file for license terms.

import (
	"log"
	"net/http"

	"gopkg.in/elazarl/goproxy.v1"
)

// UseRulesetForConnect returns the goproxy function to handle TLS-secured HTTPS.
func UseRulesetForConnect(ruleset Ruleset) func(string, *goproxy.ProxyCtx) (*goproxy.ConnectAction, string) {
	return func(host string, ctx *goproxy.ProxyCtx) (*goproxy.ConnectAction, string) {
		target := getTarget(ruleset, host, RuleForTLS)
		if target == NoOp {
			target = host
			if ruleset.Verbose {
				log.Println("!match HTTPS", host)
			}
		} else if ruleset.Verbose {
			log.Println("+HTTPS", host, ctx.Req.URL.Path)
		}

		return goproxy.OkConnect, target
	}
}

// UseRulesetForHTTP returns the goproxy function to handle plain-text HTTP.
func UseRulesetForHTTP(ruleset Ruleset) func(*http.Request, *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	return func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		host := r.URL.Host
		target := getTarget(ruleset, host, RuleForHTTP)
		if target != NoOp {
			r.URL.Host = target
			if ruleset.Verbose {
				log.Println("+plain", host, r.URL.Path)
			}
		} else if ruleset.Verbose {
			log.Println("!match plain", r.URL.Host)
		}

		return r, nil
	}
}

// UseRuleset configures devproxy to process the rules on goproxy requests.
func UseRuleset(proxy *goproxy.ProxyHttpServer, rules Ruleset) {
	proxy.OnRequest().HandleConnectFunc(UseRulesetForConnect(rules))
	proxy.OnRequest().DoFunc(UseRulesetForHTTP(rules))
}

// NewServer creates a goproxy that is configured to use the Ruleset.
func NewServer(r Ruleset, debugProxy bool) http.Handler {
	// Create the proxy and set debugging on it
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = debugProxy

	// Set up rules *in goproxy* that apply the ruleset to each request.
	UseRuleset(proxy, r)
	return proxy
}
