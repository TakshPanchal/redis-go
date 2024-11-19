#!/bin/bash
# Set the number of iterations
NUM_ITERATIONS=1000000

echo "Running $NUM_ITERATIONS ping requests to Redis..."

for i in $(seq 1 $NUM_ITERATIONS); do
    output=(redis-cli ping)
    echo "Iteration $i: $output"
done

echo "Test complete. You can now analyze the performance metrics."
