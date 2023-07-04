package electrum

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/checksum0/go-electrum/electrum"
	"github.com/ipfs/go-log"
	"golang.org/x/exp/slices"

	"github.com/keep-network/keep-common/pkg/wrappers"
	"github.com/keep-network/keep-core/pkg/bitcoin"
)

// TODO: Some problems with Electrum re-connect were detected while developing
//       integration tests: https://github.com/keep-network/keep-core/issues/3586.
//       Make sure the problem is resolved soon.

var (
	supportedProtocolVersions = []string{"1.4"}
	logger                    = log.Logger("keep-electrum")
)

// Connection is a handle for interactions with Electrum server.
type Connection struct {
	parentCtx   context.Context
	client      *electrum.Client
	clientMutex *sync.RWMutex
	config      Config
}

// Connect initializes handle with provided Config.
func Connect(parentCtx context.Context, config Config) (bitcoin.Chain, error) {
	if config.ConnectTimeout == 0 {
		config.ConnectTimeout = DefaultConnectTimeout
	}
	if config.ConnectRetryTimeout == 0 {
		config.ConnectRetryTimeout = DefaultConnectRetryTimeout
	}
	if config.RequestTimeout == 0 {
		config.RequestTimeout = DefaultRequestTimeout
	}
	if config.RequestRetryTimeout == 0 {
		config.RequestRetryTimeout = DefaultRequestRetryTimeout
	}
	if config.KeepAliveInterval == 0 {
		config.KeepAliveInterval = DefaultKeepAliveInterval
	}

	c := &Connection{
		parentCtx:   parentCtx,
		config:      config,
		clientMutex: &sync.RWMutex{},
	}

	if err := c.electrumConnect(); err != nil {
		return nil, fmt.Errorf("failed to initialize electrum client: [%w]", err)
	}

	if err := c.verifyServer(); err != nil {
		return nil, fmt.Errorf("failed to verify electrum server: [%w]", err)
	}

	// Keep the connection alive and check the connection health.
	go c.keepAlive()

	// TODO: Add reconnects on lost connection.

	return c, nil
}

func (c *Connection) electrumConnect() error {
	var client *electrum.Client
	var err error

	logger.Debug("establishing connection to electrum server...")
	client, err = connectWithRetry(
		c,
		func(ctx context.Context) (*electrum.Client, error) {
			return electrum.NewClient(ctx, c.config.URL, nil)
		},
	)

	if err == nil {
		c.client = client
	}

	return err
}

func (c *Connection) verifyServer() error {
	type Server struct {
		version  string
		protocol string
	}

	server, err := requestWithRetry(
		c,
		func(ctx context.Context, client *electrum.Client) (*Server, error) {
			serverVersion, protocolVersion, err := client.ServerVersion(ctx)
			if err != nil {
				return nil, err
			}
			return &Server{serverVersion, protocolVersion}, nil
		},
	)
	if err != nil {
		return fmt.Errorf("failed to get server version: [%w]", err)
	}

	logger.Infof(
		"connected to electrum server [version: [%s], protocol: [%s]]",
		server.version,
		server.protocol,
	)

	// Log a warning if connected to a server running an unsupported protocol version.
	if !slices.Contains(supportedProtocolVersions, server.protocol) {
		logger.Warnf(
			"electrum server [%s] runs an unsupported protocol version: [%s]; expected one of: [%s]",
			c.config.URL,
			server.protocol,
			strings.Join(supportedProtocolVersions, ","),
		)
	}

	return nil
}

func (c *Connection) keepAlive() {
	ticker := time.NewTicker(c.config.KeepAliveInterval)

	for {
		select {
		case <-ticker.C:
			_, err := requestWithRetry(
				c,
				func(ctx context.Context, client *electrum.Client) (interface{}, error) {
					return nil, client.Ping(ctx)
				},
			)
			if err != nil {
				logger.Errorf(
					"failed to ping the electrum server; "+
						"please verify health of the electrum server: [%v]",
					err,
				)
			} else {
				// Adjust ticker starting at the time of the latest successful ping.
				ticker = time.NewTicker(c.config.KeepAliveInterval)
			}
		case <-c.parentCtx.Done():
			ticker.Stop()
			c.client.Shutdown()
			return
		}
	}
}

func connectWithRetry(
	c *Connection,
	newClientFn func(ctx context.Context) (*electrum.Client, error),
) (*electrum.Client, error) {
	var result *electrum.Client
	err := wrappers.DoWithDefaultRetry(
		c.parentCtx,
		c.config.ConnectRetryTimeout,
		func(ctx context.Context) error {
			connectCtx, connectCancel := context.WithTimeout(
				ctx,
				c.config.ConnectTimeout,
			)
			defer connectCancel()

			client, err := newClientFn(connectCtx)
			if err == nil {
				result = client
			}

			return err
		},
	)

	return result, err
}

func requestWithRetry[K interface{}](
	c *Connection,
	requestFn func(ctx context.Context, client *electrum.Client) (K, error),
) (K, error) {
	var result K

	err := wrappers.DoWithDefaultRetry(
		c.parentCtx,
		c.config.RequestRetryTimeout,
		func(ctx context.Context) error {
			if err := c.reconnectIfShutdown(); err != nil {
				return err
			}

			requestCtx, requestCancel := context.WithTimeout(ctx, c.config.RequestTimeout)
			defer requestCancel()

			c.clientMutex.RLock()
			r, err := requestFn(requestCtx, c.client)
			c.clientMutex.RUnlock()

			if err != nil {
				return fmt.Errorf("request failed: [%w]", err)
			}

			result = r
			return nil
		})

	return result, err
}

func (c *Connection) reconnectIfShutdown() error {
	c.clientMutex.Lock()
	defer c.clientMutex.Unlock()

	isClientShutdown := c.client.IsShutdown()
	if isClientShutdown {
		logger.Warn("connection to electrum server is down; reconnecting...")
		err := c.electrumConnect()
		if err != nil {
			return fmt.Errorf("failed to reconnect to electrum server: [%w]", err)
		}
		logger.Info("reconnected to electrum server")
	}

	return nil
}
