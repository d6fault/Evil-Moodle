# Evil-Moodle: Automated Moodle LMS Security & Vulnerability Scanner

![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![MIT License](https://img.shields.io/badge/License-MIT-green.svg)

**Evil-Moodle** is a fast, lightweight Moodle LMS reconnaissance and 
misconfiguration scanner written in Go. Designed for penetration 
testers and security auditors, this CLI tool automates fingerprinting, 
misconfiguration discovery, and plugin enumeration via version.php 
probing using a configurable wordlist.

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
git clone https://github.com/d6fault/Evil-Moodle.git
# Navigate to the project directory
cd Evil-Moodle

# Build the optimized binary
go build -o evilmoodle
```

## Usage

Run a quick security audit against a target Moodle instance using the `scan` command:

```bash
./evilmoodle --url https://example.com
```

## Security Disclaimer

This tool is strictly intended for authorized security auditing, vulnerability assessment, and educational purposes. Do not run Evil-Moodle against targets without prior written consent. The developer assumes no liability for misuse or damage caused by this program.

## License

Distributed under the MIT License. See `LICENSE` for more information.

---
*Developed by [@d6fault](https://github.com).*
