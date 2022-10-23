package internal

import (
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger *zap.Logger
	db     *Database
)

func NewLogger() (*zap.Logger, error) {
	if logger != nil {
		return logger, nil
	}

	db = NewDatabase("logs")
	err := db.Connect()
	if err != nil {
		return nil, err
	}

	err = db.InitDatabase()
	if err != nil {
		return nil, err
	}

	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	encoder := zapcore.NewJSONEncoder(config)
	mw := mongoWriter{database: db}
	writer := zapcore.AddSync(mw)
	defaultLogLevel := zapcore.DebugLevel
	core := zapcore.NewTee(
		zapcore.NewCore(encoder, writer, defaultLogLevel),
	)

	logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	return logger, nil
}

func CloseLogger() {
	logger = nil
	db.Disconnect()
}

type logValues struct {
	Level  string `json:"level"`
	Ts     string `json:"ts"`
	Caller string `json:"caller"`
	Msg    string `json:"msg"`
}

type mongoWriter struct {
	database *Database
}

func (mw mongoWriter) Write(p []byte) (n int, err error) {
	originalLen := len(p)
	if len(p) > 0 && p[len(p)-1] == '\n' {
		p = p[:len(p)-1]
	}

	var value logValues
	if err := json.Unmarshal(p, &value); err != nil {
		return 0, err
	}

	if _, err := mw.database.Insert(bson.M{
		"level":   value.Level,
		"created": value.Ts,
		"caller":  value.Caller,
		"message": value.Msg,
	}); err != nil {
		return 0, err
	}

	return originalLen, nil
}
