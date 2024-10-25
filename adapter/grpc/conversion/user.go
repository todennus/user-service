package conversion

import (
	pbdto "github.com/todennus/proto/gen/service/dto"
	pbresource "github.com/todennus/proto/gen/service/dto/resource"
	ucdto "github.com/todennus/user-service/usecase/dto"
	ucresource "github.com/todennus/user-service/usecase/dto/resource"
	"github.com/xybor-x/snowflake"
)

func NewPbUser(user *ucresource.User) *pbresource.User {
	return &pbresource.User{
		Id:          user.ID.Int64(),
		Username:    user.Username,
		DisplayName: user.DisplayName,
		Role:        user.Role.String(),
	}
}

func NewUsecaseUserValidateRequest(req *pbdto.UserValidateRequest) *ucdto.UserValidateCredentialsRequest {
	return &ucdto.UserValidateCredentialsRequest{
		Username: req.Username,
		Password: req.Password,
	}
}

func NewPbUserValidateResponse(resp *ucdto.UserValidateCredentialsResponse) *pbdto.UserValidateResponse {
	if resp == nil {
		return nil
	}

	return &pbdto.UserValidateResponse{
		User: NewPbUser(resp.User),
	}
}

func NewUsecaseUserGetByIDRequest(req *pbdto.UserGetByIDRequest) *ucdto.UserGetByIDRequest {
	return &ucdto.UserGetByIDRequest{
		UserID: snowflake.ID(req.Id),
	}
}

func NewPbUserGetByIDResponse(resp *ucdto.UserGetByIDResponse) *pbdto.UserGetByIDResponse {
	if resp == nil {
		return nil
	}

	return &pbdto.UserGetByIDResponse{
		User: NewPbUser(resp.User),
	}
}