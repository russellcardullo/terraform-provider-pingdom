package pingdom

import "errors"

var ErrMissingId = errors.New("required field 'Id' missing")
var ErrBadResolution = errors.New("resolution must be either 'hour', 'day' or 'week'")
