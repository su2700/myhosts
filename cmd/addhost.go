// Package cmd contains all command-related logic for the CLI application.
package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// addhostCmd represents the command for adding an IP address and domain name to /etc/hosts.
var addhostCmd = &cobra.Command{
	Use:   "addhost",
	Short: "Add an IP and domain to /etc/hosts",
	Long: `Addhost adds a new IP address and domain name to the /etc/hosts file,
ensuring there are no conflicts with existing entries. It checks for either the
same IP or domain name already present in the file to prevent networking issues.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ip, _ := cmd.Flags().GetString("ip")
		domain, _ := cmd.Flags().GetString("domain")
		return addToHosts(ip, domain)
	},
}

func init() {
	rootCmd.AddCommand(addhostCmd)

	// Define flags for the IP address and domain name. Both flags are required.
	addhostCmd.Flags().StringP("ip", "i", "", "IP address to add to /etc/hosts (required)")
	addhostCmd.Flags().StringP("domain", "d", "", "Domain name to associate with the IP (required)")
	addhostCmd.MarkFlagRequired("ip")
	addhostCmd.MarkFlagRequired("domain")
}

// addToHosts adds the specified IP and domain to the /etc/hosts file after checking for conflicts.
func addToHosts(ip, domain string) error {
	// Open /etc/hosts with permissions to read and append. Create if it does not exist.
	file, err := os.OpenFile("/etc/hosts", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return err // Return an error if the file cannot be opened.
	}
	defer file.Close() // Ensure the file is closed after writing.

	scanner := bufio.NewScanner(file) // Create a scanner to read the file line by line.
	for scanner.Scan() {
		line := scanner.Text()
		// Check if the line contains either the provided IP or domain.
		if strings.Contains(line, ip) || strings.Contains(line, domain) {
			return errors.New("entry already exists with the same IP or domain")
		}
	}

	// Write the new IP and domain to the file if no conflicts were found.
	_, err = file.WriteString(fmt.Sprintf("%s\t%s\n", ip, domain))
	if err != nil {
		return err // Return any errors that occur while writing to the file.
	}

	return nil // Return nil on successful addition.
}
