#! /bin/bash
# -*- coding: utf-8 -*-

../../bin/ffctl mkpb make all --verbose --overwrite --config-file ./mkpb-sample.yaml

echo ""
echo "next:"
echo "ansible-playbook -i lxd-mic.inv -K lxd-mic.yaml"
