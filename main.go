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
		// zap.L()およびzap.S()はglobal Loggerを返すため、それを用いてlogging
		zap.S().Errorw("failed to bind", zap.Error(err), zap.String("something_key1", "something_string_value"))
		return err
	}

	// DB使用していないため、仮でIDに1を代入
	u.ID = 1

	return c.JSON(http.StatusOK, &u)
}

// --------
// conf↓
// --------

// InitDevLogger 開発環境用
func InitDevLogger() (*zap.Logger, error) {
	level := zap.NewAtomicLevel()
	level.SetLevel(zapcore.DebugLevel)

	myConfig := zap.Config{
		Level:             level,
		Development:       true,
		DisableStacktrace: true, // Stacktraceを表示すべきかはもう少し検討
		Encoding:          "console",
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

// InitPrdLogger 本番環境用
func InitPrdLogger() (*zap.Logger, error) {
	level := zap.NewAtomicLevel()

	myConfig := zap.Config{
		Level:             level,
		DisableStacktrace: true, // Stacktraceを表示すべきかはもう少し検討
		Encoding:          "json",
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:     "Msg",
			LevelKey:       "Level",
			TimeKey:        "Time",
			NameKey:        "Name",
			CallerKey:      "Caller",
			StacktraceKey:  "St",
			EncodeLevel:    zapcore.CapitalLevelEncoder, // JSON形式の場合、Logging Levelの色が表示できないため開発環境用から変更
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.FullCallerEncoder,
			EncodeName:     zapcore.FullNameEncoder,
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

	// loggerの初期化 ※InitDevLogger、InitPrdLoggerの切り替えは未実装
	logger, err := InitPrdLogger()
	if err != nil {
		return
	}
	defer logger.Sync()

	undo := zap.ReplaceGlobals(logger)
	defer undo()

	u := new(User)
	InitRouting(e, u)

	if err = e.Start(":9000"); err != nil {
		zap.S().Fatalw("HTTP Server 起動エラー", zap.Error(err))
	}
}
