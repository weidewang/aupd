#!/bin/bash
ulimit -n 655350

export GOMAXPROCS=`cat /proc/cpuinfo | grep -v grep | grep process | wc -l`
export GOGCTRACE=1
export GODEBUG=gctrace=1

GODEBUG=gctrace=1 GOGCTRACE=1 GOMAXPROCS=8 /srv/aupd/bin/aupd 