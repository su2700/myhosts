# MyHosts CLI Tool

MyHosts is a command-line interface (CLI) tool designed to simplify the management of the `/etc/hosts` file on Unix-like systems. It allows users to easily add, delete, and list entries in the `/etc/hosts` file, providing a more intuitive and manageable approach to handling host entries, especially useful for developers and system administrators.

## Features

- **Add Hosts**: Quickly add new IP and domain associations to your hosts file.
- **Delete Hosts**: Remove entries from your hosts file using either an IP address or a domain name.
- **List Hosts**: Display all current entries in your hosts file.

## Installation

### Requirements

- Go 1.15 or higher
- Root permissions for managing `/etc/hosts`

### Building from Source

Clone the repository and build the project:

```bash
git clone https://github.com/su2700/myhosts.git
cd myhosts
go build -o myhosts
