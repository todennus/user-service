package dto

import (
	"github.com/todennus/shared/errordef"
	"github.com/todennus/user-service/adapter/rest/dto/resource"
	"github.com/todennus/user-service/usecase/dto"
	"github.com/todennus/x/xerror"
	"github.com/xybor-x/snowflake"
)

func ParseUserID(meID snowflake.ID, s string) (snowflake.ID, error) {
	if s == "@me" {
		return meID, nil
	}

	return snowflake.ParseString(s)
}

// Register
type UserRegisterRequest struct {
	Username string `json:"username" example:"huykingsofm"`
	Password string `json:"password" example:"s3Cr3tP@ssW0rD"`
}

func (req UserRegisterRequest) To() *dto.UserRegisterRequest {
	return &dto.UserRegisterRequest{
		Username: req.Username,
		Password: req.Password,
	}
}

type UserRegisterResponse struct {
	*resource.User
}

func NewUserRegisterResponse(resp *dto.UserRegisterResponse) *UserRegisterResponse {
	if resp == nil {
		return nil
	}

	return &UserRegisterResponse{
		User: resource.NewUser(resp.User),
	}
}

// Register First
type UserRegisterFirstRequest struct {
	Username string `json:"username" example:"huykingsofm"`
	Password string `json:"password" example:"s3Cr3tP@ssW0rD"`
}

func (req UserRegisterFirstRequest) To() *dto.UserRegisterFirstRequest {
	return &dto.UserRegisterFirstRequest{
		Username: req.Username,
		Password: req.Password,
	}
}

type UserRegisterFirstResponse struct {
	*resource.User
}

func NewUserRegisterFirstResponse(resp *dto.UserRegisterFirstResponse) *UserRegisterFirstResponse {
	if resp == nil {
		return nil
	}

	return &UserRegisterFirstResponse{
		User: resource.NewUser(resp.User),
	}
}

// GetByID
type UserGetByIDRequest struct {
	UserID string `param:"user_id"`
}

func (req UserGetByIDRequest) To(meID snowflake.ID) (*dto.UserGetByIDRequest, error) {
	userID, err := ParseUserID(meID, req.UserID)
	if err != nil {
		return nil, xerror.Enrich(errordef.ErrRequestInvalid, "user id is invalid").
			Hide(err, "failed-to-parse-user-id", "uid", req.UserID)
	}

	return &dto.UserGetByIDRequest{UserID: userID}, nil
}

type UserGetByIDResponse struct {
	*resource.User
}

func NewUserGetByIDResponse(resp *dto.UserGetByIDResponse) *UserGetByIDResponse {
	if resp == nil {
		return nil
	}

	return &UserGetByIDResponse{
		User: resource.NewUser(resp.User),
	}
}

// GetByUsername
type UserGetByUsernameRequest struct {
	Username string `param:"username"`
}

func (req UserGetByUsernameRequest) To() *dto.UserGetByUsernameRequest {
	return &dto.UserGetByUsernameRequest{
		Username: req.Username,
	}
}

type UserGetByUsernameResponse struct {
	*resource.User
}

func NewUserGetByUsernameResponse(resp *dto.UserGetByUsernameResponse) *UserGetByUsernameResponse {
	if resp == nil {
		return nil
	}

	return &UserGetByUsernameResponse{
		User: resource.NewUser(resp.User),
	}
}

// Validate
type UserValidateRequest struct {
	Username string `json:"username" example:"huykingsofm"`
	Password string `json:"password" example:"s3Cr3tP@ssW0rD"`
}

func (req UserValidateRequest) To() *dto.UserValidateCredentialsRequest {
	return &dto.UserValidateCredentialsRequest{
		Username: req.Username,
		Password: req.Password,
	}
}

type UserValidateResponse struct {
	*resource.User
}

func NewUserValidateResponse(resp *dto.UserValidateCredentialsResponse) *UserValidateResponse {
	if resp == nil {
		return nil
	}

	return &UserValidateResponse{
		User: resource.NewUser(resp.User),
	}
}

type AvatarGetPolicyTokenRequest struct{}

func (req *AvatarGetPolicyTokenRequest) To() *dto.AvatarGetPolicyTokenRequest {
	return &dto.AvatarGetPolicyTokenRequest{}
}

type AvatarGetPolicyTokenResponse struct {
	PolicyToken string `json:"policy_token"`
}

func NewAvatarGetPolicyTokenResponse(resp *dto.AvatarGetPolicyTokenResponse) *AvatarGetPolicyTokenResponse {
	if resp == nil {
		return nil
	}

	return &AvatarGetPolicyTokenResponse{
		PolicyToken: resp.PolicyToken,
	}
}

type AvatarUpdateRequest struct {
	TemporaryFileToken string `json:"temporary_file_token"`
}

func (req *AvatarUpdateRequest) To() *dto.AvatarUpdateRequest {
	return &dto.AvatarUpdateRequest{
		TemporaryFileToken: req.TemporaryFileToken,
	}
}

type AvatarUpdateResponse struct {
	AvatarURL string `json:"avatar_url"`
}

func NewAvatarUpdateResponse(resp *dto.AvatarUpdateResponse) *AvatarUpdateResponse {
	if resp == nil {
		return nil
	}

	return &AvatarUpdateResponse{
		AvatarURL: resp.AvatarURL,
	}
}
