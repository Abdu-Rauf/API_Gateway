# API Gateway: Node.js vs. Go

## Overview
This project is an API Gateway implementation designed specifically to test and compare the fundamental architectural differences between Node.js and Go. The primary goal is to benchmark how well the Node.js single-threaded Event Loop performs against Go's multi-threaded concurrency model (Goroutines) under various load conditions.

## Architectural Comparison
By stressing both implementations with high-throughput traffic, this project evaluates:
*   **Node.js (Event Loop):** How queueing theory affects latency when a single core is saturated with CPU-bound and I/O-bound tasks.
*   **Go (Goroutines):** How lightweight threads distribute work across multiple CPU cores, and how that impacts memory footprints and concurrent connection scaling.

## Project Details & Workload
To ensure the benchmarks reflect a real-world API Gateway rather than just a simple "Hello World" test, the following features are implemented:

*   **JWT Verification:** Simulates heavy, CPU-bound cryptographic mathematical operations.
*   **Token Bucket Algorithm:** Implements robust rate-limiting, which is a core responsibility of any production API Gateway.
*   **Redis Integration:** The Token Bucket state is managed via Redis. This was done to simulate a real-world use case, reflecting how actual API Gateways handle distributed rate-limiting and network I/O downstreams.

## Metrics Evaluated
Through continuous load testing (using tools like `wrk`), the project tracks:
*   **Requests Per Second (RPS):** The absolute throughput limits.
*   **Latency:** How response times degrade as concurrent connections increase.
*   **Resource Utilization:** CPU saturation and RAM consumption.