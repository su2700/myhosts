// search.go
package cmd

import (
    "bufio"
    "fmt"
    "os"
    "strings"

    "github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
    Use:   "search [input]",
    Short: "Search for an IP or domain in /etc/hosts and return its counterpart",
    Long: `Search takes an IP address or domain name as input and searches /etc/hosts for a corresponding entry.
If an IP is provided, it returns the associated domain names. If a domain name is provided, it returns the IP.`,
    Args: cobra.ExactArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
        return searchHosts(args[0])
    },
}

func init() {
    rootCmd.AddCommand(searchCmd)
}

func searchHosts(input string) error {
    file, err := os.Open("/etc/hosts")
    if err != nil {
        return err
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    found := false

    for scanner.Scan() {
        line := scanner.Text()
        if strings.HasPrefix(line, "#") || strings.TrimSpace(line) == "" {
            continue // Skip comments and empty lines
        }

        parts := strings.Fields(line)
        if len(parts) < 2 {
            continue // Skip malformed lines
        }

        ip := parts[0]
        domains := parts[1:]

        if input == ip || contains(domains, input) {
            if input == ip {
                fmt.Println("Domains for IP", ip, ":", strings.Join(domains, ", "))
            } else {
                fmt.Println("IP for domain", input, ":", ip)
            }
            found = true
            break
        }
    }

    if !found {
        fmt.Println("Not found")
    }

    if err := scanner.Err(); err != nil {
        return err
    }

    return nil
}

func contains(slice []string, val string) bool {
    for _, item := range slice {
        if item == val {
            return true
        }
    }
    return false
}
