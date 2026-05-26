package logger

import "go.uber.org/zap"

// New создаёт production-логгер с выводом в stdout и в файл app.log.
func New() (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{"stdout", "app.log"}
	cfg.ErrorOutputPaths = []string{"stderr"}
	return cfg.Build()
}
