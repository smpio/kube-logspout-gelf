# Graylog GELF Module for Logspout adopted for Kubernetes

This module allows Logspout to send Docker logs in the GELF format to Graylog via UDP.

## Build
To build, you'll need to fork [Logspout](https://github.com/gliderlabs/logspout), add the following code to `modules.go` 

```
_ "github.com/smpio/kube-logspout-gelf"
```
and run `docker build -t $(whoami)/logspout:gelf`

## Run

```
docker run \
    -v /var/run/docker.sock:/var/run/docker.sock \
    $(whoami)/logspout:gelf \
    gelf://<graylog_host>:12201

```

## A note about GELF parameters
The following docker container attributes are mapped to the corresponding GELF extra attributes.

```
{
        "_kube_namespace": <namespace>,
        "_kube_container": <container-name>,
        "host":            <pod-name>,
}
```

You can also add extra custom fields by setting env vars with prefix `KUBE_`.

For example by setting `KUBE_NODE=node1` will add the extra field `_kube_node=node1`.



## License
MIT. See [License](LICENSE)
