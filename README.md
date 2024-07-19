# Migres

This package provides simple migration capabilities for any backend.

## System requirements

- [Go v1.21](https://go.dev/dl/)

## Basic usage

The key type in this package is `Module` which allows mapping [version strings](https://pkg.go.dev/github.com/annybs/go-version) to `Migration` interfaces. For example:

```go
import "github.com/annybs/migres"

type MyBackend struct{}

func (mb *MyBackend) Module() migres.Module {
  return migres.Module{
    "1.0.0": migres.Func(mb.upgradeV1, mb.downgradeV1),
    "2.0.0": migres.Func(mb.upgradeV2, mb.downgradeV2),
  }
}
```

Call `Module.Upgrade(from, to)` or `Module.Downgrade(from, to)` in order to execute migrations. The module ensures migrations are all run in the correct order.

## License

See [LICENSE.md](./LICENSE.md)
