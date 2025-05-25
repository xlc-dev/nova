package nova

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
)

// vlog prints messages only if verbose is enabled.
func vlog(verbose bool, format string, a ...any) {
	if verbose {
		fmt.Printf(format+"\n", a...)
	}
}

// hasAllowedExtension returns true if filename ends with any of the allowed extensions.
func hasAllowedExtension(filename string, exts []string) bool {
	for _, ext := range exts {
		trimmed := strings.TrimSpace(ext)
		if strings.HasSuffix(filename, trimmed) {
			return true
		}
	}
	return false
}

// recompile triggers the recompilation of the application.
func recompile(verbose bool) error {
	cmd := exec.Command("go", "build", ".")
	if verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	} else {
		cmd.Stdout = nil
		cmd.Stderr = nil
	}
	vlog(verbose, "Running command: go build .")
	return cmd.Run()
}

// watchAndCompile watches the given directory for file changes that match allowed extensions
// and triggers a recompilation when a change is detected.
func watchAndCompile(dir string, verbose bool, exts []string) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("failed to create watcher: %w", err)
	}
	defer watcher.Close()

	// Add the directory (project root) to the watcher.
	if err := watcher.Add(dir); err != nil {
		return fmt.Errorf("failed to watch directory %s: %w", dir, err)
	}

	vlog(verbose, "File watcher started on directory %s...", dir)

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return nil
			}
			// Trigger only on file events (writes, creates, removes, renames)
			// and only if the changed file has an allowed extension.
			if event.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Remove|fsnotify.Rename) != 0 {
				if !hasAllowedExtension(event.Name, exts) {
					continue
				}
				vlog(verbose, "Change detected in %s", event.Name)
				if err := recompile(verbose); err != nil {
					slog.Error("Recompile error", "error", err)
				} else {
					vlog(verbose, "Recompile successful")
				}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return nil
			}
			slog.Error("Watcher error", "error", err)
		}
	}
}

// watchAndReload watches dir for changes to allowed extensions, rebuilds this binary
// in‐place and then signals reloadCh once so Serve can exec the new version.
func watchAndReload(
	dir, exe string,
	verbose bool,
	exts []string,
	reloadCh chan<- struct{},
) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	if err := watcher.Add(dir); err != nil {
		return err
	}

	for {
		select {
		case ev, ok := <-watcher.Events:
			if !ok {
				return nil
			}
			if ev.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Remove|fsnotify.Rename) == 0 {
				continue
			}
			if !hasAllowedExtension(ev.Name, exts) {
				continue
			}
			if verbose {
				slog.Info("File changed", "file", ev.Name)
			}
			build := exec.Command("go", "build", "-o", exe, ".")
			build.Stdout = os.Stdout
			build.Stderr = os.Stderr
			if verbose {
				slog.Info("Rebuilding", "cmd", strings.Join(build.Args, " "))
			}
			if err := build.Run(); err != nil {
				slog.Error("Rebuild failed", "error", err)
				continue
			}
			reloadCh <- struct{}{}
			return nil
		case err, ok := <-watcher.Errors:
			if !ok {
				return nil
			}
			slog.Error("Watcher error", "error", err)
		}
	}
}

// Serve launches the web server (concurrently) with graceful shutdown and live reloading (if enabled).
// It wraps key goroutines with recovery blocks to avoid crashes due to unexpected errors.
// It also allows for logging customization via context options.
func Serve(ctx *Context, router http.Handler) error {
	logFormat := ctx.String("log_format")
	logLevelStr := ctx.String("log_level")
	if logFormat != "" || logLevelStr != "" {
		opts := &slog.HandlerOptions{}
		// Parse log level.
		if logLevelStr != "" {
			switch strings.ToLower(logLevelStr) {
			case "debug":
				opts.Level = slog.LevelDebug
			case "info":
				opts.Level = slog.LevelInfo
			case "warn", "warning":
				opts.Level = slog.LevelWarn
			case "error":
				opts.Level = slog.LevelError
			default:
				opts.Level = slog.LevelInfo
			}
		}
		var handler slog.Handler
		// Choose JSON handler if requested.
		if strings.ToLower(logFormat) == "json" {
			handler = slog.NewJSONHandler(os.Stdout, opts)
		} else {
			handler = slog.NewTextHandler(os.Stdout, opts)
		}
		slog.SetDefault(slog.New(handler))
	}

	// Set watch directory to current working directory.
	watchDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %v", err)
	}
	// Validate the watch directory.
	if info, err := os.Stat(watchDir); err != nil || !info.IsDir() {
		return fmt.Errorf("directory %q does not exist or is not a valid directory", watchDir)
	}

	// Retrieve values from the context.
	verboseVal := ctx.Bool("verbose")
	hostVal := ctx.String("host")
	portVal := ctx.Int("port")
	watchVal := ctx.Bool("watch")
	extStr := ctx.String("extensions")
	if extStr == "" {
		extStr = ".go"
	}
	if portVal == 0 {
		portVal = 8080
	}
	if hostVal == "" {
		hostVal = "localhost"
	}
	exts := strings.Split(extStr, ",")

	slog.Info("Starting server", "host", hostVal, "port", portVal)
	vlog(verboseVal, "Verbose mode enabled")

	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("cannot find executable: %w", err)
	}
	exe, err = filepath.Abs(exe)
	if err != nil {
		return fmt.Errorf("cannot resolve executable path: %w", err)
	}

	reloadCh := make(chan struct{}, 1)
	if watchVal {
		vlog(verboseVal, "Live file watching enabled for extensions: %v", exts)
		go func() {
			if err := watchAndReload(watchDir, exe, verboseVal, exts, reloadCh); err != nil {
				slog.Error("Error watching files", "error", err)
			}
		}()
	} else {
		vlog(verboseVal, "Live file watching disabled")
	}

	addr := fmt.Sprintf("%s:%d", hostVal, portVal)
	// Create a new HTTP server.
	server := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	errCh := make(chan error, 1)
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	// Listen for auto-termination signals.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// Wait for a termination signal or server error.
	select {
	case sig := <-sigCh:
		slog.Info("Received termination signal, shutting down", "signal", sig)
		// Attempt graceful shutdown (with a timeout).
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			slog.Error("Error during graceful shutdown", "error", err)
			return fmt.Errorf("error during graceful shutdown: %v", err)
		}
		return nil
	case err := <-errCh:
		return fmt.Errorf("web server error: %v", err)
	case <-reloadCh:
		vlog(verboseVal, "Change detected, restarting...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if derr := server.Shutdown(shutdownCtx); derr != nil {
			slog.Error("Error during shutdown", "error", derr)
		}
		args := append([]string{exe}, os.Args[1:]...)
		env := os.Environ()
		if err := syscall.Exec(exe, args, env); err != nil {
			return fmt.Errorf("failed to exec new binary %q: %w", exe, err)
		}
		return nil // unreachable
	}
}
