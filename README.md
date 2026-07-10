# Evil-Moodle: Automated Moodle LMS Security & Vulnerability Scanner

[![Go Version](https://shields.io)](https://golang.org)
[![License](https://shields.io)](LICENSE)

**Evil-Moodle** is a fast, lightweight, and concurrent Moodle LMS security scanner written in Go. Designed for penetration testers, security auditors, and system administrators, this CLI tool automates the discovery of misconfigurations, outdated plugins, and known vulnerabilities in Moodle learning management systems.

## Features

* **High-Performance Scanning**: Built in Go using goroutines for rapid multi-threaded target analysis.
* **Version Detection**: Accurately finger-prints core Moodle versions and deployment configurations.
* **Component Enumeration**: Identifies active plugins, themes, and blocks to map the attack surface.
* **Security Auditing**: Flags common security misconfigurations and exposed sensitive directories.

## Prerequisites

Before building Evil-Moodle from source, ensure you have the following installed:

* **Go**: Version 1.21 or higher

## Installation

### From Source

Clone the official repository and compile the binary locally:

```bash
# Clone the repository
git clone https://github.com

# Navigate to the project directory
cd Evil-Moodle

# Build the optimized binary
go build -o evilmoodle
```

## Usage

Run a quick security audit against a target Moodle instance using the `scan` command:

```bash
./evilmoodle scan --url https://example.com
```

### Command-Line Options

| Flag | Description | Example |
| :--- | :--- | :--- |
| `--url` | The target Moodle LMS URL to scan (Required) | `--url https://target.edu` |
| `--threads` | Number of concurrent requests (Default: 10) | `--threads 20` |
| `--output` | Save the scan results to a file | `--output report.txt` |

## Security Disclaimer

This tool is strictly intended for authorized security auditing, vulnerability assessment, and educational purposes. Do not run Evil-Moodle against targets without prior written consent. The developer assumes no liability for misuse or damage caused by this program.

## License

Distributed under the MIT License. See `LICENSE` for more information.

---
*Developed by [@d6fault](https://github.com).*
