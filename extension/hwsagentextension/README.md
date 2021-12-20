# hwsagent Extension

The **Hardware Sentry Agent** is the internal component which scrapes targets, collects metrics and pushes OTLP data to the OTLP receiver of the *OpenTelemetry Collector*. The `hws_agent` extension starts the **Hardware Sentry Agent** as a child process of the *OpenTelemetry Collector*.

```yaml
  hws_agent:
    grpc: http://localhost:4317
```

The `hws_agent` extension checks that **Hardware Sentry Agent** is up and running and restarts its process if needed.

The example above shows how to configure **Hardware Sentry Agent** to push metrics to the local _OTLP receiver_ using [gRPC](https://grpc.io/) on port **TCP/4317**.
If your OTLP receiver runs on another host or uses a different protocol or port, you will need to update the `grpc` value. Format: `<http|https>://<host>:<port>`.

By default, the **Hardware Sentry Agent**'s configuration file is **config/hws-config.yaml**. You can provide an alternate configuration file using the `--config` argument in the `extra_args`.

```yaml
  hws_agent:
    grpc: http://localhost:4317
    extra_args:
      - --config=config\hws-config-2.yaml
```

The following settings can be optionally configured:

- `extra_args`: The Hardware Sentry Agent arguments such us `--config=config\hws-config-2.yaml`.
- `restart_delay`: Time to wait before restarting the Hardware Sentry Agent.
- `retries`: Number of restarts due to failures.

## Example Config

```yaml
extensions:
  health_check:
  hws_agent:
    grpc: http://localhost:4317

receivers:
# {...}
service:

  extensions: [hws_agent]
```
