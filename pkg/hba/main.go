package hba

import (
	"go.uber.org/zap"
)

var (
	log *zap.SugaredLogger
)

// TODO Is this copying a structure?
func InitLogger(logger *zap.SugaredLogger) {
	log = logger
}
