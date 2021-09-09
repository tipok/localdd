#!/usr/bin/env zsh

/sbin/pfctl -ef - <<EOF
rdr-anchor "homebrew.tipok.localdd/*"
EOF

# Flush existing rules.
# Not that useful for our one-shot, but could be useful in more complicated setups
/sbin/pfctl -a homebrew.tipok.localdd -F all
/sbin/pfctl -a homebrew.tipok.localdd/redirect -F all

/sbin/pfctl -a homebrew.tipok.localdd/redirect -f /usr/local/etc/localdd/pf.anchors/redirects
