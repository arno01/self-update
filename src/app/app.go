// Source: https://stackoverflow.com/a/18914180

package main

import (
        "os"
        "os/signal"
        "syscall"
	"time"
	"fmt"
	"log"
)

// We make sigHandler receive a channel on which we will report the value of var quit
func sigHandler(q chan bool) {
        var quit bool

        c := make(chan os.Signal, 1)
        signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGUSR1, syscall.SIGUSR2)

        // foreach signal received
        for signal := range c {
                log.Printf("Signal received: "+signal.String())

                switch signal {
                case syscall.SIGINT, syscall.SIGTERM:
                        quit = true
                case syscall.SIGHUP:
                        quit = false
                case syscall.SIGUSR1, syscall.SIGUSR2:
                        quit = false
                }

                if quit {
                        quit = false
                        //              closeDb()
                        log.Printf("Terminating...")

			// Simulate slow termination
                        log.Printf("intentionally sleep for 5 seconds...")
			time.Sleep(5 * time.Second)
                        log.Printf("reached graceful stop!")
                        //              closeLog()
                        os.Exit(0)
                }
                // report the value of quit via the channel
                q <- quit
        }
}

var (
    version = "dev"
    commit  = "none"
    date    = "unknown"
    builtBy = "unknown"
)

func main() {

        fmt.Printf("app version: %s, commit: %s, built at: %s by: %s\n", version, commit, date, builtBy)

        // init two channels, one for the signals, one for the main loop
        sig := make(chan bool)
        loop := make(chan error)

        // start the signal monitoring routine
        go sigHandler(sig)


        // while vat quit is false, we keep going
        for quit := false; !quit; {
                // we start the main loop code in a goroutine
                go func() {
                        // Main loop code here
                        // we can report the error via the chan (here, nil)
			time.Sleep(1 * time.Second)
                        loop <- nil
                }()

                // We block until either a signal is received or the main code finished
                select {
                // if signal, we affect quit and continue with the loop
                case quit = <-sig:
                // if no signal, we simply continue with the loop
                case <-loop:
                }
        }
}
