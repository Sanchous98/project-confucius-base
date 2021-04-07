package stdlib

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"os/signal"
	"syscall"
	"testing"
)

func TestLogger(t *testing.T) {
	logger := new(Log)
	logger.Constructor()
	assert.NotPanics(t, func() { logger.Debug(fmt.Errorf("debug")) })
	assert.NotPanics(t, func() { logger.Info(fmt.Errorf("info")) })
	assert.NotPanics(t, func() { logger.Notice(fmt.Errorf("notice")) })
	assert.NotPanics(t, func() { logger.Warning(fmt.Errorf("warning")) })
	assert.NotPanics(t, func() { logger.Error(fmt.Errorf("error")) })
	assert.NotPanics(t, func() { logger.Critical(fmt.Errorf("critical")) })
	assert.Panics(t, func() { logger.Alert(fmt.Errorf("alert")) })
	assert.NotPanics(t, func() {
		stopChannel := make(chan os.Signal, 2)
		signal.Notify(stopChannel, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
		logger.Emergency(fmt.Errorf("emergency"))

		select {
		case <-stopChannel:
			return
		}
	})
}
