// Copyright (c) 2015-2020 The usbtmc developers. All rights reserved.
// Project site: https://github.com/gotmc/usbtmc
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/gotmc/usbtmc"
	_ "github.com/gotmc/usbtmc/driver/google"
)

var (
	debugLevel uint
	address    string
)

func init() {
	// Get the debug level from CLI flag.
	const (
		defaultLevel = 1
		debugUsage   = "USB debug level"
	)
	flag.UintVar(&debugLevel, "debug", defaultLevel, debugUsage)
	flag.UintVar(&debugLevel, "d", defaultLevel, debugUsage+" (shorthand)")

	// Get VISA address from CLI flag.
	flag.StringVar(
		&address,
		"visa",
		"USB0::2391::1031::MY44035849::INSTR",
		"VISA address of Keysight 33220A",
	)
}

func main() {
	// Parse the flags
	flag.Parse()

	// Create new USBTMC context and new device.
	start := time.Now()
	usbCtx, err := usbtmc.NewContext()
	if err != nil {
		log.Fatalf("Error creating new USB context: %s", err)
	}
	usbCtx.SetDebugLevel(int(debugLevel))

	log.Printf("Using address: %s", address)
	fg, err := usbCtx.NewDevice(address)
	if err != nil {
		log.Fatalf("NewDevice error: %s", err)
	}
	log.Printf("%.2fs to create new device.", time.Since(start).Seconds())

	// Configure function generator using different write methods.
	ctx := context.Background()
	numCycles := 131
	period := 0.112
	fg.WriteString("*CLS\n")                              // Write using usbtmc.WriteString
	io.WriteString(fg, "burst:state off\n")               // Write using io.WriteString
	fg.Write([]byte("apply:sinusoid 2340, 0.1, 0.0\n"))   // Write using byte slice
	fmt.Fprintf(fg, "burst:internal:period %f\n", period) // Write using fmt.Fprint
	fg.Command(ctx, "burst:ncycles %d", numCycles)        // Write using usbtmc.Command
	fg.Command(ctx, "burst:state on")                     // Command appends a newline.

	queries := []string{"volt?", "freq?", "volt:offs?", "volt:unit?"}

	// Query using a write and then a read.
	for _, q := range queries {
		fg.WriteString(q)
		p := make([]byte, 512)
		bytesRead, err := fg.Read(p)
		if err != nil {
			log.Printf("Error reading: %v", err)
		} else {
			log.Printf("Read %d bytes for %s = %s", bytesRead, q, p[:bytesRead])
		}
	}

	// Query using the query method
	queryRange(ctx, fg, queries)

	// Close the function generator and USBTMC context and check for errors.
	err = fg.Close()
	if err != nil {
		log.Printf("error closing fg: %s", err)
	}
	err = usbCtx.Close()
	if err != nil {
		log.Printf("Error closing context: %s", err)
	}
}

func queryRange(ctx context.Context, fg *usbtmc.Device, r []string) {
	for _, q := range r {
		s, err := fg.Query(ctx, q)
		if err != nil {
			log.Printf("Error reading: %v", err)
		} else {
			log.Printf("Query %s = %s", q, s)
		}
	}
}
