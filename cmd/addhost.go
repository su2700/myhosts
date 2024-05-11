package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

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
	addhostCmd.Flags().StringP("ip", "i", "", "IP address to add to /etc/hosts (required)")
	addhostCmd.Flags().StringP("domain", "d", "", "Domain name to associate with the IP (required)")
	addhostCmd.MarkFlagRequired("ip")
	addhostCmd.MarkFlagRequired("domain")
}

func addToHosts(ip, domain string) error {
	file, err := os.OpenFile("/etc/hosts", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	exists := false
	var existingEntry string

	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
		if strings.Contains(line, ip) || strings.Contains(line, domain) {
			exists = true
			existingEntry = line
			break
		}
	}

	if exists {
		fmt.Printf("Warning: An entry with the IP '%s' or domain '%s' already exists: %s\n", ip, domain, existingEntry)
		fmt.Println("Do you want to update it? This will overwrite the existing entry. (yes/no)")
		var response string
		fmt.Scanln(&response)
		response = strings.ToLower(response) // Normalize the input to lowercase

		// Check for any valid affirmative response
		if response != "yes" && response != "y" {
			return fmt.Errorf("operation cancelled by the user")
		}

		// Remove the existing line and continue to write the new one
		var updatedLines []string
		for _, line := range lines {
			if line != existingEntry {
				updatedLines = append(updatedLines, line)
			}
		}
		lines = updatedLines
	}

	// Reopen the file for writing to truncate and update content
	if err := file.Truncate(0); err != nil {
		return err
	}
	if _, err := file.Seek(0, 0); err != nil {
		return err
	}

	for _, line := range lines {
		if _, err := fmt.Fprintln(file, line); err != nil {
			return err
		}
	}
	// Add the new entry
	if _, err := fmt.Fprintf(file, "%s\t%s\n", ip, domain); err != nil {
		return err
	}

	return nil
}
