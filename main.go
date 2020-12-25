package main

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"go.uber.org/zap"
)

// User ...
type User struct {
	ID   uint   `json:"id" gorm:"id"`
	Name string `json:"name" gorm:"name"`
	Age  int    `json:"age" gorm:"age"`
}

// InitMiddleware ...
func InitMiddleware(e *echo.Echo) {
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
}

// InitRouting ...
func InitRouting(e *echo.Echo, u *User) {
	e.POST("user", u.CreateUser)
	e.GET("user/:id", u.GetUser)
}

// CreateUser ...
func (u *User) CreateUser(c echo.Context) error {
	if err := c.Bind(u); err != nil {
		// zap.L()はglobal Loggerを返すため、それを用いてloggingを行う。
		zap.L().Error("failed to bind", zap.Error(err), zap.String("something_key1", "something_string_value"))
		return err
	}

	u.ID = 1

	return c.JSON(http.StatusOK, &u)
}

// GetUser ...
func (u *User) GetUser(c echo.Context) error {
	id := c.Param("id")

	return c.JSON(http.StatusOK, "User ID = "+id)
}

func main() {
	e := echo.New()

	logger, err := zap.NewDevelopment()
	if err != nil {
		return
	}

	// zap.ReplaceGlobalsにloggerをセットすることで、zap.L()が任意の場所で使用できるように。
	zap.ReplaceGlobals(logger)

	InitMiddleware(e)

	u := new(User)
	InitRouting(e, u)

	e.Start(":9000")
}
