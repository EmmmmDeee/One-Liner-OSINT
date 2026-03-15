#!/bin/bash

# Example usage of the One-Liner OSINT tool

echo "=== One-Liner OSINT Examples ==="
echo

# Email analysis
echo "1. Email Analysis"
echo "   Command: osint email user@example.com"
./osint email user@example.com
echo
echo "---"
echo

# Domain analysis
echo "2. Domain Analysis"
echo "   Command: osint domain example.com"
./osint domain example.com
echo
echo "---"
echo

# IP analysis
echo "3. IP Address Analysis"
echo "   Command: osint ip 8.8.8.8"
./osint ip 8.8.8.8
echo
echo "---"
echo

# JSON output
echo "4. JSON Output Format"
echo "   Command: osint email test@gmail.com --output json"
./osint email test@gmail.com --output json
echo
echo "---"
echo

# Domain with subdomain enumeration
echo "5. Domain with Subdomain Enumeration"
echo "   Command: osint domain example.com --enumerate-subdomains"
echo "   Note: This will take longer as it checks common subdomains"
# Uncomment to run: ./osint domain example.com --enumerate-subdomains
echo "   (Skipped in this example for brevity)"
echo
echo "---"
echo

# Verbose output
echo "6. Verbose Output"
echo "   Command: osint email test@example.com --verbose"
./osint email test@example.com --verbose
echo
echo "---"
echo

echo "=== Examples Complete ==="
