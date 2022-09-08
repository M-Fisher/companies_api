package base

import "errors"

const ServerErrorCode = 500

var ErrUnauthorized = errors.New(`not authorized`)
