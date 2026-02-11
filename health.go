package seata

import (
	"fmt"

	"github.com/go-lynx/lynx/log"
)

// CheckHealth performs a health check of the Seata client.
// When disabled, returns nil. When enabled, verifies the plugin is in a valid state.
func (t *TxSeataClient) CheckHealth() error {
	if !t.IsEnabled() {
		return nil
	}

	if t.conf == nil {
		return fmt.Errorf("seata configuration is nil")
	}

	if t.conf.GetConfigFilePath() == "" {
		return fmt.Errorf("seata config file path is empty")
	}

	// Record health check in metrics if available
	if t.getMetrics() != nil {
		t.getMetrics().RecordHealthCheck("success")
	}

	log.Debugf("Seata client health check passed")
	return nil
}
