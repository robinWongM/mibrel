#!/usr/bin/env bash

USERNAME="${USERNAME:-"${_REMOTE_USER:-"automatic"}"}"

if [ "$(id -u)" -ne 0 ]; then
    echo -e 'Script must be run as root. Use sudo, su, or add "USER root" to your Dockerfile before running this script.'
    exit 1
fi

# Determine the appropriate non-root user
if [ "${USERNAME}" = "auto" ] || [ "${USERNAME}" = "automatic" ]; then
    USERNAME=""
    POSSIBLE_USERS=("vscode" "node" "codespace" "$(awk -v val=1000 -F ":" '$3==val{print $1}' /etc/passwd)")
    for CURRENT_USER in "${POSSIBLE_USERS[@]}"; do
        if id -u ${CURRENT_USER} > /dev/null 2>&1; then
            USERNAME=${CURRENT_USER}
            break
        fi
    done
    if [ "${USERNAME}" = "" ]; then
        USERNAME=root
    fi
elif [ "${USERNAME}" = "none" ] || ! id -u ${USERNAME} > /dev/null 2>&1; then
    USERNAME=root
fi

apt update && apt install -y protobuf-compiler libprotobuf-dev && apt-get clean

su $USERNAME -s /bin/bash -c '
    # buf
    go install github.com/bufbuild/buf/cmd/buf@latest
    # grpcurl
    go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
    # protoc-gen-go
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
    # protoc-gen-connect-go
    go install github.com/bufbuild/connect-go/cmd/protoc-gen-connect-go@latest

    # bun
    curl -fsSL https://bun.sh/install | bash
'