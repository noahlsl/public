package logsx

import (
	"fmt"
	"go.uber.org/zap/zapcore"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"go.uber.org/zap"
)

const (
	callerSkipOffset  = 3
	defaultServerName = "service"
)

type ZapWriter struct {
	logger *zap.Logger
}

func NewZapWriter(c logx.LogConf, fields ...zap.Field) (logx.Writer, *zap.Logger) {
	if c.ServiceName != "" {
		fields = append(fields, zap.String(defaultServerName, c.ServiceName))
	}
	writer, logger, err := newZapWriter(c, zap.Fields(fields...))
	if err != nil {
		panic(err)
	}
	return writer, logger
}

func newZapWriter(c logx.LogConf, opts ...zap.Option) (logx.Writer, *zap.Logger, error) {
	opts = append(opts, zap.AddCallerSkip(callerSkipOffset))
	cfg := zap.NewProductionConfig()
	cfg.DisableCaller = true
	cfg.DisableStacktrace = true
	cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.DateTime)
	cfg.Level = zap.NewAtomicLevelAt(getZapLevel(c.Level))
	logger, err := cfg.Build(opts...)
	if err != nil {
		return nil, nil, err
	}

	return &ZapWriter{
		logger: logger,
	}, logger, nil
}

func (w *ZapWriter) Alert(v interface{}) {
	w.logger.Error(fmt.Sprint(v))
}

func (w *ZapWriter) Close() error {
	return w.logger.Sync()
}

func (w *ZapWriter) Debug(v interface{}, fields ...logx.LogField) {
	w.logger.Debug(fmt.Sprint(v), toZapFields(fields...)...)
}

func (w *ZapWriter) Error(v interface{}, fields ...logx.LogField) {
	w.logger.Error(fmt.Sprint(v), toZapFields(fields...)...)
}

func (w *ZapWriter) Info(v interface{}, fields ...logx.LogField) {
	w.logger.Info(fmt.Sprint(v), toZapFields(fields...)...)
}

func (w *ZapWriter) Severe(v interface{}) {
	w.logger.Fatal(fmt.Sprint(v))
}

func (w *ZapWriter) Slow(v interface{}, fields ...logx.LogField) {
	w.logger.Warn(fmt.Sprint(v), toZapFields(fields...)...)
}

func (w *ZapWriter) Stack(v interface{}) {
	w.logger.Error(fmt.Sprint(v), zap.Stack("stack"))
}

func (w *ZapWriter) Stat(v interface{}, fields ...logx.LogField) {
	w.logger.Info(fmt.Sprint(v), toZapFields(fields...)...)
}

func toZapFields(fields ...logx.LogField) []zap.Field {
	zapFields := make([]zap.Field, 0, len(fields))
	for _, f := range fields {
		zapFields = append(zapFields, zap.Any(f.Key, f.Value))
	}
	return zapFields
}

func getZapLevel(in string) zapcore.Level {
	switch strings.ToLower(in) {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "panic":
		return zapcore.PanicLevel
	case "fatal":
		return zapcore.FatalLevel
	}
	return zapcore.ErrorLevel
}
