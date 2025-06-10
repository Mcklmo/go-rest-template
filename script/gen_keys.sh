if [ "$#" -ne 2 ]; then
    echo "Usage: $0 <private_key_filename> <public_key_filename>"
    exit 1
fi
PRIVATE_KEY_FILE=$1
PUBLIC_KEY_FILE=$2
openssl genpkey -algorithm RSA -out "$PRIVATE_KEY_FILE" -pkeyopt rsa_keygen_bits:2048
openssl rsa -pubout -in "$PRIVATE_KEY_FILE" -out "$PUBLIC_KEY_FILE"
echo "RSA keys generated and stored in $PRIVATE_KEY_FILE and $PUBLIC_KEY_FILE"
