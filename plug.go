// Package seata provides a Seata distributed-transaction plugin for the go-lynx
// framework. It wraps seata-go and exposes global transaction management via
// TxSeataClient, supporting AT/TCC/SAGA/XA modes as configured in the Seata
// configuration file. The plugin integrates with the lynx runtime contract,
// Prometheus metrics, and context-aware lifecycle hooks.
package seata

import (
	"github.com/go-lynx/lynx/pkg/factory"
	"github.com/go-lynx/lynx/plugins"
)

// init registers the Seata plugin with the global factory on import.
func init() {
	factory.GlobalTypedFactory().RegisterPlugin(pluginName, confPrefix, func() plugins.Plugin {
		return NewTxSeataClient()
	})
}
