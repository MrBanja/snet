package snet

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

// ListenAndServe starts the server and listens for signals to shut down the server.
func ListenAndServe(
	ctx context.Context,
	server *http.Server,
	logger *slog.Logger,
	signals ...os.Signal,
) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	var (
		wg       sync.WaitGroup
		errorsCh = make(chan error, 2)
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer cancel()
		select {
		case <-ctx.Done():
			logger.Info("start server: ctx canceled")
			return
		default:
		}

		logger.Info("server listen and serve", slog.String("Addr", server.Addr))
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server listen and serve", "error", err)
			errorsCh <- err
		}
	}()

	// Run signal handler
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer cancel()

		c := make(chan os.Signal)
		signal.Notify(c, signals...)

		select {
		case <-ctx.Done():
		case sig := <-c:
			logger.Warn("received signal", slog.String("signal", sig.String()))
		}

		ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		logger.Warn("shutting down server")
		if err := server.Shutdown(ctx); err != nil {
			logger.Error("server shutdown", "error", err)
			errorsCh <- err
		}
		logger.Warn("shut down")
	}()

	wg.Wait()
	close(errorsCh)
	var err error = nil
	for e := range errorsCh {
		err = errors.Join(err, e)
	}

	return err
}
