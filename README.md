# localdd

A simple tool inspired by [pow.cx](https://github.com/basecamp/pow) written in Go.

## Configuration

### Configuration Path

#### macOS
```
~/Library/Application Support/localdd
```

### Configure new domain

1. Place a new file with a **test** domain name into the config folder e.g. **abc.test**
2. Insert the destination into this file following format is supported
    * Port on localhost e.g. **8080**
    * Complete url e.g. http://example.com:8080

## Manual installation

### macOS

#### Install Binary

1. Build macOS binary `make darwin` and rename the binary in **build** to *localdd*
2. Place the binary in `/usr/local/Cellar/localdd/sbin/`
3. Place `configuration/macos/launchd/agent/homebrew.tipok.localldd.plist` in `/Library/LaunchAgents/`

#### Configure DNS

1. Place `configuration/macos/resolv` in `/etc/resolver/` and rename it to `test`
2. Place `configuration/macos/dnsmasql.conf` in `/usr/local/etc/dnsmasq.d/` and rename it to `localdd.conf`
3. Add the line `conf-dir=/usr/local/etc/dnsmasq.d/,*.conf` to `/usr/local/etc/dnsmasq.conf`

#### Configure PF

1. Place `configuration/macos/pf.anchors/redirects` in `/usr/local/etc/localdd/pf.anchors/`
2. Place `configuration/macos/sbin/pf-init.sh` in `/usr/local/Cellar/localdd/sbin/`
3. Place `configuration/macos/launchd/daemon/homebrew.tipok.localldd.pf.plist` in `/System/Library/LaunchDaemons/`
