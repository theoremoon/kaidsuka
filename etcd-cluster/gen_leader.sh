#!/bin/bash

name="leader-$(cat /dev/urandom | tr -d -c a-z | head -c 8)"
private_ip="$1"

echo "LEADER_NAME=${name}" > leader_config
echo "LEADER_IP=${private_ip}" >> leader_config

echo "[+] ${name}"

mkdir "${name}"
sed -e "s/<hostname>/${name}/" \
    -e "s/<private_ip>/${private_ip}/" \
    Vagrantfile  > "${name}/Vagrantfile"

cp -r ansible/leader "${name}/ansible"
cp hosts ansible.cfg "${name}/ansible"

(cd "${name}";
vagrant up;
vagrant ssh-config > ansible/ssh_config
)

(cd "${name}/ansible";
ansible-playbook -i hosts main.yml -e "private_ip=${private_ip}" -e "name=${name}"
)



