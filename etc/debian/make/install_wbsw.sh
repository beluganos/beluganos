#! /bin/bash -e
# -*- coding: utf-8 -*-

if [ -n "$1" ]; then
    cd $1
fi

dpkg -i ./*.deb
