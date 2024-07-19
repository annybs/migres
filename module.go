package migres

import (
	"slices"
	"sort"

	"github.com/annybs/go-version"
)

// Module provides migrations keyed by version string.
// This helps to organise migrations execute upgrades or downgrades in the correct version order.
//
// For example:
//
//	mod := Module{"1": migration1, "3": migration3, "2": migration2}
//	mod.Upgrade("1", "3")
//
// Here, the module's three migrations will be sorted into the proper order of 1-2-3, then their upgrades will be performed in that same order.
type Module map[string]Migration

// Downgrade migrations at versions after from until (and including) to.
func (mod Module) Downgrade(from, to string) error {
	fromVersion, err := version.Parse(from)
	if err != nil {
		return err
	}

	toVersion, err := version.Parse(to)
	if err != nil {
		return err
	}

	versions, err := mod.Versions()
	if err != nil {
		return err
	}

	versions = versions.Match(&version.Constraint{
		Lte: fromVersion,
		Gt:  toVersion,
	})
	sort.Stable(versions)
	slices.Reverse(versions)

	var lastVersion *version.Version
	for _, version := range versions {
		m := mod[version.Text]
		if err := m.Downgrade(); err != nil {
			return failMigration(err, version, lastVersion)
		}
		lastVersion = version
	}

	return nil
}

// Upgrade migrations at versions after from until (and including) to.
func (mod Module) Upgrade(from, to string) error {
	fromVersion, err := version.Parse(from)
	if err != nil {
		return err
	}

	toVersion, err := version.Parse(to)
	if err != nil {
		return err
	}

	versions, err := mod.Versions()
	if err != nil {
		return err
	}

	versions = versions.Match(&version.Constraint{
		Gt:  fromVersion,
		Lte: toVersion,
	})
	sort.Stable(versions)

	var lastVersion *version.Version
	for _, version := range versions {
		m := mod[version.Text]
		if err := m.Upgrade(); err != nil {
			return failMigration(err, version, lastVersion)
		}
		lastVersion = version
	}

	return nil
}

// Versions gets a list of all versions in the module.
func (mod Module) Versions() (version.List, error) {
	list := version.List{}

	for str := range mod {
		v, err := version.Parse(str)
		if err != nil {
			return nil, err
		}
		list = append(list, v)
	}

	return list, nil
}
