#!/bin/bash

name="client-$(cat /dev/urandom | tr -d -c a-z | head -c 8)"
private_ip="$1"

if [ ! -f "leader_config" ]; then
    echo "[-] leader_config is not found"
    exit 1
fi
source "leader_config"

echo "[+] ${name}"

mkdir "${name}"
sed -e "s/<hostname>/${name}/" \
    -e "s/<private_ip>/${private_ip}/" \
    Vagrantfile  > "${name}/Vagrantfile"
cp -r ansible/client "${name}/ansible"
cp hosts ansible.cfg "${name}/ansible"

(cd "${name}";
vagrant up;
vagrant ssh-config > ansible/ssh_config
)

ETCD_TEMP=$(mktemp)
(cd "${LEADER_NAME}";
vagrant ssh -c "etcdctl member add ${name} --peer-urls 'http://${private_ip}:2380'") > "${ETCD_TEMP}"
sed -i 1,2d "${ETCD_TEMP}"
source "${ETCD_TEMP}"

(cd "${name}/ansible";
ansible-playbook -i hosts main.yml -e "private_ip=${private_ip}" -e "name=${name}" -e "initial_cluster=${ETCD_INITIAL_CLUSTER}"
)

