#!/bin/zsh

export PATH="/opt/homebrew/opt/go/libexec/bin:/Users/obeng/go/bin:/opt/homebrew/bin:/opt/homebrew/sbin:/usr/local/bin:/System/Cryptexes/App/usr/bin:/usr/bin:/bin:/usr/sbin:/sbin:/var/run/com.apple.security.cryptexd/codex.system/bootstrap/usr/local/bin:/var/run/com.apple.security.cryptexd/codex.system/bootstrap/usr/bin:/var/run/com.apple.security.cryptexd/codex.system/bootstrap/usr/appleinternal/bin:/usr/local/go/bin
"
# Optional: debug output
echo "PATH=$PATH"
which mockgen
mockgen -version

go generate ./...