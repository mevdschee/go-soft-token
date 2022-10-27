#!/usr/bin/env bash

gawk -i inplace -F '=' '/CURRENT_BUILD_VERSION/{$2=$2+1}1' OFS='='  version.env

source version.env

package_name="go-soft-token"

#the full list of the platforms run: go tool dist list
platforms=(
"darwin/amd64"
"linux/amd64"
"windows/amd64"
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
