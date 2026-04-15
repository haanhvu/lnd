package main

import (
	"fmt"
	"os"

	"github.com/btcsuite/btcd/wire"
	"github.com/jessevdk/go-flags"
	"github.com/lightningnetwork/lnd"
	"github.com/lightningnetwork/lnd/fn/v2"
	"github.com/lightningnetwork/lnd/lnwallet/chancloser"
	"github.com/lightningnetwork/lnd/lnwallet/types"
	"github.com/lightningnetwork/lnd/lnwire"
	"github.com/lightningnetwork/lnd/signal"
)

func main() {
	// Hook interceptor for os signals.
	shutdownInterceptor, err := signal.Intercept()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Load the configuration, and parse any command line options. This
	// function will also set up logging properly.
	loadedConfig, err := lnd.LoadConfig(shutdownInterceptor)
	if err != nil {
		if e, ok := err.(*flags.Error); !ok || e.Type != flags.ErrHelp {
			// Print error if not due to help request.
			err = fmt.Errorf("failed to load config: %w", err)
			_, _ = fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		// Help was requested, exit normally.
		os.Exit(0)
	}
	implCfg := loadedConfig.ImplementationConfig(shutdownInterceptor)
	//implCfg.AuxChanCloser = fn.Some[chancloser.AuxChanCloser](&mockAuxChanCloser{})

	// Call the "real" main in a nested manner so the defers will properly
	// be executed in the case of a graceful shutdown.
	if err = lnd.Main(
		loadedConfig, lnd.ListenerCfg{}, implCfg, shutdownInterceptor,
	); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

type mockAuxChanCloser struct{}

func (m *mockAuxChanCloser) ShutdownBlob(
	req types.AuxShutdownReq,
) (fn.Option[lnwire.CustomRecords], error) {

	return fn.None[lnwire.CustomRecords](), nil
}

func (m *mockAuxChanCloser) AuxCloseOutputs(
	desc types.AuxCloseDesc) (fn.Option[chancloser.AuxCloseOutputs], error) {

	// Implement later
	return fn.None[chancloser.AuxCloseOutputs](), nil
}

func (m *mockAuxChanCloser) FinalizeClose(desc types.AuxCloseDesc,
	closeTx *wire.MsgTx) error {

	return nil
}
