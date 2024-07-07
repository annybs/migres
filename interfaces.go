package migres

// Migration is anything that can upgrade or downgrade external state - commonly, but not limited to, database schemas.
//
// Each migration SHOULD be able to upgrade or downgrade freely, allowing any changes to be reverted with ease.
//
// Of course, this is not always possible.
// In the case of irreversible state change, the opposite function should return an error e.g. if an upgrade deletes something irrecoverably, have the corresponding downgrade function throw a descriptive error.
type Migration interface {
	Downgrade() error // Perform a downgrade.
	Upgrade() error   // Perform an upgrade.
}
