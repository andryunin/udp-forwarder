package cmd

import (
	"errors"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/andryunin/udp-forwarder/internal"
	"github.com/spf13/cobra"
)

var listenAddr *net.UDPAddr
var targetAddr *net.UDPAddr
var waitResponseSeconds int

func parseAddr(s string) (*net.UDPAddr, error) {
	addr, err := net.ResolveUDPAddr("udp", s)

	if err != nil {
		return nil, err
	}

	if addr.IP == nil {
		return nil, errors.New("missing host in address")
	}

	return addr, nil
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "udp-forwarder [listen-addr] [target-addr]",
	Short: "Forwards all UDP packages to specified address",

	Args: func(cmd *cobra.Command, args []string) error {
		var err error

		if err = cobra.ExactArgs(2)(cmd, args); err != nil {
			return err
		}

		listenAddr, err = parseAddr(args[0])

		if err != nil {
			return err
		}

		targetAddr, err = parseAddr(args[1])

		if err != nil {
			return err
		}

		return nil
	},

	Run: func(cmd *cobra.Command, args []string) {
		err := internal.RunServer(listenAddr, targetAddr, time.Duration(waitResponseSeconds)*time.Second)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().IntVarP(
		&waitResponseSeconds, "timeout", "t", 60, "How many seconds wait for target response",
	)
}
