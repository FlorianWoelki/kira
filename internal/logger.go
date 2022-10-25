package internal

import (
	"encoding/json"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	Logger         *zap.Logger
	db             *Database
	collectionName string
)

func NewLogger() (*zap.Logger, error) {
	if Logger != nil {
		return Logger, nil
	}

	collectionName = "logs_" + time.Now().Format(time.RFC3339)
	db = NewDatabase(collectionName)
	err := db.Connect()
	if err != nil {
		return nil, err
	}

	err = db.CreateCollection()
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

	Logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	return Logger, nil
}

func RotateDatabase() error {
	collectionName = "logs_" + time.Now().Format(time.RFC3339)
	return db.CreateCollection()
}

func CloseLogger() {
	Logger = nil
	db.Disconnect()
}

type logValues struct {
	Level        string `json:"level"`
	Ts           string `json:"ts"`
	Caller       string `json:"caller"`
	Msg          string `json:"msg"`
	RequestBody  string `json:"requestBody"`
	ResponseBody string `json:"responseBody"`
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
		"level":        value.Level,
		"created":      value.Ts,
		"caller":       value.Caller,
		"message":      value.Msg,
		"requestBody":  value.RequestBody,
		"responseBody": value.ResponseBody,
	}); err != nil {
		return 0, err
	}

	return originalLen, nil
}
