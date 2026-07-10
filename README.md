# Evil-Moodle
![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![MIT License](https://img.shields.io/badge/MIT-green?style=for-the-badge)

**Evil-Moodle** is a fast, lightweight Moodle LMS reconnaissance and 
misconfiguration scanner written in Go. Designed for penetration 
testers and security auditors, this CLI tool automates fingerprinting, 
misconfiguration discovery, and plugin enumeration via version.php 
probing using a configurable wordlist.

## Features

* **High-Performance Scanning**: Built in Go using goroutines for rapid multi-threaded target analysis.
* **Version Detection**: Accurately finger-prints core Moodle versions and deployment configurations.
* **Component Enumeration**: Identifies active plugins, themes, and blocks to map the attack surface.

## Prerequisites

Before building Evil-Moodle from source, ensure you have the following installed:

* **Go**: Version 1.21 or higher

## Installation

### From Source

Clone the official repository and compile the binary locally:

```bash
# Clone the repository
git clone https://github.com/d6fault/Evil-Moodle.git
# Navigate to the project directory
cd Evil-Moodle

# Build the optimized binary
go build -o evilmoodle
```

### Wordlist

For Plugin probing it will be regularly updated for new moodle plugins when they get added. The wordlist is located at ![Wordlist](https://github.com/d6fault/Evil-Moodle/blob/main/plugins.txt)

## Usage

```bash
./evilmoodle -h
Evil Moodle is a CLI tool for reconnaissance and security scanning of Moodle LMS instances.

It fingerprints the target, and enumerates installed
plugins via probing using a configurable wordlist.

Usage:
  evilmoodle --url <target> [flags]

Examples:
  evilmoodle --url https://moodle.example.com
  evilmoodle --url moodle.example.com --wordlist plugins.txt -v
  evilmoodle --url https://example.com --wordlist plugins.txt --rate 10 --delay 100

Flags:
      --url string        Moodle URL (required)
      --wordlist string   Path to plugins wordlist
  -v, --verbose           Verbose output
  -r, --rate int          Rate limit in requests per second (0 = unlimited)
  -d, --delay int         Delay between requests in milliseconds (default 0)

Usage:
  evilmoodle [flags]

Flags:
  -d, --delay int         Delay between requests in milliseconds
  -h, --help              help for evilmoodle
  -r, --rate int          Rate limit in requests per second (0 = unlimited)
      --url string        Moodle URL
  -v, --verbose           Verbose output
      --wordlist string   Path to plugins wordlist
```

## Security Disclaimer

This tool is strictly intended for authorized security auditing, vulnerability assessment, and educational purposes. Do not run Evil-Moodle against targets without prior written consent. The developer assumes no liability for misuse or damage caused by this program.

## Roadmap

Track the upcoming features and development milestones for Evil-Moodle. Feel free to open an issue to suggest new capabilities.

- [x] **Core CLI Framework**: Basic argument parsing and flag configuration.
- [x] **Version Detection**: Reliable extraction of core Moodle versions via standard artifacts.
- [ ] **Repository additions**: Get evilmoodle Added as a package in major Linux distros (Kali Linux, Black Arch, etc...)
- [ ] **CVE Matching**: Automate CVEs referencing via the Moodle version

*Developed by [@d6fault](https://github.com/d6fault).*
