package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// showhostsCmd represents the command to show all entries in /etc/hosts.
var showhostsCmd = &cobra.Command{
	Use:   "showhosts",
	Short: "Show all IP and host entries in /etc/hosts",
	Long: `Showhosts displays all entries from the /etc/hosts file,
enabling the user to view all IP addresses and associated domain names.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return showHosts()
	},
}

func init() {
	rootCmd.AddCommand(showhostsCmd)
}

// showHosts reads and prints each line from the /etc/hosts file.
func showHosts() error {
	// Open /etc/hosts with permissions to read.
	file, err := os.Open("/etc/hosts")
	if err != nil {
		return err // Return an error if the file cannot be opened.
	}
	defer file.Close() // Ensure the file is closed after reading.

	scanner := bufio.NewScanner(file) // Create a scanner to read the file line by line.
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line) // Print each line to the console.
	}

	// Check if there were any errors during scanning, such as a malformed line.
	if err := scanner.Err(); err != nil {
		return err
	}

	return nil // Return nil if everything was successful.
}
