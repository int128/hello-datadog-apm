# hello-datadog-apm [![go](https://github.com/int128/hello-datadog-apm/actions/workflows/go.yaml/badge.svg)](https://github.com/int128/hello-datadog-apm/actions/workflows/go.yaml)

This is an example application using Datadog APM.

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

## Run in GitHub Actions

See [the workflow](.github/workflows/go.yaml).

## Amazon ECS Fargate Task

See [the official doc](https://docs.datadoghq.com/integrations/ecs_fargate/?tab=cloudformation).

```json
{
  "containerDefinitions": [
    {
      "name": "main",
      "image": "ghcr.io/int128/hello-datadog-apm/go:main",
      "essential": true,
      "environment": [
        {
          "name": "DD_ENV",
          "value": "int128-sandbox"
        }
      ],
      "logConfiguration": {
        "logDriver": "awslogs",
        "options": {
          "awslogs-create-group": "true",
          "awslogs-group": "/ecs/int128-hello-datadog-apm",
          "awslogs-region": "ap-northeast-1",
          "awslogs-stream-prefix": "ecs"
        }
      }
    },
    {
      "name": "datadog-agent",
      "image": "public.ecr.aws/datadog/agent:7",
      "essential": true,
      "environment": [
        {
          "name": "ECS_FARGATE",
          "value": "true"
        },
        {
          "name": "DD_API_KEY",
          "valueFrom": "arn:aws:secretsmanager:ap-northeast-1:123456789012:secret:int128-hello-datadog-apm-abcdef:datadog_api_key::"
        }
      ],
      "logConfiguration": {
        "logDriver": "awslogs",
        "options": {
          "awslogs-create-group": "true",
          "awslogs-group": "/ecs/int128-hello-datadog-apm",
          "awslogs-region": "ap-northeast-1",
          "awslogs-stream-prefix": "ecs"
        }
      }
    }
  ],
  "networkMode": "awsvpc",
  "requiresCompatibilities": ["FARGATE"],
  "cpu": "256",
  "memory": "512"
}
```

The following log should be written:

```
2023/04/01 03:11:49 Datadog Tracer v1.48.0 INFO: DATADOG TRACER CONFIGURATION
```

```json
{
    "date": "2023-04-01T03:11:49Z",
    "os_name": "Linux (Unknown Distribution)",
    "os_version": "Debian GNU/Linux 11 (bullseye)",
    "version": "v1.48.0",
    "lang": "Go",
    "lang_version": "go1.20.2",
    "env": "int128-sandbox",
    "service": "hello-datadog-apm",
    "agent_url": "http://localhost:8126/v0.4/traces",
    "agent_error": "",
    "debug": false,
    "analytics_enabled": false,
    "sample_rate": "NaN",
    "sample_rate_limit": "disabled",
    "sampling_rules": null,
    "sampling_rules_error": "",
    "service_mappings": null,
    "tags": {
        "runtime-id": "********"
    },
    "runtime_metrics_enabled": false,
    "health_metrics_enabled": false,
    "profiler_code_hotspots_enabled": true,
    "profiler_endpoints_enabled": true,
    "dd_version": "",
    "architecture": "amd64",
    "global_service": "hello-datadog-apm",
    "lambda_mode": "false",
    "appsec": false,
    "agent_features": {
        "DropP0s": true,
        "Stats": true,
        "StatsdPort": 0
    }
}
```
