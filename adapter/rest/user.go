package rest

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/todennus/shared/errordef"
	"github.com/todennus/shared/middleware"
	"github.com/todennus/shared/response"
	"github.com/todennus/shared/xcontext"
	"github.com/todennus/user-service/adapter/abstraction"
	"github.com/todennus/user-service/adapter/rest/dto"
	"github.com/todennus/x/xhttp"
)

type UserAdapter struct {
	userUsecase   abstraction.UserUsecase
	avatarUsecase abstraction.AvatarUsecase
}

func NewUserAdapter(userUsecase abstraction.UserUsecase, avatarUsecase abstraction.AvatarUsecase) *UserAdapter {
	return &UserAdapter{userUsecase: userUsecase, avatarUsecase: avatarUsecase}
}

func (a *UserAdapter) Router(r chi.Router) {
	r.Post("/", middleware.RequireAuthentication(a.Register()))
	r.Post("/validate", middleware.RequireAuthentication(a.Validate()))

	r.Get("/{user_id}", middleware.RequireAuthentication(a.GetByID()))
	r.Get("/username/{username}", middleware.RequireAuthentication(a.GetByUsername()))

	r.Get("/{user_id}/avatar/upload_token", middleware.RequireAuthentication(a.GetAvatarUploadToken()))
	r.Put("/{user_id}/avatar", middleware.RequireAuthentication(a.UpdateAvatar()))
}

// @Summary Register a new user
// @Description Register a new user by providing username and password. <br>
// @Description Require `todennus/admin:create:user` scope.
// @Tags User
// @Security OAuth2Application[todennus/admin:create:user]
// @Accept json
// @Produce json
// @Param user body dto.UserRegisterRequest true "User registration data"
// @Success 201 {object} response.SwaggerSuccessResponse[dto.UserRegisterResponse] "User registered successfully"
// @Failure 400 {object} response.SwaggerBadRequestErrorResponse "Bad request"
// @Failure 403 {object} response.SwaggerForbiddenErrorResponse "Forbidden"
// @Failure 409 {object} response.SwaggerDuplicatedErrorResponse "Duplicated"
// @Router /users [post]
func (a *UserAdapter) Register() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		request, err := xhttp.ParseHTTPRequest[dto.UserRegisterRequest](r)
		if err != nil {
			response.RESTWriteAndLogInvalidRequestError(ctx, w, err)
			return
		}

		user, err := a.userUsecase.Register(ctx, request.To())
		response.NewRESTResponseHandler(ctx, dto.NewUserRegisterResponse(user), err).
			Map(http.StatusBadRequest, errordef.ErrRequestInvalid).
			Map(http.StatusForbidden, errordef.ErrForbidden).
			Map(http.StatusConflict, errordef.ErrDuplicated).
			WithDefaultCode(http.StatusCreated).
			WriteHTTPResponse(ctx, w)
	}
}

// @Summary Get user by id
// @Description Get an user information by user id.
// @Tags User
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} response.SwaggerSuccessResponse[dto.UserGetByIDResponse] "Get user successfully"
// @Failure 400 {object} response.SwaggerBadRequestErrorResponse "Bad request"
// @Failure 404 {object} response.SwaggerNotFoundErrorResponse "Not found"
// @Router /users/{user_id} [get]
func (a *UserAdapter) GetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		req, err := xhttp.ParseHTTPRequest[dto.UserGetByIDRequest](r)
		if err != nil {
			response.RESTWriteAndLogInvalidRequestError(ctx, w, err)
			return
		}

		ucReq, err := req.To(xcontext.RequestSubjectID(ctx))
		if err != nil {
			response.RESTWriteAndLogInvalidRequestError(ctx, w, err)
			return
		}

		resp, err := a.userUsecase.GetByID(ctx, ucReq)
		response.NewRESTResponseHandler(ctx, dto.NewUserGetByIDResponse(resp), err).
			Map(http.StatusBadRequest, errordef.ErrRequestInvalid).
			Map(http.StatusNotFound, errordef.ErrNotFound).
			WriteHTTPResponse(ctx, w)
	}
}

// @Summary Get user by username
// @Description Get an user information by user username. <br>
// @Tags User
// @Produce json
// @Param username path string true "Username"
// @Success 200 {object} response.SwaggerSuccessResponse[dto.UserGetByUsernameResponse] "Get user successfully"
// @Failure 400 {object} response.SwaggerBadRequestErrorResponse "Bad request"
// @Failure 404 {object} response.SwaggerNotFoundErrorResponse "Not found"
// @Router /users/username/{username} [get]
func (a *UserAdapter) GetByUsername() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		req, err := xhttp.ParseHTTPRequest[dto.UserGetByUsernameRequest](r)
		if err != nil {
			response.RESTWriteAndLogInvalidRequestError(ctx, w, err)
			return
		}

		resp, err := a.userUsecase.GetByUsername(ctx, req.To())
		response.NewRESTResponseHandler(ctx, dto.NewUserGetByUsernameResponse(resp), err).
			Map(http.StatusBadRequest, errordef.ErrRequestInvalid).
			Map(http.StatusNotFound, errordef.ErrNotFound).
			WriteHTTPResponse(ctx, w)
	}
}

