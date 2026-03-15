# OSINT Examples

This directory contains example usage of the One-Liner OSINT tool.

## Basic Usage Examples

### Email Analysis

```bash
# Basic email analysis
osint email user@example.com

# Email analysis with JSON output
osint email user@example.com --output json

# Email analysis with breach checking (requires API key)
osint email user@example.com --check-breaches
```

### Domain Analysis

```bash
# Basic domain analysis
osint domain example.com

# Domain analysis with subdomain enumeration
osint domain example.com --enumerate-subdomains

# Domain analysis with CSV output
osint domain example.com --output csv
```

### IP Address Analysis

```bash
# Basic IP analysis
osint ip 8.8.8.8

# IP analysis with geolocation
osint ip 8.8.8.8 --geolocation

# IP analysis with JSON output
osint ip 8.8.8.8 --output json
```

### Image Metadata Extraction

```bash
# Extract metadata from local image
osint image /path/to/image.jpg

# Extract metadata from remote image
osint image https://example.com/image.jpg

# Extract with GPS coordinates
osint image /path/to/image.jpg --extract-gps
```

## Advanced Examples

### Custom Configuration

```bash
# Use custom config file
osint email user@example.com --config ~/.my-osint-config.yaml

# Increase worker count for parallel processing
osint domain example.com --workers 20

# Set custom timeout
osint ip 8.8.8.8 --timeout 60
```

### Output Formats

```bash
# Text output (default)
osint email user@example.com --output text

# JSON output
osint email user@example.com --output json

# CSV output
osint email user@example.com --output csv

# Table output
osint email user@example.com --output table
```

### Verbose Logging

```bash
# Enable verbose logging for debugging
osint email user@example.com --verbose

# Disable colored output
osint email user@example.com --no-color
```

## API Integration Examples

### Have I Been Pwned

```bash
# Set API key in config file
echo "api_keys:
  haveibeenpwned: YOUR_API_KEY" > ~/.osint.yaml

# Check email for breaches
osint email user@example.com --check-breaches
```

### Multiple Analyses

```bash
# Analyze multiple targets
for email in user1@example.com user2@example.com user3@example.com; do
    osint email $email --output json >> results.json
done

# Analyze multiple domains
for domain in example.com test.com sample.org; do
    osint domain $domain --output csv >> domains.csv
done
```

## Scripting Examples

### Batch Processing

```bash
#!/bin/bash

# Read emails from file and analyze each
while IFS= read -r email; do
    echo "Analyzing: $email"
    osint email "$email" --output json >> batch_results.json
done < emails.txt
```

### Combined Analysis

```bash
#!/bin/bash

TARGET="example.com"

echo "Analyzing domain: $TARGET"
osint domain $TARGET --output json > domain_results.json

# Extract IPs and analyze them
IPS=$(osint domain $TARGET --output json | jq -r '.[] | select(.type=="A") | .data.records[]')

for ip in $IPS; do
    echo "Analyzing IP: $ip"
    osint ip $ip --output json >> ip_results.json
done
```

## Integration with Other Tools

### With jq (JSON processing)

```bash
# Extract specific fields from JSON output
osint email user@example.com --output json | jq '.[] | select(.type=="format") | .data'

# Get all DNS records
osint domain example.com --output json | jq '.[] | select(.type=="A" or .type=="MX")'
```

### With grep

```bash
# Filter results
osint domain example.com | grep -A 5 "MX"

# Extract email addresses
osint domain example.com | grep -oE '[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}'
```

### Pipeline Examples

```bash
# Analyze domain and save results
osint domain example.com --output json | \
    jq '.[] | {type: .type, data: .data}' > results.json

# Extract IPs and geolocate them
osint domain example.com --output json | \
    jq -r '.[] | select(.type=="A") | .data.records[]' | \
    xargs -I {} osint ip {} --output json
```

## Performance Tuning

```bash
# Increase workers for parallel processing
osint domain example.com --workers 50

# Decrease timeout for faster failures
osint ip 192.168.1.1 --timeout 5

# Use caching for repeated queries
# (Cache is automatic, stored in .osint_cache/)
```

## Best Practices

1. **Rate Limiting**: Use appropriate worker counts to avoid rate limiting
2. **API Keys**: Store API keys in config file, not command line
3. **Output Format**: Use JSON for programmatic processing
4. **Batch Processing**: Process multiple targets efficiently with loops
5. **Error Handling**: Check exit codes and handle errors appropriately

## Running the Examples

```bash
# Make the example script executable
chmod +x examples/basic_usage.sh

# Run the examples
./examples/basic_usage.sh
```

## Contributing Examples

If you have useful examples, please contribute them! Examples should be:
- Clear and well-documented
- Practical and useful
- Include expected output
- Follow best practices
