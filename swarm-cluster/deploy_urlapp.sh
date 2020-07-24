#!/bin/bash

./build-image.py --registry 192.168.33.254:5000 --project urlapp/

docker_compose=$(cd urlapp; find -name "docker-compose_*.yml" -printf "%P\n" | head -n 1)
batch=$(mktemp)
cat <<EOF > $batch
put urlapp/${docker_compose}
EOF

d=$(find -name "leader-*")
sftp -b "${batch}" -F "${d}/ansible/ssh_config" default

(cd "$d";
vagrant ssh -c "docker stack deploy --compose-file ${docker_compose} urlapp"
)
