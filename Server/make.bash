#!/usr/bin/env bash

set -e -u -x

BASE_DIR=$( cd "$(dirname "$0")" ; pwd -P )
cd $BASE_DIR

if [ ! -f make.bash ]; then
  echo "make.bash must be run from $GOPATH/src/github.com/payske-dev/vaipn-tunnel-core/Server"
  exit 1
fi

# $1, if specified, is go build tags
if [ -z ${1+x} ]; then BUILD_TAGS=""; else BUILD_TAGS="$1"; fi

export GOCACHE=/tmp

prepare_build () {
  BUILDINFOFILE="vaipnd_buildinfo.txt"
  BUILDDATE=$(date -Iseconds)
  BUILDREPO=$(git config --get remote.origin.url)
  BUILDREV=$(git rev-parse --short HEAD)
  GOVERSION=$(go version | perl -ne '/go version (.*?) / && print $1')

  LDFLAGS="\
  -linkmode external -extldflags \"-static\" \
  -X github.com/payske-dev/vaipn-tunnel-core/vaipn/common/buildinfo.buildDate=$BUILDDATE \
  -X github.com/payske-dev/vaipn-tunnel-core/vaipn/common/buildinfo.buildRepo=$BUILDREPO \
  -X github.com/payske-dev/vaipn-tunnel-core/vaipn/common/buildinfo.buildRev=$BUILDREV \
  -X github.com/payske-dev/vaipn-tunnel-core/vaipn/common/buildinfo.goVersion=$GOVERSION \
  "
  echo -e "${BUILDDATE}\n${BUILDREPO}\n${BUILDREV}\n" > $BUILDINFOFILE

  echo "Variables for ldflags:"
  echo " Build date: ${BUILDDATE}"
  echo " Build repo: ${BUILDREPO}"
  echo " Build revision: ${BUILDREV}"
  echo " Go version: ${GOVERSION}"
  echo ""
}

build_for_linux () {
  prepare_build linux
  GOOS=linux GOARCH=amd64 go build -v -x -tags "${BUILD_TAGS}" -ldflags "$LDFLAGS" -o vaipnd
  if [ $? != 0 ]; then
    echo "...'go build' failed, exiting"
    exit $?
  fi
  chmod 555 vaipnd

  if [ "$1" == "generate" ]; then
    ./vaipnd --ipaddress 0.0.0.0 --web 3000 --protocol SSH:3001 --protocol OSSH:3002 --logFilename /var/log/vaipnd/vaipnd.log generate

    chmod 666 vaipnd.config
    chmod 666 vaipnd-traffic-rules.config
    chmod 666 vaipnd-osl.config
    chmod 666 vaipnd-tactics.config
    chmod 666 server-entry.dat
  fi

}

build_for_linux generate
echo "Done"
