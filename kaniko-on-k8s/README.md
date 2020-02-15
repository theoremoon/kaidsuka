```
$ kubectl apply -f registry.yaml
$ kubectl apply -f kaniko.yaml
$ docker run --rm -p8888:80 192.168.99.102:32000/theoremoon/hello:latest
```

Note: [just\_login](https://github.com/theoremoon/ctf4y/tree/master/1/just_login) assumes it use with mounted directory so this is do not work well
