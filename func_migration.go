package migres

// FuncMigration enables creating a functional migration using callback functions.
//
//	import "github.com/annybs/migres"
//
//	type MyBackend struct{}
//
//	func (mb *MyBackend) Module() migres.Module {
//	  return migres.Module{
//	    "1.0.0": migres.Func(mb.upgradeV1, mb.downgradeV1),
//	    "2.0.0": migres.Func(mb.upgradeV2, mb.downgradeV2),
//	  }
//	}
//
// In this example MyModule can be defined with multiple upgrade/downgrade functions.
// This may be simpler than defining separate migration structs in many cases.
type FuncMigration struct {
	D func() error
	U func() error
}

func (fm *FuncMigration) Downgrade() error {
	return fm.D()
}

func (fm *FuncMigration) Upgrade() error {
	return fm.U()
}

// Func creates a functional migration.
func Func(up, down func() error) *FuncMigration {
	return &FuncMigration{
		D: down,
		U: up,
	}
}
