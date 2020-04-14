# DevProxy

Use your production URLs to reach your development environment.


# Functionality

devproxy reads all connection requests, and either sends them to their normal
destination, or forwards them to an alternate server. Once the connection is
made, data is copied back and forth with no interference.

HTTPS is handled by making the connection, then letting the browser do the TLS
handshake. In this case, the alternate server should have a valid certificate
for the domain, because the TLS connection is made between server and browser,
not devproxy.


# Why?

A staging environment should be as close to production as physically possible.
Every single difference may be a source of bugs in production that *cannot* be
detected in staging.

Modifying `/etc/hosts` all the time to point the production DNS at the staging
server (or comment said redirection back out) is tedious, invisible, and
requires frequent privilege escalations.

With the proxy switcher, switching the browser between truth and lie is improved
in all respects: a fast user-level action with a status indicator.  All that you
need is a proxy to transparently connect the staging backend when a production
URL is requested… and that proxy is devproxy.


# Examples

## Basic Example

Given the configuration:

    [[rules]]
    match_host = '\bexample.com$'
    send_to = local

    [servers.local]
    address = 127.0.0.1
    http_port = 8080
    https_port = 8443

When the browser requests http://example.com/gnorc from the proxy, then the
following happens:

1. devproxy connects to 127.0.0.1:8080
2. devproxy requests `/gnorc`
3. devproxy copies the request headers, including `Host: example.com`
4. devproxy copies the full response to the browser

## TLS Example

Given the same configuration as above, but the browser requests a connection
to `example.com:443` (because the user asked for https://example.com/rhinoc
for example):

1. devproxy connects to 127.0.0.1:8443 and responds with “OK”
2. The browser and the server do a TLS handshake, which succeeds without
   warnings if the server has the `example.com` certificate and key available
3. The browser makes a request with `Host: example.com` inside the tunnel
4. The server returns a response via the tunnel

Swapping the backend server is the extent of devproxy’s meddling, so it has no
need to break the secure connection and read the encrypted data. It simply
copies the encrypted data between browser and server.


# Using it

1. Clone anywhere with go 1.13 or newer.  You should end up with a
	`devproxy2` directory, containing this README.md file.
2. Edit [devproxy.toml](./devproxy.toml) to configure devproxy for your
    environment. This is a minimal file; a much longer example with all possible
    options, and comments for them, is in the
    [devproxy-full.toml](./devproxy-full.toml) file.
3. Run `go build`
4. Run `./devproxy` (Linux/OS X) or `devproxy.exe` (Windows).
5. Set your web proxy to 127.0.0.1:8111, or whatever was configured in your
    `devproxy.toml` file.

I use [Proxy Switcher and Manager](https://addons.mozilla.org/en-US/firefox/addon/proxy-switcher-and-manager/)
with Firefox so that I can easily switch between using devproxy or not, and
see at a glance whether I _am_ using it.


# Command line flags

## -config

Specifies a configuration file to use. By default, devproxy looks in the
directory containing the binary, and ../etc relative to the binary. That is, if
devproxy were installed as `/usr/local/bin/devproxy`, then the default locations
are:

1. `/usr/local/bin/devproxy.toml`
2. `/usr/local/etc/devproxy.toml`

Using `-config=/home/utena/devproxy.toml` causes that specific file to be
loaded, avoiding the search entirely.

## -verbose and -debug

With `-verbose`, devproxy logs requests it receives, and the decisions taken
overall.  Logs may also be enabled on a per-rule level with `debug_rule = true`;
this causes the related rule to log its match decisions in more detail, even
without `-verbose` specified.

With `-debug`, devproxy tells goproxy to log what _it_ is doing.

These options are fully independent.  Neither affects nor implies the other.


# Compatibility

devproxy2 freezes a new API and configuration file format; changes in devproxy2
will be backwards-compatible with version 2.0.0.

devproxy2 supersedes [devproxy](https://github.com/sapphirecat/devproxy) and
makes the following major changes:

- Configuration is done by file, instead of built into the binary
- Matching occurs on host and port explicitly, instead of trying one regexp
    against either a "host" or "host:port" format string
- Matching a rule results in a specific redirection, not running a function
- The `-target`, `-listen`, and `-port` command-line options are removed
- `debug_rule` configuration option added to allow verbosity per-rule
- Converted to Go Modules
- Restructured as a single application, not a library


# License

SPDX BSD-3-Clause; full text in [LICENSE](./LICENSE).
