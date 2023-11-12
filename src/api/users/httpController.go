package users

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UserController ...
var (
	UserController UserControllerInterface = &userController{}
)

type UserControllerInterface interface {
	Create(c *gin.Context)
	Login(c *gin.Context)
	Logout(c *gin.Context)
	GetOne(c *gin.Context)
	GetAll(c *gin.Context)
	Update(c *gin.Context)
	PasswordUpdate(c *gin.Context)
	Delete(c *gin.Context)
}

type userController struct {
	service UserServiceInterface
}

func NewUserController(ser UserServiceInterface) UserControllerInterface {
	return &userController{
		ser,
	}
}

// ///////controllers/////////////////
// Create godoc
// @Summary Create a user
// @Description Create a new user item
// @Tags users
// @Accept json
// @Produce json
// @Success 201 {object} User
// @Failure 400 {object} support.HttpError
// @Router /register [post]
func (controller userController) Create(c *gin.Context) {

	user := &User{}
	user.Firstname = c.PostForm("firstname")
	user.Lastname = c.PostForm("lastname")
	user.Phone = c.PostForm("phone")
	user.Username = c.PostForm("username")
	user.Address = c.PostForm("address")
	user.Email = c.PostForm("email")
	user.Password = c.PostForm("password")
	u, err1 := controller.service.Create(user)
	if err1 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": err1.Message()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": u})
}

// Login godoc
// @Summary Login a user
// @Description Login a new user item
// @Tags users
// @Accept json
// @Produce json
// @Success 201 {object} User
// @Failure 400 {object} support.HttpError
// @Router /login [post]
func (controller userController) Login(c *gin.Context) {
	user := &LoginUser{}
	user.Email = c.PostForm("email")
	user.Password = c.PostForm("password")
	auth, problem := controller.service.Login(user)
	if problem != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": problem.Message()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": auth})
}

// Logout godoc
// @Summary Logout a user
// @Description Logout a new user item
// @Tags users
// @Accept json
// @Produce json
// @Success 201 {object} User
// @Failure 400 {object} support.HttpError
// @Router /api/logout [post]
func (controller userController) Logout(c *gin.Context) {
	token := string(c.Param("token"))
	_, problem := controller.service.Logout(token)
	if problem != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": problem.Message()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "succeessifully logged out"})
}

// GetOne godoc
// @Summary GetOne a user
// @Description GetOne a new user item
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} User
// @Failure 400 {object} support.HttpError
// @Router /api/users/{code} [get]
func (controller userController) GetOne(c *gin.Context) {
	code := c.Param("code")
	user, problem := controller.service.GetOne(code)
	if problem != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": problem.Message()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": user})
}

// GetAll godoc
// @Summary GetAll a user
// @Description GetAll a new user item
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} User
// @Failure 400 {object} support.HttpError
// @Router /api/users [get]
func (controller userController) GetAll(c *gin.Context) {
	search := c.Param("search")
	users, problem := controller.service.GetAll(search)
	if problem != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": problem.Message()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"users": users})
}

// Update godoc
// @Summary Update a user
// @Description Update a new user item
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} User
// @Failure 400 {object} support.HttpError
// @Router /api/users/{code} [put]
func (controller userController) Update(c *gin.Context) {
	user := &User{}
	user.Firstname = c.PostForm("firstname")
	user.Lastname = c.PostForm("lastname")
	user.Username = c.PostForm("username")
	user.Phone = c.PostForm("phone")
	user.Address = c.PostForm("address")
	user.Email = c.PostForm("email")
	// user.Business = c.FormValue("business")
	code := c.Param("code")
	_, problem := controller.service.Update(code, user)
	if problem != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": problem.Message()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "Updated successifuly"})
}

// PasswordUpdate godoc
// @Summary PasswordUpdate a user
// @Description PasswordUpdate a new user item
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} User
// @Failure 400 {object} support.HttpError
// @Router /api/users/password [put]
func (controller userController) PasswordUpdate(c *gin.Context) {
	fmt.Println("-----------------0")
	oldpassword := c.PostForm("oldpassword")
	email := c.PostForm("email")
	newpassword := c.PostForm("newpassword")
	_, _, problem := controller.service.PasswordUpdate(oldpassword, email, newpassword)
	if problem != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": problem.Message()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "Updated successifuly"})
}

// Delete godoc
// @Summary Delete a user
// @Description Delete a new user item
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} User
// @Failure 400 {object} support.HttpError
// @Router /api/users/{code} [delete]
func (controller userController) Delete(c *gin.Context) {
	id := string(c.Param("id"))
	success, failure := controller.service.Delete(id)
	if failure != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": failure.Message()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": success})

}
