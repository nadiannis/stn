package models

import "errors"

var ErrDuplicateEmail = errors.New("models: duplicate email")
var ErrDuplicateBackHalf = errors.New("models: duplicate back-half")
var ErrInvalidCredentials = errors.New("models: invalid credentials")
var ErrNoRecord = errors.New("models: no matching record found")