// @Summary Validate user credentials
// @Description Validate the user credentials and returns the user information. <br>
// @Description Require `todennus/admin:validate:user` scope.
// @Tags User
// @Security OAuth2Application[todennus/admin:validate:user]
// @Accept json
// @Produce json
// @Param body body dto.UserValidateRequest true "Validation data"
// @Success 200 {object} response.SwaggerSuccessResponse[dto.UserValidateResponse] "Validate successfully"
// @Failure 400 {object} response.SwaggerBadRequestErrorResponse "Bad request"
// @Failure 403 {object} response.SwaggerForbiddenErrorResponse "Forbidden"
// @Router /users/validate [post]
func (a *UserAdapter) Validate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		req, err := xhttp.ParseHTTPRequest[dto.UserValidateRequest](r)
		if err != nil {
			response.RESTWriteAndLogInvalidRequestError(ctx, w, err)
			return
		}

		resp, err := a.userUsecase.ValidateCredentials(ctx, req.To())
		response.NewRESTResponseHandler(ctx, dto.NewUserValidateResponse(resp), err).
			Map(http.StatusBadRequest, errordef.ErrRequestInvalid, errordef.ErrCredentialsInvalid).
			Map(http.StatusForbidden, errordef.ErrForbidden).
			WriteHTTPResponse(ctx, w)
	}
}

// @Summary Get an avatar upload_token.
// @Description Get the upload_token used for updating the avatar image. <br>
// @Description Require `todennus/update:user.avatar` scope.
// @Tags User
// @Security OAuth2Application[todennus/update:user.avatar]
// @Produce json
// @Param user_id path string true "user_id"
// @Success 200 {object} response.SwaggerSuccessResponse[dto.AvatarGetUploadTokenResponse] "Get token successfully"
// @Failure 403 {object} response.SwaggerForbiddenErrorResponse "Forbidden"
// @Router /users/{user_id}/avatar/upload_token [get]
func (a *UserAdapter) GetAvatarUploadToken() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		req, err := xhttp.ParseHTTPRequest[dto.AvatarGetUploadTokenRequest](r)
		if err != nil {
			response.RESTWriteAndLogInvalidRequestError(ctx, w, err)
			return
		}

		ucreq, err := req.To(xcontext.RequestSubjectID(ctx))
		if err != nil {
			response.RESTWriteAndLogInvalidRequestError(ctx, w, err)
			return
		}

		resp, err := a.avatarUsecase.GetUploadToken(ctx, ucreq)
		response.NewRESTResponseHandler(ctx, dto.NewAvatarGetUploadTokenResponse(resp), err).
			Map(http.StatusForbidden, errordef.ErrForbidden).
			WriteHTTPResponse(ctx, w)
	}
}

// @Summary Update avatar.
// @Description Use a temporary_file_token to update user avatar. <br>
// @Description Require `todennus/update:user.avatar` scope.
// @Tags User
// @Security OAuth2Application[todennus/update:user.avatar]
// @Accept json
// @Produce json
// @Param user_id path string true "user_id"
// @Param body body dto.AvatarUpdateRequest true "Avatar update request"
// @Success 200 {object} response.SwaggerSuccessResponse[dto.AvatarUpdateResponse] "Update successfully"
// @Failure 400 {object} response.SwaggerBadRequestErrorResponse "Bad request"
// @Failure 403 {object} response.SwaggerForbiddenErrorResponse "Forbidden"
// @Router /users/{user_id}/avatar [put]
func (a *UserAdapter) UpdateAvatar() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		req, err := xhttp.ParseHTTPRequest[dto.AvatarUpdateRequest](r)
		if err != nil {
			response.RESTWriteAndLogInvalidRequestError(ctx, w, err)
			return
		}

		ucreq, err := req.To(xcontext.RequestSubjectID(ctx))
		if err != nil {
			response.RESTWriteAndLogInvalidRequestError(ctx, w, err)
			return
		}

		resp, err := a.avatarUsecase.Update(ctx, ucreq)
		response.NewRESTResponseHandler(ctx, dto.NewAvatarUpdateResponse(resp), err).
			Map(http.StatusBadRequest, errordef.ErrRequestInvalid).
			Map(http.StatusForbidden, errordef.ErrForbidden).
			WriteHTTPResponse(ctx, w)
	}
}
