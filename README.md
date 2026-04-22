# cvt-plugin-registry

Placeholder CVT `RegistryProvider` plugin. Returns a hardcoded inline
OpenAPI spec from `FetchSchema` and logs `RegisterConsumerUsage` calls
to stderr.

This repo exists to unblock end-to-end testing of the CVT CLI wiring
landing in [sahina/cvt#83](https://github.com/sahina/cvt/issues/83).
The real central-registry HTTP contract is not yet defined; once it
is, this plugin will be replaced with a proper adapter.

## Build

```sh
go build -o cvt-plugin-registry .
cvt plugins install ./cvt-plugin-registry
```

## Configure

Add to `~/.cvt/config.yaml`:

```yaml
plugins:
  registry:
    binary: ~/.cvt/plugins/cvt-plugin-registry
    timeout: 5s
    on_error: fail_closed
hooks:
  fetch_schema: registry
  register_consumer_usage: registry
```

## Try it

```sh
cvt validate --schema hello --consumer demo --interaction interaction.json
```

Any `--schema <id>` value resolves to the same inline spec. Consumer
usage records are printed to CVT's structured log.

## License

MIT.
