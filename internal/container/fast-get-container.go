package container

import "go.uber.org/zap"

func GetLog() *zap.Logger {
	return DiContainer.Get("zap-logger").(*zap.Logger)
}
