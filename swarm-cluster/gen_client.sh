#!/bin/bash

name="client-$(cat /dev/urandom | tr -d -c a-z | head -c 8)"
private_ip="$1"

if [ ! -f "leader_config" ]; then
    echo "[-] leader_config is not found"
    exit 1
fi
source "leader_config"


mkdir "${name}"
sed -e "s/<hostname>/${name}/" \
    -e "s/<private_ip>/${private_ip}/" \
    Vagrantfile  > "${name}/Vagrantfile"
cp -r ansible/ "${name}/ansible"
cp hosts ansible.cfg "${name}/ansible"

(cd "${name}";
vagrant up;
vagrant ssh-config > ansible/ssh_config
)

(cd "${name}/ansible";
ansible-playbook -i hosts main.yml -e "private_ip=${private_ip}" -e "name=${name}"
)

(cd "${name}";
vagrant ssh -c "${JOIN_COMMAND}")

echo "[+] ${name}"
