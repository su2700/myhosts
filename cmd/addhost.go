package cmd

import (
    "bufio"
    "errors"
    "fmt"
    "os"
    "strings"

    "github.com/spf13/cobra"
)

// addhostCmd represents the addhost command
var addhostCmd = &cobra.Command{
    Use:   "addhost",
    Short: "Add an IP and domain to /etc/hosts",
    Long: `This command adds a new IP address and domain name to the /etc/hosts file,
ensuring there are no conflicts with existing entries.`,
    RunE: func(cmd *cobra.Command, args []string) error {
        ip, _ := cmd.Flags().GetString("ip")
        domain, _ := cmd.Flags().GetString("domain")
        return addToHosts(ip, domain)
    },
}

func init() {
    rootCmd.AddCommand(addhostCmd)

    // Here you will define your flags and configuration settings.
    addhostCmd.Flags().StringP("ip", "i", "", "IP address to add (required)")
    addhostCmd.Flags().StringP("domain", "d", "", "Domain name to add (required)")
    addhostCmd.MarkFlagRequired("ip")
    addhostCmd.MarkFlagRequired("domain")
}

func addToHosts(ip, domain string) error {
    file, err := os.OpenFile("/etc/hosts", os.O_RDWR|os.O_APPEND, 0644)
    if err != nil {
        return err
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()
        if strings.Contains(line, ip) && strings.Contains(line, domain) {
            return errors.New("entry already exists with the same IP and domain")
        }
    }

    _, err = file.WriteString(fmt.Sprintf("%s\t%s\n", ip, domain))
    return err
}
