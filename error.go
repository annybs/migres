package migres

import (
	"fmt"

	"github.com/annybs/go-version"
)

// Migration error.
var (
	ErrMigrationFailed = Error{Message: "migration failed at version %q: %s"}
)

// Error reflects an error that occurred during a migration.
type Error struct {
	Message       string           // Error message template.
	PreviousError error            // Original error encountered during the migration.
	Version       *version.Version // Version at which the migration error occured.
	LastVersion   *version.Version // Version
}

// Error retrieves the message of a migration error.
// This does not necessarily include all information about the error, such as the last
func (e *Error) Error() string {
	return fmt.Sprintf(e.Message, e.Version, e.PreviousError)
}

// Is determines whether the Error is an instance of the target.
// https://pkg.go.dev/errors#Is
//
// This implementation does not compare versions.
func (e *Error) Is(target error) bool {
	if t, ok := target.(*Error); ok {
		return t.Message == e.Message
	}
	return false
}

func failMigration(err error, v, last *version.Version) *Error {
	return &Error{
		Message:       ErrMigrationFailed.Message,
		PreviousError: err,
		Version:       v,
		LastVersion:   last,
	}
}
