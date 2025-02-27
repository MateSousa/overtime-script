#!/bin/bash

# Check if correct number of arguments
if [ $# -ne 2 ]; then
  echo "Usage: $0 <ticket-url> <minutes>"
  exit 1
fi

TICKET_URL=$1
MINUTES=$2
TIMESTAMP=$(date +%Y%m%d%H%M%S)
CM_NAME="overtime-${TIMESTAMP}"

# Generate ConfigMap YAML
cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: ConfigMap
metadata:
  name: ${CM_NAME}
  namespace: personal-scripts
  labels:
    app: overtime
    created: "$(date +%Y-%m-%d)"
data:
  ticket_url: "${TICKET_URL}"
  minutes: "${MINUTES}"
EOF

echo "Created ConfigMap ${CM_NAME} with ticket ${TICKET_URL} and ${MINUTES} minutes"

