package domain

import "errors"

var (
	ErrUserNotFound         = errors.New("user not found")
	ErrUserAlreadyActive    = errors.New("user already active")
	ErrUserAlreadyInactive  = errors.New("user already inactive")
	ErrUserStatusInactive   = errors.New("user status is inactive")
	ErrUserExistsActive     = errors.New("user already exists and active")
	ErrUserExistsInactive   = errors.New("user already exists and inactive")
	ErrNothingToChange      = errors.New("nothing to change")
	ErrNoUsers              = errors.New("no users exists")
	ErrTooManyArguments     = errors.New("too many arguments")
	ErrMissingLoginAndName  = errors.New("no <login> <name> provided")
	ErrMissingName          = errors.New("no <name> provided")
	ErrMissingID            = errors.New("no user <id> provided")
	ErrMissingNewUserFields = errors.New("no user <id>, new <login>, new <name> provided")
	ErrMissingNewLoginName  = errors.New("no user new <login>, new <name> provided")
	ErrMissingNewName       = errors.New("no user new <name> provided")
	ErrMissingActiveUserID  = errors.New("no new active user <id> provided")
)
