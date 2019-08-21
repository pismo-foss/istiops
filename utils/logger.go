package utils

import (
	"os"

	logJSON "github.com/sirupsen/logrus"
	"github.com/snowzach/rotatefilehook"
)

type Fields map[string]interface{}

var (
	log     *logJSON.Logger
	System  string
	Env     string
	Version string
)

func init() {
	System = os.Getenv("SYSTEM")
	Version = os.Getenv("SYSTEM_VERSION")
	Env = os.Getenv("ENV")

	log = logJSON.New()
	formatter := &logJSON.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05.000",
		FieldMap: logJSON.FieldMap{
			logJSON.FieldKeyTime: "logdate",
			logJSON.FieldKeyMsg:  "text",
		},
	}
	log.SetFormatter(formatter)
	log.SetOutput(os.Stdout)
	rotateFileHook, err := rotatefilehook.NewRotateFileHook(rotatefilehook.RotateFileConfig{
		Filename:   "/var/log/app/json.log",
		MaxSize:    1000,
		MaxBackups: 7,
		MaxAge:     7,
		Level:      logJSON.InfoLevel,
		Formatter:  formatter,
	})
	if err != nil {
		log.Error(err.Error())
		return
	}
	log.AddHook(rotateFileHook)
}

func output(cid string, fields []Fields) *logJSON.Entry {
	var f Fields
	if len(fields) < 1 {
		f = Fields{}
	} else {
		f = fields[0]
	}
	f["systemVersion"] = Version
	f["cid"] = cid

	return log.WithFields(logJSON.Fields(f))
}

func Debug(message string, cid string, fields ...Fields) {
	output(cid, fields).Debug(message)
}

func Info(message string, cid string, fields ...Fields) {
	output(cid, fields).Info(message)
}

func Warn(message string, cid string, fields ...Fields) {
	output(cid, fields).Warn(message)
}

func Fatal(message string, cid string, fields ...Fields) {
	output(cid, fields).Fatal(message)
}

func Error(message string, cid string, fields ...Fields) {
	output(cid, fields).Error(message)
}
