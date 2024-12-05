package dto

import "github.com/xybor-x/snowflake"

type AvatarGetUploadTokenRequest struct {
	UserID snowflake.ID
}

type AvatarGetUploadTokenResponse struct {
	UploadToken string
}

func NewAvatarGetUploadTokenResponse(uploadToken string) *AvatarGetUploadTokenResponse {
	return &AvatarGetUploadTokenResponse{UploadToken: uploadToken}
}

type AvatarUpdateRequest struct {
	UserID    snowflake.ID
	FileToken string
}

type AvatarUpdateResponse struct {
}

func NewAvatarUpdateResponse() *AvatarUpdateResponse {
	return &AvatarUpdateResponse{}
}
