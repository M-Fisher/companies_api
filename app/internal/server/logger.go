package server

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func (s *Server) configureLogger() {
	level := zap.NewAtomicLevel()
	logconf := zap.Config{
		Level:       level,
		Encoding:    "json",
		Development: false,
		OutputPaths: []string{"stdout"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:     "msg",
			LevelKey:       "severity",
			TimeKey:        "timestamp",
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.FullCallerEncoder,
		},
	}

	if s.Config.DevMode {
		level.SetLevel(zapcore.DebugLevel)
		logconf.Encoding = "console"
		logconf.Development = true
		logconf.EncoderConfig.EncodeLevel = zapcore.LowercaseColorLevelEncoder
	}

	logconf.Level = level

	log, err := logconf.Build()
	if err != nil {
		panic(err)
	}

	s.Log = zap.New(customCore{log.Core()})

}

type customCore struct {
	zapcore.Core
}

func (c customCore) Write(e zapcore.Entry, f []zapcore.Field) error {
	var levelName string

	switch e.Level {
	case zapcore.DebugLevel, zapcore.InfoLevel:
		levelName = "low"
	case zapcore.WarnLevel:
		levelName = "medium"
	case zapcore.ErrorLevel, zapcore.DPanicLevel,
		zapcore.PanicLevel, zapcore.FatalLevel:
		levelName = "high"
	default:
		levelName = "low"
	}

	f = append(f, zap.String("severity_level", levelName))

	return c.Core.Write(e, f)
}

func (c customCore) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Core.Enabled(ent.Level) {
		return ce.AddCore(ent, c)
	}

	return ce
}

func (c customCore) With(fields []zapcore.Field) zapcore.Core {
	return customCore{c.Core.With(fields)}
}
