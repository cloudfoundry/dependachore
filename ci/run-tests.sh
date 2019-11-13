#!/usr/bin/env bash
set -e

export GOOGLE_APPLICATION_CREDENTIALS="$PWD/$GOOGLE_APPLICATION_CREDENTIALS"

cd dependachore
ginkgo -mod vendor -randomizeAllSpecs -randomizeSuites -race -keepGoing -r .
