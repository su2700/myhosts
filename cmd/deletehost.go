package cmd

import (
    "bufio"
    "fmt"
    "os"
    "strings"

    "github.com/spf13/cobra"
)

var deletehostCmd = &cobra.Command{
    Use:   "deletehost <ip|hostname>",
    Short: "Delete an IP or hostname from /etc/hosts",
    Long: `Deletehost removes entries from the /etc/hosts file based on the provided IP address or hostname.
You must specify either an IP or a hostname to locate and remove the corresponding entries.`,
    Args: cobra.ExactArgs(1), // Ensures exactly one argument is provided
    RunE: func(cmd *cobra.Command, args []string) error {
        return deleteFromHosts(args[0])
    },
}

func init() {
    rootCmd.AddCommand(deletehostCmd)
}

func deleteFromHosts(entry string) error {
    // Open /etc/hosts with read permissions
    file, err := os.Open("/etc/hosts")
    if err != nil {
        return err
    }
    defer file.Close()

    var lines, removedLines []string
    scanner := bufio.NewScanner(file)
    found := false

    // Scan each line and determine if it contains the entry
    for scanner.Scan() {
        line := scanner.Text()
        // Check if the line contains the entry; consider using strings.Fields for more accurate matching
        if strings.Contains(line, entry) {
            found = true
            removedLines = append(removedLines, line) // Save removed lines for confirmation
        } else {
            lines = append(lines, line) // Save lines that do not match
        }
    }

    if !found {
        return fmt.Errorf("no entry found matching %s", entry)
    }

    // Optionally, ask for confirmation before deletion
    fmt.Printf("The following entries will be removed:\n%s\n", strings.Join(removedLines, "\n"))
    fmt.Println("Are you sure you want to proceed? (y/n):")
    var response string
    fmt.Scanln(&response)
    if response != "y" {
        return fmt.Errorf("operation cancelled by user")
    }

    if err := scanner.Err(); err != nil {
        return err
    }

    // Rewrite the /etc/hosts file without the deleted entries
    return rewriteHostsFile(lines)
}

func rewriteHostsFile(lines []string) error {
    file, err := os.OpenFile("/etc/hosts", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
    if err != nil {
        return err
    }
    defer file.Close()

    for _, line := range lines {
        if _, err := file.WriteString(line + "\n"); err != nil {
            return err
        }
    }
    return nil
}