#!/bin/bash

if [ "$1" = "build" ]; then

    ./gen_leader.sh 192.168.33.100

    ./gen_client.sh 192.168.33.101
    ./gen_client.sh 192.168.33.102
    ./gen_client.sh 192.168.33.103

elif [ "$1" = "registry" ]; then
    cp -r ansible/ registry/
    cp hosts ansible.cfg registry/ansible
    (cd registry;
    vagrant up;
    vagrant ssh-config > ansible/ssh_config)
    (cd registry/ansible;
    ansible-playbook -i hosts main.yml)
    (cd registry;
    vagrant ssh -c "docker run -d -p 5000:5000 registry")

elif [ "$1" = "clean" ]; then
    for d in $(find -name "client-*"); do
        (cd "$d"; vagrant destroy --force)
        rm -rf "$d"
    done

    d=$(find -name "leader-*")
    (cd "$d"; vagrant destroy --force)
    rm -rf "$d"
    rm leader_config

    (cd registry;
    vagrant destroy --force)

    (cd urlapp;
    find -name "docker-compose_*.yml" -exec rm \{\} \; )
fi
