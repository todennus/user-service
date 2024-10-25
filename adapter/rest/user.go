package rest

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/todennus/shared/errordef"
	"github.com/todennus/shared/middleware"
	"github.com/todennus/shared/response"
	"github.com/todennus/user-service/adapter/abstraction"
	"github.com/todennus/user-service/adapter/rest/dto"
	"github.com/todennus/x/xcontext"
	"github.com/todennus/x/xhttp"
)

type UserRESTAdapter struct {
	userUsecase abstraction.UserUsecase
}

func NewUserAdapter(userUsecase abstraction.UserUsecase) *UserRESTAdapter {
	return &UserRESTAdapter{userUsecase: userUsecase}
}

func (a *UserRESTAdapter) Router(r chi.Router) {
	r.Post("/", a.Register())
	r.Post("/validate", a.Validate())

	r.Get("/{user_id}", middleware.RequireAuthentication(a.GetByID()))
	r.Get("/username/{username}", middleware.RequireAuthentication(a.GetByUsername()))
}

// @Summary Register a new user
// @Description Register a new user by providing username and password
// @Tags User
// @Accept json
// @Produce json
// @Param user body dto.UserRegisterRequest true "User registration data"
// @Success 201 {object} response.SwaggerSuccessResponse[dto.UserRegisterResponse] "User registered successfully"
// @Failure 400 {object} response.SwaggerBadRequestErrorResponse "Bad request"
// @Failure 409 {object} response.SwaggerDuplicatedErrorResponse "Duplicated"
// @Router /users [post]
func (a *UserRESTAdapter) Register() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		request, err := xhttp.ParseHTTPRequest[dto.UserRegisterRequest](r)
		if err != nil {
			response.RESTWriteAndLogInvalidRequestError(ctx, w, err)
			return
		}

		user, err := a.userUsecase.Register(ctx, request.To())
		response.NewRESTResponseHandler(ctx, dto.NewUserRegisterResponse(user), err).
			WithDefaultCode(http.StatusCreated).
			Map(http.StatusConflict, errordef.ErrDuplicated).
			Map(http.StatusBadRequest, errordef.ErrRequestInvalid).
			WriteHTTPResponse(ctx, w)
	}
}

// @Summary Get user by id
// @Description Get an user information by user id. <br>
// @Tags User
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} response.SwaggerSuccessResponse[dto.UserGetByIDResponse] "Get user successfully"
// @Failure 400 {object} response.SwaggerBadRequestErrorResponse "Bad request"
// @Failure 404 {object} response.SwaggerNotFoundErrorResponse "Not found"
// @Router /users/{user_id} [get]
func (a *UserRESTAdapter) GetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		req, err := xhttp.ParseHTTPRequest[dto.UserGetByIDRequest](r)
		if err != nil {
			response.RESTWriteAndLogInvalidRequestError(ctx, w, err)
			return
		}

		ucReq, err := req.To(xcontext.RequestUserID(ctx))
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
func (a *UserRESTAdapter) GetByUsername() http.HandlerFunc {
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
// @Description Validate the user credentials and returns the user information.
// @Tags User
// @Accept json
// @Produce json
// @Param body body dto.UserValidateRequest true "Validation data"
// @Success 200 {object} response.SwaggerSuccessResponse[dto.UserValidateResponse] "Validate successfully"
// @Failure 400 {object} response.SwaggerBadRequestErrorResponse "Bad request"
// @Failure 401 {object} response.SwaggerInvalidCredentialsErrorResponse "Invalid credentials"
// @Router /users/validate [post]
func (a *UserRESTAdapter) Validate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		req, err := xhttp.ParseHTTPRequest[dto.UserValidateRequest](r)
		if err != nil {
			response.RESTWriteAndLogInvalidRequestError(ctx, w, err)
			return
		}

		resp, err := a.userUsecase.ValidateCredentials(ctx, req.To())
		response.NewRESTResponseHandler(ctx, dto.NewUserValidateResponse(resp), err).
			Map(http.StatusBadRequest, errordef.ErrRequestInvalid).
			Map(http.StatusUnauthorized, errordef.ErrCredentialsInvalid).
			WriteHTTPResponse(ctx, w)
	}
}
