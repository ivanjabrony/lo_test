package logger

import (
	"context"
	"fmt"
	"io"
	"log"
	"sync"
	"time"
)

type AsyncLogger struct {
	w  io.Writer
	in chan string
	wg *sync.WaitGroup
}

// NewAsync creates a new AsyncLogger.
//
// AsyncLogger main idea is to have a separate goroutine that reads logs from a channel
// and writes it to a writer.
// Standart writer is a os.Stdout
func NewAsync(ctx context.Context, w io.Writer) AsyncLogger {
	in := make(chan string)
	wg := sync.WaitGroup{}
	logger := AsyncLogger{w, in, &wg}

	wg.Add(1)

	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case data := <-in:
				_, err := w.Write([]byte(data))
				if err != nil {
					log.Printf("error while writing logs: %v", err)
				}
			}
		}
	}()

	go func() {
		wg.Wait()
		close(in)
	}()

	return logger
}

// Log passes data into a main channel
func (al AsyncLogger) Log(format string, info ...any) {
	al.wg.Add(1)

	go func() {
		time := time.Now().Format("2006-01-02 15:04:05")
		log := fmt.Sprintf(format, info...)
		al.in <- fmt.Sprintf("[%v] == %v\n", time, log)
		al.wg.Done()
	}()
}
