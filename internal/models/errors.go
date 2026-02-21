package models

import "errors"

// ErrNoRecord is returned when a database query does not find a matching record.
var ErrNoRecord = errors.New("models: no matching record found.")
