# Seata Plugin for Lynx

`lynx-seata` is a thin Lynx integration for `seata-go`. Its runtime scope is intentionally small:

- decide whether the Seata client should start
- pass an external Seata config file path into `client.InitPath(...)`
- expose `WithGlobalTx(...)` so service code can mark transaction boundaries explicitly

It does not introduce a Lynx-owned YAML surface for registry, coordinator, auth, Saga, TCC, XA, or transport tuning. Those settings stay in the external Seata config file referenced by `config_file_path`.

## Runtime facts

- Go module: `github.com/go-lynx/lynx-seata`
- Config prefix: `lynx.seata`
- Runtime plugin name: `seata.server`
- Public APIs: `GetPlugin()`, `GetConfig()`, `GetConfigFilePath()`, `IsEnabled()`, `WithGlobalTx(...)`, `GetMetrics()`

## Configuration

The shipped example in [`conf/example_config.yml`](./conf/example_config.yml) is the authoritative runtime shape:

```yaml
lynx:
  seata:
    enabled: true
    config_file_path: "./conf/seata.yml"
```

Field behavior:

- `enabled`: gates Seata startup and resource registration
- `config_file_path`: filesystem path passed to `client.InitPath(...)`

## Usage

```go
import (
	"context"
	"time"

	seata "github.com/go-lynx/lynx-seata"
)

func createOrder(ctx context.Context) error {
	plugin := seata.GetPlugin()
	if plugin == nil || !plugin.IsEnabled() {
		return nil
	}

	return plugin.WithGlobalTx(ctx, "CreateOrderTx", 30*time.Second, func(ctx context.Context) error {
		return doBusiness(ctx)
	})
}
```

## Operational notes

- Keep all real Seata topology and client settings in the external file pointed to by `config_file_path`.
- `WithGlobalTx(...)` is the intended API for business-level transaction boundaries.
- `client.InitPath(...)` is now guarded as a process-wide one-time initialization. Repeated startup with the same config path is allowed; switching to a different config path in the same process returns an explicit error.
- `CheckHealth()` only validates that the plugin is enabled and the path is non-empty. It is not a full connectivity probe for the Seata control plane.
