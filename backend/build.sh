#!/usr/bin/env bash
go mod tidy && go mod download

build_with_osarch() {
    osarch=$1
    goos=$(echo "$osarch" | cut -d '/' -f1)
    if [ -z "$goos" ]; then
        goos=$(go env GOOS)
    fi
    goarch=$(echo "$osarch" | cut -d '/' -f2)
    if [ -z "$goarch" ] || [ "$goarch" = "$osarch" ]; then
        goarch=$(go env GOARCH)
    fi
    export GOOS=$goos
    export GOARCH=$goarch
    output=$2
    # mod_basename=$(basename "$(go list -m)")
    mod_basename=$(basename "$(head -n 1 go.mod | cut -f 2 -d " ")")
    if [ -d "$output" ]; then
        if [ "${output: -1}" = "/" ]; then
            output="${output}#MOD#_#OS#_#ARCH#"
        else
            output="${output}/#MOD#_#OS#_#ARCH#"
        fi
    fi
    if [ -z "$output" ]; then
        output="#MOD#_#OS#_#ARCH#"
    fi
    suffix=""
    if [ "$GOOS" = "windows" ] && [ "${output: -4}" != ".exe" ]; then
        suffix=".exe"
    fi
    output=$(echo "${output}${suffix}" | sed "s/#MOD#/${mod_basename}/g" | sed "s/#OS#/${GOOS}/g" | sed "s/#ARCH#/${GOARCH}/g")
    echo "building for $(go env GOOS)/$(go env GOARCH)  ==>  ${output}"
    go build -ldflags "-s -w $LDFLAGS" \
        -gcflags="all=-trimpath=${PWD}" \
        -asmflags="all=-trimpath=${PWD}" \
        -o "${output}"
}

output_dir="./bin/"

if [ -d $output_dir ]; then
    rm -rf $output_dir
fi

mkdir -p $output_dir

if [ -z "$LDFLAGS" ]; then
    LDFLAGS="-X website-testing/config.Version=$(date -Iseconds -u)"
fi
targets=("linux/amd64" "linux/arm64" "windows/arm64" "windows/amd64" "darwin/amd64" "darwin/arm64")

for osarch in "${targets[@]}"; do
    build_with_osarch "$osarch" "$output_dir" &
done
wait
echo "build completed"
