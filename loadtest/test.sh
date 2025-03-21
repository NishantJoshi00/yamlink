#!/bin/bash

# Holistic test script for Waypoint load testing
# This script runs a complete test cycle including:
# - Server startup
# - Load testing with Locust
# - Results summary
# - Server shutdown
# - Cleanup

set -e  # Exit on any error

# Configuration
SERVER_PORT=8000
LOCUST_USERS=100
LOCUST_SPAWN_RATE=10
LOCUST_RUN_TIME="1m"
SERVER_LOG="waypoint-test.log"
TEST_RESULTS="loadtest-results.log"

# Color codes for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Helper functions
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check dependencies
check_dependencies() {
    log_info "Checking dependencies..."
    
    # Check for Go
    if ! command -v go &> /dev/null; then
        log_error "Go is not installed. Please install Go first."
        exit 1
    fi
    
    # Check for lsof (needed for port checking)
    if ! command -v lsof &> /dev/null; then
        log_warn "lsof is not installed. This is needed for port checking."
        if command -v apt-get &> /dev/null; then
            log_warn "Attempting to install lsof with apt-get..."
            sudo apt-get update && sudo apt-get install -y lsof || {
                log_error "Could not install lsof. Please install it manually or check port usage manually."
            }
        elif command -v yum &> /dev/null; then
            log_warn "Attempting to install lsof with yum..."
            sudo yum install -y lsof || {
                log_error "Could not install lsof. Please install it manually or check port usage manually."
            }
        elif command -v brew &> /dev/null; then
            log_warn "Attempting to install lsof with brew..."
            brew install lsof || {
                log_error "Could not install lsof. Please install it manually or check port usage manually."
            }
        else
            log_warn "Could not determine package manager. Please install lsof manually."
        fi
    fi
    
    # Check for Python/Locust
    if ! command -v locust &> /dev/null; then
        log_warn "Locust is not installed. Attempting to install..."
        pip install locust || {
            log_error "Failed to install Locust. Please install it manually."
            exit 1
        }
    fi
}

# Check if port is available
check_port() {
    if lsof -i:$SERVER_PORT > /dev/null 2>&1; then
        log_error "Port $SERVER_PORT is already in use."
        log_info "Checking for processes using port $SERVER_PORT..."
        lsof -i:$SERVER_PORT
        return 1
    fi
    return 0
}

# Start the Waypoint server
start_server() {
    # Check if port is available first
    check_port || {
        log_warn "Attempting to find an available port..."
        # Try to find an available port
        for PORT_ATTEMPT in $(seq 8001 8020); do
            SERVER_PORT=$PORT_ATTEMPT
            log_info "Trying port $SERVER_PORT..."
            if ! lsof -i:$SERVER_PORT > /dev/null 2>&1; then
                log_info "Port $SERVER_PORT is available."
                # Update config.yaml with new port
                sed -i.bak "s/port: [0-9]*/port: $SERVER_PORT/g" config.yaml
                log_info "Updated config.yaml to use port $SERVER_PORT"
                break
            fi
            # If we reach the last attempt and still no port is available
            if [ $PORT_ATTEMPT -eq 8020 ] && lsof -i:$SERVER_PORT > /dev/null 2>&1; then
                log_error "Could not find an available port in range 8000-8020. Please free up a port or modify the script."
                exit 1
            fi
        done
    }
    
    log_info "Starting Waypoint server on port $SERVER_PORT..."
    CONFIG_FILE=config.yaml go run ../cmd/waypoint/main.go > "$SERVER_LOG" 2>&1 &
    SERVER_PID=$!
    echo $SERVER_PID > waypoint.pid
    log_info "Server started with PID: $SERVER_PID"
    
    # Wait for server to start
    sleep 2
    
    # Verify server is running
    if ! ps -p $SERVER_PID > /dev/null; then
        log_error "Server failed to start. Check $SERVER_LOG for details."
        cat "$SERVER_LOG"
        exit 1
    fi
    
    # Verify server is actually listening on the port
    if ! lsof -i:$SERVER_PORT -p $SERVER_PID > /dev/null 2>&1; then
        log_error "Server started but is not listening on port $SERVER_PORT."
        log_error "Check $SERVER_LOG for details."
        cat "$SERVER_LOG"
        kill $SERVER_PID
        exit 1
    fi
    
    log_info "Server successfully started and is listening on port $SERVER_PORT"
}

