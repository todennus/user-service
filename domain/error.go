package domain

import (
	"fmt"

	"github.com/todennus/shared/errordef"
)

var (
	ErrUsernameInvalid    = fmt.Errorf("%w%s", errordef.ErrDomainKnown, "invalid username")
	ErrDisplayNameInvalid = fmt.Errorf("%w%s", errordef.ErrDomainKnown, "invalid display name")
	ErrPasswordInvalid    = fmt.Errorf("%w%s", errordef.ErrDomainKnown, "invalid password")

	ErrMismatchedPassword = fmt.Errorf("%w%s", errordef.ErrDomainKnown, "mismatched password")
)
