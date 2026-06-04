package cmd

import (
	"fmt"
	"log/slog"
	"net"
	"os"
	"time"

	"github.com/lmittmann/tint"
	"github.com/spf13/cobra"
)

var rootOptions struct {
	verbose bool
	count   int
	rate    int
	dst     net.IP
	port    int
	len     int
}

var rootCmd = &cobra.Command{
	Use:   "xudp",
	Short: "xudp is a simple UDP packet generator",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		options := &tint.Options{AddSource: true, TimeFormat: time.TimeOnly}
		if rootOptions.verbose {
			options.Level = slog.LevelDebug
		}
		logger := slog.New(tint.NewHandler(os.Stderr, options))
		slog.SetDefault(logger)
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := validateRootOptions(); err != nil {
			return err
		}

		dst := &net.UDPAddr{IP: rootOptions.dst, Port: rootOptions.port}
		conn, err := net.ListenUDP(udpNetwork(rootOptions.dst), nil)
		if err != nil {
			return fmt.Errorf("open UDP socket: %w", err)
		}
		defer conn.Close()

		var tick <-chan time.Time
		if rootOptions.rate > 0 {
			interval := time.Second / time.Duration(rootOptions.rate)
			tick = time.Tick(interval)
		}

		payload := make([]byte, rootOptions.len)
		failures := 0

		for i := 0; i < rootOptions.count; i++ {
			if tick != nil {
				<-tick
			}
			n, err := conn.WriteToUDP(payload, dst)
			if err != nil {
				slog.Error("Failed to send packet", "seq", i, "error", err)
				failures++
				continue
			}
			if n != len(payload) {
				slog.Error("Sent partial packet", "seq", i, "sent", n, "want", len(payload))
				failures++
				continue
			}
			slog.Debug("Sent packet", "seq", i, "size", n)
		}

		slog.Info("Finished sending packets", "total", rootOptions.count)
		if failures > 0 {
			return fmt.Errorf("failed to send %d of %d packets", failures, rootOptions.count)
		}
		return nil
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&rootOptions.verbose, "verbose", "v", false, "Verbose output")
	rootCmd.Flags().IntVarP(&rootOptions.count, "count", "c", 1000, "Number of packets to send")
	rootCmd.Flags().IntVarP(&rootOptions.rate, "rate", "r", 1000, "Packets per second")
	rootCmd.Flags().IPVarP(&rootOptions.dst, "dst", "d", net.IPv6loopback, "Destination IP address")
	rootCmd.Flags().IntVarP(&rootOptions.port, "port", "p", 4321, "Destination port")
	rootCmd.Flags().IntVarP(&rootOptions.len, "len", "l", 64, "Length of the UDP payload")
}

func validateRootOptions() error {
	if rootOptions.count < 0 {
		return fmt.Errorf("count must be greater than or equal to 0")
	}
	if rootOptions.rate < 0 {
		return fmt.Errorf("rate must be greater than or equal to 0")
	}
	if rootOptions.dst == nil {
		return fmt.Errorf("dst must be a valid IP address")
	}
	if rootOptions.port < 1 || rootOptions.port > 65535 {
		return fmt.Errorf("port must be between 1 and 65535")
	}
	if rootOptions.len < 0 {
		return fmt.Errorf("len must be greater than or equal to 0")
	}
	return nil
}

func udpNetwork(ip net.IP) string {
	if ip.To4() != nil {
		return "udp4"
	}
	return "udp6"
}
