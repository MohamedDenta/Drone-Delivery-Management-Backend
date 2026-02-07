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
TIMESTAMP=$(date +%s)
DRONE_NAME="drone-$TIMESTAMP"
echo "2. Registering Drone '$DRONE_NAME'"
DRONE_RESP=$(curl -s -X POST $BASE_URL/api/v1/drones \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -d "{\"name\":\"$DRONE_NAME\"}")
echo "Response: $DRONE_RESP"

# Extract Drone ID
DRONE_ID=$(echo $DRONE_RESP | jq -r .id)
if [ "$DRONE_ID" == "null" ]; then
    echo "Drone might already exist or registration failed. Attempting to login anyway..."
fi

echo "---------------------------------------------------"
echo "3. Login as '$DRONE_NAME'"
DRONE_TOKEN=$(curl -s -X POST $AUTH_URL -d "{\"name\":\"$DRONE_NAME\",\"user_type\":\"drone\"}" | jq -r .access_token)
echo "Drone Token: ${DRONE_TOKEN:0:20}..."

echo "---------------------------------------------------"
echo "4. Drone Heartbeat & Location Update"
# Update location (this also sets heartbeat in Redis)
curl -s -X POST $BASE_URL/api/v1/drones/location \
  -H "Authorization: Bearer $DRONE_TOKEN" \
  -d '{"latitude":30.0,"longitude":31.0}' | jq .

echo "---------------------------------------------------"
echo "5. Verify Redis Cache (Location)"
# Optional: Try to check Redis if running in docker
if command -v docker &> /dev/null; then
    echo "Checking Redis for drone location..."
    REDIS_VAL=$(docker exec drone_redis redis-cli GET drone:$DRONE_ID:location || echo "N/A")
    echo "Redis Value: $REDIS_VAL"
fi

echo "---------------------------------------------------"
echo "6. Creating an Order (Async Flow)"
ORDER_RESP=$(curl -s -X POST $BASE_URL/api/v1/orders \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -d '{"origin_lat":30.0,"origin_lon":31.0,"dest_lat":30.1,"dest_lon":31.1}')
ORDER_ID=$(echo $ORDER_RESP | jq -r .id)
echo "Created Order ID: $ORDER_ID (Published to RabbitMQ)"

echo "---------------------------------------------------"
echo "7. Waiting for Async Dispatcher..."
# Poll for status change to RESERVED
MAX_RETRIES=10
COUNT=0
STATUS="PENDING"

while [ "$STATUS" == "PENDING" ] && [ $COUNT -lt $MAX_RETRIES ]; do
    echo "Polling order status (Attempt $((COUNT+1)))..."
    ORDER_STATUS_RESP=$(curl -s -X GET $BASE_URL/api/v1/orders/$ORDER_ID \
      -H "Authorization: Bearer $ADMIN_TOKEN")
    
    STATUS=$(echo $ORDER_STATUS_RESP | jq -r .status)
    DRONE_ID_ASSIGNED=$(echo $ORDER_STATUS_RESP | jq -r .drone_id)
    
    echo "Current Status: $STATUS, Assigned Drone: $DRONE_ID_ASSIGNED"
    
    if [ "$STATUS" == "RESERVED" ]; then
        echo "Order successfully assigned to drone!"
        break
    fi
    
    sleep 2
    COUNT=$((COUNT+1))
done

if [ "$STATUS" == "PENDING" ]; then
    echo "Timed out waiting for order assignment."
    exit 1
fi

echo "Verification Complete! Check the server logs for 'Successfully assigned order' messages."
