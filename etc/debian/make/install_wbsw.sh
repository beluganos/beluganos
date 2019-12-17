<<<<<<< HEAD
#! /bin/sh
=======
#! /bin/bash -e
>>>>>>> develop
# -*- coding: utf-8 -*-

if [ -n "$1" ]; then
    cd $1
fi

dpkg -i ./*.deb
