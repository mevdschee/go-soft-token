#!/usr/bin/env bash

gawk -i inplace -F '=' '/CURRENT_BUILD_VERSION/{$2=$2+1";"}1' OFS='='  version.env

source version.env

package_name="go-soft-token"

#the full list of the platforms run: go tool dist list
platforms=(
#"aix/ppc64"
#"android/386"
#"android/amd64"
#"android/arm"
#"android/arm64"
#"darwin/386"
"darwin/amd64"
#"darwin/arm"
#"darwin/arm64"
#"dragonfly/amd64"
#"freebsd/386"
#"freebsd/amd64"
#"freebsd/arm"
#"illumos/amd64"
#"js/wasm"
#"linux/386"
"linux/amd64"
#"linux/arm"
#"linux/arm64"
#"linux/mips"
#"linux/mips64"
#"linux/mips64le"
#"linux/mipsle"
#"linux/ppc64"
#"linux/ppc64le"
#"linux/s390x"
#"nacl/386"
#"nacl/amd64p32"
#"nacl/arm"
#"netbsd/386"
#"netbsd/amd64"
#"netbsd/arm"
#"netbsd/arm64"
#"openbsd/386"
#"openbsd/amd64"
#"openbsd/arm"
#"openbsd/arm64"
#"plan9/386"
#"plan9/amd64"
#"plan9/arm"
#"solaris/amd64"
#"windows/386"
"windows/amd64"
#"windows/arm"
)

VERSION=$CURRENT_MAJOR_VERSION'.'$CURRENT_MINOR_VERSION'.'$CURRENT_BUILD_VERSION

echo Building version $VERSION

for platform in "${platforms[@]}"
do
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}
    output_name=$package_name'-'$GOOS'-'$GOARCH
    if [ $GOOS = "windows" ]; then
        output_name+='.exe'
    fi

    mkdir -p dist/$VERSION
    env GOOS=$GOOS GOARCH=$GOARCH go build -ldflags "-X main.version=$VERSION" -o dist/$VERSION/$output_name
    if [ $? -ne 0 ]; then
        echo 'An error has occurred! Aborting the script execution...'
        exit 1
    fi
done

go build -ldflags "-X main.version=$VERSION"
