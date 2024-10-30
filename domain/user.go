package domain

import (
	"fmt"
	"time"

	"github.com/todennus/shared/enumdef"
	"github.com/todennus/x/enum"
	"github.com/todennus/x/xstring"
	"github.com/xybor-x/snowflake"
)

const (
	MinimumDisplayNameLength = 3
	MaximumDisplayNameLength = 32

	MinimumUsernameLength = 4
	MaximumUsernameLength = 20

	MinimumPasswordLength = 8
	MaximumPassowrdLength = 32
)

type User struct {
	ID          snowflake.ID
	DisplayName string
	Username    string
	HashedPass  string
	Role        enum.Enum[enumdef.UserRole]
	UpdatedAt   time.Time
}

type UserDomain struct {
	Snowflake *snowflake.Node
}

func NewUserDomain(snowflake *snowflake.Node) (*UserDomain, error) {
	return &UserDomain{Snowflake: snowflake}, nil
}

func (domain *UserDomain) New(username, password string) (*User, error) {
	if err := domain.validateUsername(username); err != nil {
		return nil, err
	}

	if err := domain.validatePassword(password); err != nil {
		return nil, err
	}

	hashedPass, err := HashPassword(password)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:          domain.Snowflake.Generate(),
		DisplayName: username,
		Username:    username,
		HashedPass:  string(hashedPass),
		Role:        enumdef.UserRoleUser,
	}, nil
}

func (domain *UserDomain) NewFirst(username, password string) (*User, error) {
	user, err := domain.New(username, password)
	if err != nil {
		return nil, err
	}

	user.Role = enumdef.UserRoleAdmin
	return user, nil
}

func (domain *UserDomain) Validate(hashedPassword, password string) error {
	return ValidatePassword(hashedPassword, password)
}

func (domain *UserDomain) SetDisplayName(user *User, displayname string) error {
	if err := domain.validateDisplayName(displayname); err != nil {
		return err
	}

	user.DisplayName = displayname
	return nil
}

func (domain *UserDomain) validateDisplayName(displayname string) error {
	if len(displayname) > MaximumDisplayNameLength {
		return fmt.Errorf("%w: require at most %d characters", ErrDisplayNameInvalid, MaximumDisplayNameLength)
	}

	if len(displayname) < MinimumDisplayNameLength {
		return fmt.Errorf("%w: require at least %d characters", ErrDisplayNameInvalid, MinimumDisplayNameLength)
	}

	for _, c := range displayname {
		if !xstring.IsNumber(c) && !xstring.IsLetter(c) && !xstring.IsUnderscore(c) && !xstring.IsSpace(c) {
			return fmt.Errorf("%w: got an invalid character %c", ErrUsernameInvalid, c)
		}
	}

	return nil
}

func (domain *UserDomain) validateUsername(username string) error {
	if len(username) > MaximumUsernameLength {
		return fmt.Errorf("%w: require at most %d characters", ErrUsernameInvalid, MaximumUsernameLength)
	}

	if len(username) < MinimumUsernameLength {
		return fmt.Errorf("%w: require at least %d characters", ErrUsernameInvalid, MinimumUsernameLength)
	}

	for _, c := range username {
		if !xstring.IsNumber(c) && !xstring.IsLetter(c) && !xstring.IsUnderscore(c) {
			return fmt.Errorf("%w: got an invalid character %c", ErrUsernameInvalid, c)
		}
	}

	return nil
}

func (domain *UserDomain) validatePassword(password string) error {
	if len(password) > MaximumPassowrdLength {
		return fmt.Errorf("%w: require at most %d characters", ErrPasswordInvalid, MaximumPassowrdLength)
	}

	if len(password) < MinimumPasswordLength {
		return fmt.Errorf("%w: require at least %d characters", ErrPasswordInvalid, MinimumPasswordLength)
	}

	haveLowercase := false
	haveUppercase := false
	haveNumber := false
	haveSpecial := false

	for _, c := range password {
		switch {
		case xstring.IsLowerCaseLetter(c):
			haveLowercase = true
		case xstring.IsUpperCaseLetter(c):
			haveUppercase = true
		case xstring.IsNumber(c):
			haveNumber = true
		case xstring.IsSpecialCharacter(c):
			haveSpecial = true
		default:
			return fmt.Errorf("%w: got an invalid character %c", ErrPasswordInvalid, c)
		}
	}

	if !haveLowercase {
		return fmt.Errorf("%w: require at least a lowercase letter", ErrPasswordInvalid)
	}

	if !haveUppercase {
		return fmt.Errorf("%w: require at least an uppercase letter", ErrPasswordInvalid)
	}

	if !haveNumber {
		return fmt.Errorf("%w: require at least a number", ErrPasswordInvalid)
	}

	if !haveSpecial {
		return fmt.Errorf("%w: require at least a special character", ErrPasswordInvalid)
	}

	return nil
}
