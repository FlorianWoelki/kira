package pkg

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/florianwoelki/kira/internal"
	"github.com/robfig/cron/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	ROTATE_MINUTE = "*/1 * * * *"
	ROTATE_HOUR   = "0 * * * *"
	ROTATE_DAY    = "0 0 * * *"
	ROTATE_WEEK   = "0 0 * * 0"
	ROTATE_MONTH  = "0 0 1 * *"
)

var (
	Logger  *zap.Logger
	db      *internal.Database
	cronJob *cron.Cron
	// An internal logger that does only log functionalities like rotating.
	logger *log.Logger = log.New(os.Stdout, "language: ", log.LstdFlags|log.Lshortfile)
)

// NewLogger creates a new logger with a specific rotation. It also connects to the
// database and creates a new `zap` logger.
func NewLogger(rotation string) (*zap.Logger, error) {
	if Logger != nil {
		return Logger, nil
	}

	// Create new database with new collection inside.
	db = internal.NewDatabase()
	err := db.Connect()
	if err != nil {
		return nil, err
	}

	collectionName := fmt.Sprintf("logs_%d", time.Now().UnixMilli())
	err = db.CreateCollection(collectionName)
	if err != nil {
		return nil, err
	}

	// Create cron job for rotating collection.
	cronJob = cron.New()
	_, err = cronJob.AddFunc(rotation, func() {
		rotationName, err := rotateDatabase()
		if err != nil {
			logger.Printf("Something went wrong while executing cron job: %v+\n", err)
		}

		logger.Printf("Rotating database to %s...", rotationName)
	})
	if err != nil {
		return nil, err
	}
	cronJob.Start()

	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	mw := mongoWriter{database: db}

	// Create all logging encoder.
	jsonEncoder := zapcore.NewJSONEncoder(config)
	consoleEncoder := zapcore.NewConsoleEncoder(config)

	// Creates the core zap logger with console and mongodb writing.
	defaultLogLevel := zapcore.DebugLevel
	core := zapcore.NewTee(
		zapcore.NewCore(jsonEncoder, zapcore.AddSync(mw), defaultLogLevel),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), defaultLogLevel),
	)
	Logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	return Logger, nil
}

// rotateDatabase rotates the database by defining the new rotated collection name and
// then creating it in the database.
func rotateDatabase() (string, error) {
	collectionName := fmt.Sprintf("logs_%d", time.Now().UnixNano())
	return collectionName, db.CreateCollection(collectionName)
}

// CloseLogger sets the initialized logger to `nil`, stops the cronjob and disconnects
// from the database.
func CloseLogger() {
	Logger = nil
	cronJob.Stop()
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

// mongoWriter implements the `io.Writer` interface and it is being used for logging to
// mongodb.
type mongoWriter struct {
	database *internal.Database
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