# Stop the server
stop_server() {
    if [ -f waypoint.pid ]; then
        SERVER_PID=$(cat waypoint.pid)
        log_info "Stopping Waypoint server (PID: $SERVER_PID)..."
        kill $SERVER_PID 2>/dev/null || true
        rm waypoint.pid
    else
        log_warn "No waypoint.pid file found. Server may not be running."
    fi
}

# Run the load test
run_load_test() {
    log_info "Running load test with $LOCUST_USERS users at $LOCUST_SPAWN_RATE users/sec for $LOCUST_RUN_TIME..."
    log_info "Testing against http://localhost:$SERVER_PORT"
    
    # Checking if server is responsive before starting load test
    if ! curl -s --max-time 5 http://localhost:$SERVER_PORT/health > /dev/null; then
        log_error "Server is not responding to health check. Checking server logs..."
        cat "$SERVER_LOG"
        return 1
    fi
    
    log_info "Server is responsive, proceeding with load test..."
    locust --host=http://localhost:$SERVER_PORT \
           --headless \
           -u $LOCUST_USERS \
           -r $LOCUST_SPAWN_RATE \
           --run-time $LOCUST_RUN_TIME \
           --csv=test_results > "$TEST_RESULTS" 2>&1
    
    TEST_EXIT_CODE=$?
    
    if [ $TEST_EXIT_CODE -eq 0 ]; then
        log_info "Load test completed successfully."
    else
        log_error "Load test failed with exit code $TEST_EXIT_CODE."
        log_error "Load test output:"
        cat "$TEST_RESULTS"
        log_error "Server logs:"
        cat "$SERVER_LOG"
    fi
    
    return $TEST_EXIT_CODE
}

# Analyze and report test results
analyze_results() {
    log_info "Analyzing test results..."
    
    if [ -f test_results_stats.csv ]; then
        # Print some basic stats
        log_info "Test Statistics:"
        echo "--------------------------------------------"
        echo "Requests statistics:"
        head -1 test_results_stats.csv | tr ',' '\n' | nl
        grep -m1 "Aggregated" test_results_stats.csv | tr ',' '\n' | nl
        echo "--------------------------------------------"
        
        # Check for failures
        FAILURES=$(grep -m1 "Aggregated" test_results_stats.csv | cut -d',' -f8)
        if [ "$FAILURES" -gt 0 ]; then
            log_error "Test had $FAILURES failures."
            return 1
        else
            log_info "No failures detected in the test."
            return 0
        fi
    else
        log_warn "No test results file found. Cannot analyze results."
        return 1
    fi
}

# Clean up resources
cleanup() {
    log_info "Cleaning up..."
    stop_server
    
    # Restore original config.yaml if it was modified
    if [ -f config.yaml.bak ]; then
        log_info "Restoring original config.yaml..."
        mv config.yaml.bak config.yaml
    fi
    
    # Remove test result files
    if [ -f "$TEST_RESULTS" ]; then
        log_info "Removing test result files..."
        rm -f "$TEST_RESULTS" test_results_*.csv
    fi
    
    # Remove log file
    if [ -f "$SERVER_LOG" ]; then
        log_info "Removing server log file..."
        rm -f "$SERVER_LOG"
    fi
}

# Main execution
main() {
    log_info "Starting Waypoint holistic test..."
    
    # Setup trap to ensure cleanup on exit
    trap cleanup EXIT
    
    # Run test steps
    check_dependencies
    start_server
    run_load_test
    local test_result=$?
    
    if [ $test_result -eq 0 ]; then
        analyze_results
        test_result=$?
    fi
    
    # Report final status
    if [ $test_result -eq 0 ]; then
        log_info "🎉 Holistic test PASSED! 🎉"
    else
        log_error "❌ Holistic test FAILED! ❌"
        exit 1
    fi
}

# Run the main function
main "$@"
