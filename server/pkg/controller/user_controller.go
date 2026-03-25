package controller

import (
	"github.com/labstack/echo/v4"
	"ims-intro/pkg/controller/request"
	"ims-intro/pkg/controller/response"
	"ims-intro/pkg/middleware"
	"ims-intro/pkg/service"
	"net/http"
	"strconv"
	"time"
)

type UserController struct {
	userService service.IUserService
}

func NewUserController(userService service.IUserService) *UserController {
	return &UserController{userService}
}

func (controller *UserController) RegisterUserRoutes(e *echo.Echo) {
	e.POST("/login", controller.Login)
	e.POST("/signup", controller.SignUp)
	e.POST("/logout", controller.Logout)

	usersGroup := e.Group("/users")
	usersGroup.Use(middleware.AuthMiddleware)
	usersGroup.Use(middleware.Authorize([]string{"admin"}))
	usersGroup.GET("", controller.GetAllUsers)
	usersGroup.PUT("/:id/role", controller.UpdateUserRole)
}

func (controller *UserController) Login(c echo.Context) error {
	var loginRequest request.LoginRequest
	err := c.Bind(&loginRequest)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request: unable to bind the provided data to the user structure"))
	}

	loginResult, err := controller.userService.Login(loginRequest.Username, loginRequest.Password)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, response.NewErrorResponse(err.Error()))
	}

	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = loginResult.Token
	cookie.Expires = time.Now().Add(24 * time.Hour)
	cookie.Path = "/"
	cookie.HttpOnly = false
	cookie.SameSite = http.SameSiteLaxMode
	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, response.NewLoginResponse(loginResult.Role))
}

func (controller *UserController) SignUp(c echo.Context) error {
	var signUpRequest request.SignUpRequest
	err := c.Bind(&signUpRequest)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request: unable to bind the provided data to the user structure"))
	}

	err = controller.userService.SignUp(signUpRequest.ToDtoModel())
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, response.NewErrorResponse(err.Error()))
	}

	return c.NoContent(http.StatusCreated)
}

func (controller *UserController) Logout(c echo.Context) error {
	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = ""
	cookie.Expires = time.Now().Add(-24 * time.Hour)
	cookie.Path = "/"
	cookie.HttpOnly = false
	cookie.SameSite = http.SameSiteLaxMode
	c.SetCookie(cookie)

	return c.NoContent(http.StatusOK)
}


func (controller *UserController) GetAllUsers(c echo.Context) error {
	users, err := controller.userService.GetAllUsers()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.NewErrorResponse(err.Error()))
	}

	result := make([]response.UserSummaryResponse, 0, len(users))
	for _, user := range users {
		result = append(result, response.UserSummaryResponse{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Role:     user.Role,
		})
	}

	return c.JSON(http.StatusOK, result)
}

func (controller *UserController) UpdateUserRole(c echo.Context) error {
	param := c.Param("id")
	if param == "" {
		return c.JSON(http.StatusBadRequest, response.NewErrorResponse("invalid request: no user id specified"))
	}

	userID, err := strconv.Atoi(param)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.NewErrorResponse("invalid request: user id must be an integer"))
	}

	updateRequest := new(request.UpdateUserRoleRequest)
	err = c.Bind(updateRequest)
	if err != nil || updateRequest == nil {
		return c.JSON(http.StatusBadRequest, response.NewErrorResponse("invalid request: unable to bind role payload"))
	}

	err = controller.userService.UpdateUserRole(int64(userID), updateRequest.Role)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, response.NewErrorResponse(err.Error()))
	}

	return c.NoContent(http.StatusOK)
}
