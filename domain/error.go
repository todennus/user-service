package domain

import (
	"fmt"

	"github.com/todennus/shared/errordef"
)

var (
	ErrUsernameInvalid    = fmt.Errorf("%winvalid username", errordef.ErrDomainKnown)
	ErrDisplayNameInvalid = fmt.Errorf("%winvalid display name", errordef.ErrDomainKnown)
	ErrPasswordInvalid    = fmt.Errorf("%winvalid password", errordef.ErrDomainKnown)
	ErrMismatchedPassword = fmt.Errorf("%wmismatched password", errordef.ErrDomainKnown)
)
