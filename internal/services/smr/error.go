package smr

import "errors"

var (
	ErrNoLink = errors.New("no link")
	ErrParse  = errors.New("parse error")
	ErrScheme = errors.New("scheme error")
)
