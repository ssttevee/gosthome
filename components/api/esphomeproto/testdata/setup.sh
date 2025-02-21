#!/bin/bash
# Copyright 2021 The Periph Authors. All rights reserved.
# Use of this source code is governed under the Apache License, Version 2.0
# that can be found in the LICENSE file.

set -eu

pushd "$(dirname $0)" > /dev/null

# Set it to -v for verbosity.
QUIET=-q
# QUIET=

if ! ([ -f venv/bin/activate ] && diff -q requirements.txt venv/requirements.txt > /dev/null 2>&1  ); then
    (
        if [ ! -d venv ]; then
        mkdir venv
        fi

        if [ ! -f ./venv/bin/activate ]; then
        python3 -m venv venv
        fi

        echo "- Activating virtualenv"
        source venv/bin/activate

        echo "- Installing requirements"
        pip3 install $QUIET -U pip
        pip3 install $QUIET -U -r requirements.txt
        cp requirements.txt venv/requirements.txt

        echo ""
        echo "Congratulations! Everything is inside ./venv/"
        echo "To access esphome, run:"
        echo "  source $PWD/venv/bin/activate"
        deactivate
    ) > .venv.init.log
fi
source venv/bin/activate
popd > /dev/null
exec python3 "$@"
