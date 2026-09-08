package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"text/tabwriter"
)

var (
	host      = flag.String("host", "", "Host to bind to (default: all interfaces)")
	port      = flag.Int("port", 8080, "Port to listen on")
	directory = flag.String("directory", ".", "Directory to serve")
	help      = flag.Bool("help", false, "Show this help message")
)

func main() {
	flag.Parse()

	// Show help if requested
	if *help {
		showHelp()
		return
	}

	// Resolve the absolute path of the directory
	dir, err := filepath.Abs(*directory)
	if err != nil {
		log.Fatalf("Invalid directory: %v", err)
	}

	// Change to the specified directory
	if err := os.Chdir(dir); err != nil {
		log.Fatalf("Cannot change to directory %s: %v", dir, err)
	}

	// Create the handler
	handler := &Handler{
		Directory: dir,
		ThemeManager: NewThemeManager(),
	}

	// Set the host for binding
	bindHost := *host
	if bindHost == "" {
		bindHost = "0.0.0.0"
	}

	// Display host for logging (use 0.0.0.0 for binding but show localhost for convenience)
	displayHost := bindHost
	if displayHost == "0.0.0.0" {
		displayHost = "localhost"
	}

	// Start the server
	addr := fmt.Sprintf("%s:%d", *host, *port)
	fmt.Printf("Light Web Server @ http://%s:%d\n", displayHost, *port)
	fmt.Printf("Serving: %s\n", dir)
	
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatal(err)
	}
}

// showHelp displays the help message with proper flag formatting
func showHelp() {
	fmt.Printf("Light Web Server - A simple web server for Markdown files and source code\n\n")
	fmt.Printf("Usage:\n")
	fmt.Printf("  lightweb [OPTIONS]\n\n")
	fmt.Printf("Options:\n")
	
	// Use tabwriter for aligned output
	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(tw, "  --host HOST\tHost to bind to (default: all interfaces)\n")
	fmt.Fprintf(tw, "  --port PORT\tPort to listen on (default: 8080)\n")
	fmt.Fprintf(tw, "  --directory DIR\tDirectory to serve (default: current directory)\n")
	fmt.Fprintf(tw, "  --help\tShow this help message\n")
	tw.Flush()
	
	fmt.Printf("\n")
}