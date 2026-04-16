#!/bin/bash

OUTPUT_FILE="results.csv"
echo "Server,Connections,ReqPerSec,Latency,Timeouts,ConnErr" > $OUTPUT_FILE

# 1. Build the Go binary
echo "Compiling Go binary..."
cd go
go build -o apigateway_go server.go verifyToken.go
cd ..

# Helper function to run tests and parse output
run_tests() {
    local server_name=$1
    local port=$2

    for conns in 10 100 1000 5000; do
        echo "Testing $server_name with $conns connections..."
        
        # Updated path to point to auth/wrk_script.lua
        OUTPUT=$(wrk -t8 -c$conns -d30s -s auth/wrk_script.lua http://localhost:$port)
        
        # Extract Latency (avg) and Requests/sec using awk
        LATENCY=$(echo "$OUTPUT" | grep "Latency" | awk '{print $2}')
        RPS=$(echo "$OUTPUT" | grep "Requests/sec:" | awk '{print $2}')
        
        # Extract errors (Default to 0 if there are none)
        TIMEOUTS=$(echo "$OUTPUT" | grep "Socket errors" | awk -F'timeout ' '{print $2}' | tr -d ' ')
        CONN_ERRORS=$(echo "$OUTPUT" | grep "Socket errors" | awk -F'connect ' '{print $2}' | awk -F',' '{print $1}')
        
        [ -z "$TIMEOUTS" ] && TIMEOUTS=0
        [ -z "$CONN_ERRORS" ] && CONN_ERRORS=0
        
        # Append to CSV
        echo "$server_name,$conns,$RPS,$LATENCY,$TIMEOUTS,$CONN_ERRORS" >> $OUTPUT_FILE
        sleep 2
    done
}


# 2. Start Node.js, benchmark, and kill
echo "Starting Node.js Server on 1 Core..."
cd node
taskset -c 0 node server.js &
NODE_PID=$!
cd .. # RETURN TO ROOT before running tests
sleep 3 # Wait for server and Redis to be ready

run_tests "NodeJS" 3000

echo "Stopping Node.js Server..."
kill $NODE_PID
wait $NODE_PID 2>/dev/null

# 3. Start Go, benchmark, and kill
echo "Starting Go Server on 1 Core..."
cd go # Now we are at root, so just cd go
taskset -c 0 ./apigateway_go &
GO_PID=$!
cd .. # RETURN TO ROOT before running tests
sleep 3 # Wait for server and Redis to be ready

run_tests "Go" 8080

echo "Stopping Go Server..."
kill $GO_PID
wait $GO_PID 2>/dev/null

echo "Benchmarking complete! Data saved to $OUTPUT_FILE"
cat $OUTPUT_FILE