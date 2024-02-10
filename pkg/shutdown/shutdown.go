package shutdown

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// HandleSignals реєструє обробники для сигналів завершення роботи
// та чекає на їх отримання або завершення контексту.
func HandleSignals(ctx context.Context, cancel context.CancelFunc) error {
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

	select {
	case <-signalCh:
		// Handle graceful shutdown
		return GracefulShutdown(ctx, cancel)
	case <-ctx.Done():
		// Context cancelled, shutdown initiated elsewhere
		return nil
	}
}

// GracefulShutdown викликається при отриманні сигналу завершення роботи.
func GracefulShutdown(ctx context.Context, cancel context.CancelFunc) error {
	cancel()
	// You may perform additional cleanup or graceful shutdown steps here
	return nil
}
