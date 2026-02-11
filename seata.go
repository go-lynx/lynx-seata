package seata

import (
	"context"
	"fmt"
	"time"

	"github.com/go-lynx/lynx"
	"github.com/go-lynx/lynx-seata/conf"
	"github.com/go-lynx/lynx/log"
	"github.com/go-lynx/lynx/plugins"
	"github.com/seata/seata-go/pkg/client"
	"github.com/seata/seata-go/pkg/tm"
)

// Plugin metadata
const (
	// pluginName is the unique identifier for the Seata plugin, used to identify the plugin in the plugin system.
	pluginName = "seata.server"

	// pluginVersion indicates the current version of the Seata plugin.
	pluginVersion = "v2.0.0"

	// pluginDescription briefly describes the functionality of the Seata plugin.
	pluginDescription = "Seata distributed transaction plugin for Lynx framework"

	// confPrefix is the configuration prefix used when loading Seata configuration.
	confPrefix = "lynx.seata"

	// defaultTransactionTimeout is the default global transaction timeout.
	defaultTransactionTimeout = 60 * time.Second
)

type TxSeataClient struct {
	// Embed base plugin, inherit common properties and methods of the plugin
	*plugins.BasePlugin
	// Seata configuration information
	conf *conf.Seata
}

// NewTxSeataClient creates a new Seata plugin instance.
func NewTxSeataClient() *TxSeataClient {
	return &TxSeataClient{
		BasePlugin: plugins.NewBasePlugin(
			plugins.GeneratePluginID("", pluginName, pluginVersion),
			pluginName,
			pluginDescription,
			pluginVersion,
			confPrefix,
			90,
		),
	}
}

// InitializeResources loads and initializes the Seata plugin configuration.
func (t *TxSeataClient) InitializeResources(rt plugins.Runtime) error {
	t.conf = &conf.Seata{}

	err := rt.GetConfig().Value(confPrefix).Scan(t.conf)
	if err != nil {
		return err
	}

	// Set default configuration
	if t.conf.ConfigFilePath == "" {
		t.conf.ConfigFilePath = "./conf/seata.yml"
	}

	return nil
}

// StartupTasks initializes the Seata client when enabled.
func (t *TxSeataClient) StartupTasks() error {
	log.Infof("Initializing seata")
	if t.conf.GetEnabled() {
		client.InitPath(t.conf.GetConfigFilePath())
	} else {
		log.Infof("Seata client is disabled")
		return nil
	}
	log.Infof("Seata successfully initialized")
	return nil
}

// CleanupTasks performs cleanup during plugin shutdown.
// Seata-go does not expose a public shutdown API; connections will be released when the process exits.
func (t *TxSeataClient) CleanupTasks() error {
	if t.conf.GetEnabled() {
		log.Infof("Seata client shutting down")
	}
	return nil
}

// Configure updates the plugin configuration. Overrides base to apply *conf.Seata.
func (t *TxSeataClient) Configure(cfg any) error {
	if cfg == nil {
		return nil
	}
	c, ok := cfg.(*conf.Seata)
	if !ok {
		return fmt.Errorf("invalid configuration type, expected *conf.Seata")
	}
	t.conf = c
	return nil
}

// GetConfig returns the Seata configuration.
func (t *TxSeataClient) GetConfig() *conf.Seata {
	return t.conf
}

// GetConfigFilePath returns the Seata configuration file path.
func (t *TxSeataClient) GetConfigFilePath() string {
	if t.conf == nil {
		return "./conf/seata.yml"
	}
	return t.conf.GetConfigFilePath()
}

// IsEnabled returns whether the Seata plugin is enabled.
func (t *TxSeataClient) IsEnabled() bool {
	if t.conf == nil {
		return false
	}
	return t.conf.GetEnabled()
}

// WithGlobalTx executes the given business function within a global transaction.
// It wraps seata-go's tm.WithGlobalTx with sensible defaults.
func (t *TxSeataClient) WithGlobalTx(ctx context.Context, name string, timeout time.Duration, business func(context.Context) error) error {
	if !t.IsEnabled() {
		return fmt.Errorf("seata plugin is disabled")
	}
	if timeout == 0 {
		timeout = defaultTransactionTimeout
	}
	gc := &tm.GtxConfig{
		Timeout: timeout,
		Name:    name,
	}
	return tm.WithGlobalTx(ctx, gc, business)
}

// GetPlugin obtains the TxSeataClient plugin instance from the application's plugin manager.
func GetPlugin() *TxSeataClient {
	if lynx.Lynx() == nil || lynx.Lynx().GetPluginManager() == nil {
		return nil
	}
	pl := lynx.Lynx().GetPluginManager().GetPlugin(pluginName)
	if pl == nil {
		return nil
	}
	tc, ok := pl.(*TxSeataClient)
	if !ok {
		return nil
	}
	return tc
}
