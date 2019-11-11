#!/usr/bin/env bash
set -e

cd dependachore
ginkgo -mod vendor -randomizeAllSpecs -randomizeSuites -race -keepGoing -r .
