package usecase

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"

	"github.com/go-chi/render"
	"github.com/rbo13/write-it/app"
	"github.com/rbo13/write-it/app/response"
)

type userUsecase struct {
	userService app.UserService
}

// UserResponse represents a user response
type UserResponse struct {
	StatusCode uint        `json:"status_code"`
	Message    string      `json:"message"`
	Success    bool        `json:"success"`
	Data       interface{} `json:"data"`
}

type loginResponse struct {
	UserResponse UserResponse `json:"user_response"`
	AuthToken    string       `json:"auth_token"`
}

// NewUser ...
func NewUser(userService app.UserService) app.UserHandler {
	return &userUsecase{
		userService,
	}
}

func (u *userUsecase) Create(w http.ResponseWriter, r *http.Request) {
	var user app.User

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		config := response.Configure(err.Error(), http.StatusUnprocessableEntity, nil)
		response.JSONError(w, r, config)
		return
	}

	err = u.userService.CreateUser(&user)

	if err != nil {
		config := response.Configure(err.Error(), http.StatusBadRequest, nil)
		response.JSONError(w, r, config)
		return
	}

	config := response.Configure("User successfully registered", http.StatusOK, user)
	response.JSONOK(w, r, config)
	return
}

func (u *userUsecase) Login(w http.ResponseWriter, r *http.Request) {
	var user app.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		loginResp := loginResponse{
			UserResponse: errorResponse(http.StatusUnprocessableEntity, err.Error()),
			AuthToken:    "",
		}

		config := response.Configure(err.Error(), http.StatusUnprocessableEntity, &loginResp)
		response.JSONError(w, r, config)
		return
	}

	userResp, err := u.userService.Login(user.EmailAddress, user.Password)

	if err != nil {
		loginResp := loginResponse{
			UserResponse: errorResponse(http.StatusNotFound, err.Error()),
			AuthToken:    "",
		}

		config := response.Configure(err.Error(), http.StatusUnprocessableEntity, &loginResp)
		response.JSONError(w, r, config)
		return
	}

	authToken, err := u.userService.GenerateAuthToken(userResp)

	if err != nil {
		loginResp := loginResponse{
			UserResponse: errorResponse(http.StatusBadRequest, err.Error()),
			AuthToken:    "",
		}

		config := response.Configure(err.Error(), http.StatusUnprocessableEntity, &loginResp)
		response.JSONError(w, r, config)
		return
	}

	loginResp := map[string]interface{}{
		"user":       userResp,
		"auth_token": authToken,
	}

	config := response.Configure("Logged in sucessfully", http.StatusOK, loginResp)
	response.JSONOK(w, r, config)
}

func (u *userUsecase) Get(w http.ResponseWriter, r *http.Request) {
	users, err := u.userService.Users()

	if err != nil || users == nil {
		config := response.Configure(err.Error(), http.StatusNotFound, users)
		response.JSONError(w, r, config)
		return
	}

	config := response.Configure("Users successfully retrieved", http.StatusOK, users)
	response.JSONOK(w, r, config)
}

func (u *userUsecase) GetByID(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	if err != nil || userID <= 0 {
		config := response.Configure(err.Error(), http.StatusUnprocessableEntity, nil)
		response.JSONError(w, r, config)
		return
	}

	user, err := u.userService.User(userID)

	if err != nil {
		config := response.Configure(err.Error(), http.StatusNotFound, nil)
		response.JSONError(w, r, config)
		return
	}

	config := response.Configure("User successfully retrieved", http.StatusOK, user)
	response.JSONOK(w, r, config)
}

func (u *userUsecase) Update(w http.ResponseWriter, r *http.Request) {
	var user app.User
	userID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	if err != nil {
		config := response.Configure(err.Error(), http.StatusUnprocessableEntity, nil)
		response.JSONError(w, r, config)
		return
	}

	user.ID = userID
	user.UpdatedAt = time.Now().Unix()

	err = json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		config := response.Configure(err.Error(), http.StatusUnprocessableEntity, nil)
		response.JSONError(w, r, config)
		return
	}

	err = u.userService.UpdateUser(&user)

	if err != nil {
		config := response.Configure(err.Error(), http.StatusUnprocessableEntity, nil)
		response.JSONError(w, r, config)
		return
	}

	config := response.Configure("User successfully updated", http.StatusOK, user)
	response.JSONOK(w, r, config)
}

func (u *userUsecase) Delete(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	if err != nil {
		deleteResponse := UserResponse{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    err.Error(),
			Success:    false,
			Data:       nil,
		}
		render.JSON(w, r, &deleteResponse)
		return
	}

	err = u.userService.DeleteUser(userID)

	if err != nil {
		deleteResponse := UserResponse{
			StatusCode: http.StatusNotFound,
			Message:    err.Error(),
			Success:    false,
			Data:       nil,
		}
		render.JSON(w, r, &deleteResponse)
		return
	}

	deleteResponse := UserResponse{
		StatusCode: http.StatusNoContent,
		Message:    "User successfully deleted",
		Success:    true,
		Data:       nil,
	}
	render.JSON(w, r, &deleteResponse)
}

func errorResponse(statusCode uint, message string) (errResponse UserResponse) {
	errResponse = UserResponse{
		StatusCode: statusCode,
		Message:    message,
		Success:    false,
		Data:       nil,
	}

	return errResponse
}

func okResponse(statusCode uint, data interface{}, message string) (okResponse UserResponse) {
	okResponse = UserResponse{
		StatusCode: statusCode,
		Message:    message,
		Success:    true,
		Data:       data,
	}

	return okResponse
}
