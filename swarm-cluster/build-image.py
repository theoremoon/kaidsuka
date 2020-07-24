#!/usr/bin/env python3
import yaml
import os
import sys
import random
import string
import subprocess

# parse commandline
def show_help():
  print("--registry  Docker registry url")
  print("--project   Project directory")

i = 1
project = None
registry = None
while i < len(sys.argv):
  if sys.argv[i] == "--registry":
    registry = sys.argv[i+1]
    i += 2
  elif sys.argv[i] == "--project":
    project = sys.argv[i+1].rstrip("/")
    i += 2
  else:
    break

if project is None or registry is None:
  show_help()
  exit(1)

# load docker-compose.yml
compose_file = os.path.join(project, "docker-compose.yml")
if not os.path.exists(compose_file):
  print("[-] docker-compose.yml does not found")
  exit(1)

with open(compose_file) as f:
  compose = yaml.safe_load(f)

# set image url
for service in compose['services'].keys():
  if "image" in service:
    continue

  compose["services"][service]["image"] = "{}/{}_{}:latest".format(registry, project, service)

# save as new compose file
new_compose_name = "docker-compose_{}.yml".format("".join(random.choices(string.ascii_letters, k=8)))
new_compose_file = os.path.join(project, new_compose_name)
with open(new_compose_file, "w") as f:
  f.write(yaml.dump(compose, default_flow_style=False))

# build & push iamge using new compose file
subprocess.run(["docker-compose", "-f", new_compose_name, "build"], cwd=project, stdout=sys.stdout, stderr=sys.stderr)
subprocess.run(["docker-compose", "-f", new_compose_name, "push"], cwd=project, stdout=sys.stdout, stderr=sys.stderr)

