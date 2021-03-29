package stdlib

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLogger(t *testing.T) {
	logger := new(Log)
	logger.Constructor()
	assert.NotPanics(t, func() { logger.Debug(fmt.Errorf("error")) })
	assert.NotPanics(t, func() { logger.Info(fmt.Errorf("error")) })
	assert.NotPanics(t, func() { logger.Notice(fmt.Errorf("error")) })
	assert.NotPanics(t, func() { logger.Warning(fmt.Errorf("error")) })
	assert.NotPanics(t, func() { logger.Error(fmt.Errorf("error")) })
	assert.NotPanics(t, func() { logger.Critical(fmt.Errorf("error")) })
	assert.NotPanics(t, func() { logger.Alert(fmt.Errorf("error")) })
	assert.NotPanics(t, func() { logger.Emergency(fmt.Errorf("error")) })
}
