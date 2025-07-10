#!/bin/bash

# Generate VRF Key Script
# This script generates a new secp256k1 private key for use with OmniVRF

set -e

# Check if OpenSSL is available
if ! command -v openssl &> /dev/null; then
    echo "Error: OpenSSL is required but not installed."
    echo "Please install OpenSSL and try again."
    exit 1
fi

# Generate key file name with timestamp
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
KEY_FILE="vrf_key_${TIMESTAMP}.pem"
HEX_FILE="vrf_key_${TIMESTAMP}.hex"

echo "Generating new VRF private key..."

# Generate new secp256k1 key
openssl ecparam -genkey -name secp256k1 -out "${KEY_FILE}"

if [ $? -ne 0 ]; then
    echo "Error: Failed to generate private key"
    exit 1
fi

echo "✓ Generated private key: ${KEY_FILE}"

# Extract private key in hex format
openssl ec -in "${KEY_FILE}" -text -noout | grep priv -A 3 | tail -n +2 | tr -d '\n[:space:]:' > "${HEX_FILE}"

if [ $? -ne 0 ]; then
    echo "Error: Failed to extract hex private key"
    exit 1
fi

echo "✓ Extracted hex private key: ${HEX_FILE}"

# Display the hex key
HEX_KEY=$(cat "${HEX_FILE}")
echo "✓ Your VRF private key (hex): ${HEX_KEY}"

echo ""
echo "IMPORTANT SECURITY NOTES:"
echo "========================"
echo "1. Store this private key securely"
echo "2. Never share or commit this key to version control"
echo "3. Use environment variables or secure key management systems"
echo "4. Consider using AWS Secrets Manager or similar for production"
echo ""
echo "To use this key with OmniVRF:"
echo "export VRF_PRIVATE_KEY=${HEX_KEY}"
echo ""
echo "Files created:"
echo "- ${KEY_FILE} (PEM format - keep secure)"
echo "- ${HEX_FILE} (Hex format - for environment variable)"
echo ""
echo "You can delete the .pem file after copying the hex key to a secure location."