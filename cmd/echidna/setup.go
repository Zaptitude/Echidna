package echidna

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/gookit/color"
)

var ctx context.Context

func createEchidnaDirs() error {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("Could not get current working directory in CreateDirs() with error\n%s", err)
	}

	current := filepath.Join(dir, "current")
	err = os.MkdirAll(current, os.ModePerm)
	if err != nil {
		return fmt.Errorf("Could not create directory %s in createdirs() with error\n%s", current, err)
	}
	inspect := filepath.Join(dir, "inspect")
	err = os.MkdirAll(inspect, os.ModePerm)
	if err != nil {
		return fmt.Errorf("Could not create directory %s in createdirs() with error\n%s", inspect, err)
	}

	return nil

}

func deleteCurrentDir() error {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("Could not get current working directory in deleteCurrentDir(). Please manually remove the folder")
	}

	current := filepath.Join(dir, "current")
	err = os.RemoveAll(current)
	if err != nil {
		return fmt.Errorf("Failed to remove Current/ directory and some subfiles. Please remove them manually")
	}

	return nil
}

func greeting() {
	color.LightCyan.Println(`
	   _     _     _             
          | |   (_)   | |            
  ___  ___| |__  _  __| |_ __   __ _ 
 / _ \/ __| '_ \| |/ _' | '_ \ / _' |
|  __/ (__| | | | | (_| | | | | (_| |
 \___|\___|_| |_|_|\__,_|_| |_|\__,_|`)

	color.Yellow.Println("Echidna Scanner running. Browse to http://127.0.0.1:8080 to view status.")
}

func setupCloseHandler() {
	// set up context for cancelling goroutines
	var cancel context.CancelFunc
	ctx, cancel = context.WithCancel(context.Background())
	// set up goroutine to catch CTRL+C and execute cleanup
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cancel()
		fmt.Println("Ctrl+C detected. Cancelling scanning goroutines and current web requests.")
		// Give the goroutines time to return and free up access to the zip files
		// so when we delete them we have access
		time.Sleep(2 * time.Second)
		fmt.Println("Attempting to remove current/ directory")
		err := deleteCurrentDir()
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
}

func errorHandler(ctx context.Context, errChan chan error) {

	for {
		err := <-errChan
		writeErr := ioutil.WriteFile("errors.txt", []byte(err.Error()), os.ModePerm)
		if writeErr != nil {
			log.Fatalf("Failed to log error %s into errors.txt. Closing program since error handling seems broken.", err)
		}
	}
}
