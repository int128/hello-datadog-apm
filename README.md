# hello-datadog-apm
Hello world with Datadog APM.

## Local development

To start [datadog-agent](https://docs.datadoghq.com/containers/docker/) locally:

```sh
docker run --rm --name dd-agent \
  -v /var/run/docker.sock:/var/run/docker.sock:ro \
  -v /proc/:/host/proc/:ro \
  -v /sys/fs/cgroup/:/host/sys/fs/cgroup:ro \
  -e DD_API_KEY= \
  -e DD_APM_NON_LOCAL_TRAFFIC=true \
  -p 127.0.0.1:8126:8126/tcp \
  gcr.io/datadoghq/agent:7
```
