package pingdom

import "errors"

// ErrMissingId is an error for when a required Id field is missing.
var ErrMissingId = errors.New("required field 'Id' missing")

// ErrBadResolution is an error for when an invalid resolution is specified.
var ErrBadResolution = errors.New("resolution must be either 'hour', 'day' or 'week'")
