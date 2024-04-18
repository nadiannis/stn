package models

import "errors"

var ErrDuplicateEmail = errors.New("models: duplicate email")
var ErrInvalidCredentials = errors.New("models: invalid credentials")
var ErrNoRecord = errors.New("models: no matching record found")
