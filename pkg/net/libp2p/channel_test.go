package libp2p

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/keep-network/keep-core/pkg/net"
)

func TestMain(m *testing.M) {
	// TODO: Invoke code that builds up state
	code := m.Run()
	// TODO: close open conns
	os.Exit(code)
}

func newTestContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 3*time.Second)
}

func TestRegisterIdentifier(t *testing.T) {
	t.Parallel()

	ctx, cancel := newTestContext()
	defer cancel()

	tests := map[string]struct {
		tIdentifier net.TransportIdentifier
		pIdentifier net.ProtocolIdentifier
		errorString string
	}{
		"invalid transport identifier":        {},
		"transport identifier already exists": {},
		"protocol identifier already exists":  {},
		"new, valid transport identifier":     {},
	}

	for name, tt := range tests {
	}
}

type testProtocolIdentifier interface{}
type testTransportIdentifier interface{}

func (t testIdentifier) ProviderName() string
