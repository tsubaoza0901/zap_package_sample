package main

import (
	"net/http"

	"github.com/labstack/echo"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// --------
// model↓
// --------

// User ...
type User struct {
	ID   uint   `json:"id" gorm:"id"`
	Name string `json:"name" gorm:"name"`
	Age  int    `json:"age" gorm:"age"`
}

// --------
// router↓
// --------

// InitRouting ...
func InitRouting(e *echo.Echo, u *User) {
	e.POST("user", u.CreateUser)
}

// --------
// handler↓
// --------

// CreateUser ...
func (u *User) CreateUser(c echo.Context) error {
	if err := c.Bind(u); err != nil {
		// zap.L()はglobal Loggerを返すため、それを用いてlogging
		zap.L().Error("failed to bind", zap.Error(err), zap.String("something_key1", "something_string_value"))
		return err
	}

	// DB使用していないため、仮でIDに1を代入
	u.ID = 1

	return c.JSON(http.StatusOK, &u)
}

// --------
// conf↓
// --------

// InitLogger ...
func InitLogger() (*zap.Logger, error) {
	level := zap.NewAtomicLevel()
	level.SetLevel(zapcore.DebugLevel)

	myConfig := zap.Config{
		Level:    level,
		Encoding: "console",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "Time",
			LevelKey:       "Level",
			NameKey:        "Name",
			CallerKey:      "Caller",
			MessageKey:     "Msg",
			StacktraceKey:  "St",
			EncodeLevel:    zapcore.CapitalColorLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.FullCallerEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}
	return myConfig.Build()
}

// --------
// main.go↓
// --------

// main
func main() {
	e := echo.New()

	// loggerの初期化
	logger, err := InitLogger()
	if err != nil {
		return
	}

	// zap.ReplaceGlobalsにloggerをセットすることで、zap.L()が任意の場所で使用できるように
	zap.ReplaceGlobals(logger)

	u := new(User)
	InitRouting(e, u)

	e.Start(":9000")
}
