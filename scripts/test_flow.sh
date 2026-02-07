#!/bin/bash
set -e

BASE_URL="http://localhost:8081"
AUTH_URL="$BASE_URL/auth/token"

# Ensure jq is installed
if ! command -v jq &> /dev/null; then
    echo "jq could not be found, please install it to run this script"
    exit 1
fi

echo "Waiting for server to be ready..."
sleep 2

echo "---------------------------------------------------"
echo "1. Login as Admin"
ADMIN_TOKEN=$(curl -s -X POST $AUTH_URL -d '{"name":"admin","user_type":"admin"}' | jq -r .access_token)
echo "Admin Token: ${ADMIN_TOKEN:0:20}..."

echo "---------------------------------------------------"
echo "2. Registering Drone 'drone-001'"
DRONE_RESP=$(curl -s -X POST $BASE_URL/api/v1/drones \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -d '{"name":"drone-001"}')
echo "Response: $DRONE_RESP"

# Extract Drone ID (Assuming response contains the full drone object)
DRONE_ID=$(echo $DRONE_RESP | jq -r .id)
echo "Registered Drone ID: $DRONE_ID"

echo "---------------------------------------------------"
echo "3. Login as 'drone-001'"
DRONE_TOKEN=$(curl -s -X POST $AUTH_URL -d '{"name":"drone-001","user_type":"drone"}' | jq -r .access_token)
echo "Drone Token: ${DRONE_TOKEN:0:20}..."

echo "---------------------------------------------------"
echo "4. Creating an Order"
ORDER_RESP=$(curl -s -X POST $BASE_URL/api/v1/orders \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -d '{"origin_lat":10.0,"origin_lon":20.0,"dest_lat":10.1,"dest_lon":20.1}')
ORDER_ID=$(echo $ORDER_RESP | jq -r .id)
echo "Created Order ID: $ORDER_ID"

echo "---------------------------------------------------"
echo "5. Drone Heartbeat (Update Location)"
curl -s -X POST $BASE_URL/api/v1/drones/location \
  -H "Authorization: Bearer $DRONE_TOKEN" \
  -d '{"latitude":10.0,"longitude":20.0}' | jq .

echo "---------------------------------------------------"
echo "6. Drone Reserves Job"
RESERVE_RESP=$(curl -s -X POST $BASE_URL/api/v1/drones/jobs/reserve \
  -H "Authorization: Bearer $DRONE_TOKEN" \
  -d "{\"drone_id\":\"$DRONE_ID\"}")
echo "Reservation: $RESERVE_RESP"

echo "---------------------------------------------------"
echo "Verification Complete!"
