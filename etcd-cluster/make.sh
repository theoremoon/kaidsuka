#!/bin/bash

if [ "$1" = "build" ]; then
    ./gen_leader.sh 192.168.33.100
    ./gen_client.sh 192.168.33.101
    ./gen_client.sh 192.168.33.102
    ./gen_client.sh 192.168.33.103

elif [ "$1" = "clean" ]; then
    for d in $(find -name "client-*"); do
        (cd "$d"; vagrant destroy --force)
        rm -rf "$d"
    done

    d=$(find -name "leader-*")
    (cd "$d"; vagrant destroy --force)
    rm -rf "$d"
    rm leader_config

fi
