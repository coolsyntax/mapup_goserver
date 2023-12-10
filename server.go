package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"sync"
	"time"
)

type SortRequest struct {
	ToSort [][]int `json:"to_sort"` // This struct defines the expected JSON format of the request body.
}

type SortResponse struct {
	SortedArrays [][]int `json:"sorted_arrays"` // This struct defines the expected JSON format of the response body.
	TimeNs       int64   `json:"time_ns"`       // This will hold the time taken for sorting in nanoseconds.
}

func main() {
	http.HandleFunc("/process-single", handleSequential)     // Register the function to handle "/process-single" endpoint.
	http.HandleFunc("/process-concurrent", handleConcurrent) // Register the function to handle "/process-concurrent" endpoint.

	fmt.Println("Server listening on port 8000...") // Print a message indicating the server is started.
	http.ListenAndServe(":8000", nil)               // Start the server on port 8000.
}

func handleSequential(w http.ResponseWriter, r *http.Request) {
	// This function handles the "/process-single" endpoint for sequential sorting.

	// Decode the JSON request body into a SortRequest struct.
	var req SortRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error: %v", err)
		return
	}

	// Start measuring time.
	startTime := time.Now()

	// Create an empty slice to store the sorted arrays.
	sortedArrays := make([][]int, len(req.ToSort))

	// Iterate over each sub-array in the request.
	for i, arr := range req.ToSort {
		// Sort the sub-array using the `sort.Ints` function.
		sortedArrays[i] = sort.IntSlice(arr)
	}

	// Calculate the elapsed time in nanoseconds.
	elapsedTime := time.Since(startTime).Nanoseconds()

	// Create a SortResponse object with sorted arrays and elapsed time.
	resp := SortResponse{
		SortedArrays: sortedArrays,
		TimeNs:       elapsedTime,
	}

	// Encode the SortResponse object into JSON and send it as the response.
	json.NewEncoder(w).Encode(resp)
}

func handleConcurrent(w http.ResponseWriter, r *http.Request) {
	// This function handles the "/process-concurrent" endpoint for concurrent sorting using goroutines and channels.

	// Decode the JSON request body into a SortRequest struct.
	var req SortRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error: %v", err)
		return
	}

	// Start measuring time.
	startTime := time.Now()

	// Create a WaitGroup to manage concurrent goroutines.
	wg := sync.WaitGroup{}

	// Create an empty slice to store the sorted arrays.
	sortedArrays := make([][]int, len(req.ToSort))

	// Create a channel to receive the sorted sub-arrays from goroutines.
	ch := make(chan [][]int)

	// Iterate over each sub-array in the request.
	for i, arr := range req.ToSort {
		// Add 1 to the WaitGroup counter for each goroutine.
		wg.Add(1)

		// Launch a goroutine to sort the sub-array concurrently.
		go func(i int, arr []int) {
			// Sort the sub-array.
			sortedArrays[i] = sort.IntSlice(arr)

			// Signal completion of sorting the sub-array.
			wg.Done()
		}(i, arr)
	}

	// Launch another goroutine to wait for all sub-arrays to be sorted and receive them.
	go func() {
		// Wait for all goroutines to finish sorting.
		wg.Wait()

		// Send the sorted arrays through the channel.
		ch <- sortedArrays
	}()

	// Receive the sorted arrays from the channel.
	sortedArrays = <-ch

	// Calculate the elapsed time in nanoseconds.
	elapsedTime := time.Since(startTime).Nanoseconds()

	// Create a SortResponse object with sorted arrays and elapsed time.
	resp := SortResponse{
		SortedArrays: sortedArrays,
		TimeNs:       elapsedTime,
	}

	// Encode the SortResponse object into JSON and send it as the response.
	json.NewEncoder(w).Encode(resp)
}
