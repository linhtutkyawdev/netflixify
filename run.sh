#!/bin/ash

# Start the first process
dumb-init /doctron --config /doctron.yaml &

# Start the second process
echo "Running Netflixify Web Server on localhost:$PORT" && /netflixify &

# Wait for any process to exit
wait -n

# Exit with status of process that exited first
exit $?