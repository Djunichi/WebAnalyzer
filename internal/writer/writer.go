package writer

import "go.uber.org/zap"

type ZapWriter struct {
	logger *zap.SugaredLogger
}

func NewZapWriter(logger *zap.SugaredLogger) *ZapWriter {
	return &ZapWriter{logger: logger}
}

func (z *ZapWriter) Write(p []byte) (n int, err error) {
	z.logger.Info(string(p))
	return len(p), nil
}
