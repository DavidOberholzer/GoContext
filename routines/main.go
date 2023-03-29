package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type routine struct {
	timeout time.Duration
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}

	routines := []routine{
		{timeout: 5 * time.Second},
		{timeout: 2 * time.Second},
		{timeout: 0 * time.Second},
		{timeout: 0 * time.Second},
		{timeout: 3 * time.Second},
	}

	for i, r := range routines {
		wg.Add(1)
		go func(ctx context.Context, routineNum int, r routine) {
			defer wg.Done()

			switch {
			case r.timeout != 0:
				timeoutCtx, timeoutCancel := context.WithTimeout(ctx, r.timeout)
				defer timeoutCancel()

				<-timeoutCtx.Done()
			default:
				<-ctx.Done()
			}

			fmt.Printf("Goodbye to routine %d\n", routineNum)
		}(ctx, i, r)
	}

	sigs := make(chan os.Signal)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)

	// Wait for 10s to let routines go.
	ctxWait, waitCancel := context.WithTimeout(ctx, 10*time.Second)
	defer waitCancel()

	select {
	case <-ctxWait.Done():
		fmt.Println("Timed out!")
	case <-sigs:
		fmt.Println("Program killed")
	}

	cancel()
	// Wait for ALL routines to report done!
	wg.Wait()
}
