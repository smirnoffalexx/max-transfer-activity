package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type GetBlockResponse struct {
	Id      string      `json:"id"`
	JSONRPC string      `json:"jsonrpc"`
	Result  interface{} `json:"result"`
}

func sendRequest(method string, params string) (GetBlockResponse, error) {
	requestBody := `{"jsonrpc": "2.0", "method": "` + method + `", "params": [` + params + `], "id": "getblock.io"}`

	req, err := http.NewRequest(http.MethodPost, os.Getenv("GET_BLOCK_URL"), strings.NewReader(requestBody))
	if err != nil {
		return GetBlockResponse{}, err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")

	client := http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return GetBlockResponse{}, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return GetBlockResponse{}, err
	}

	if resp.StatusCode > 202 {
		return GetBlockResponse{},
			errors.New("Invalid response status: " + strconv.Itoa(resp.StatusCode) + ". Body: " + string(body))
	}

	var data GetBlockResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return GetBlockResponse{}, err
	}

	return data, nil
}

func quickSort(arr []AddressActivity, low, high int) []AddressActivity {
	if low < high {
		var p int
		arr, p = partition(arr, low, high)
		arr = quickSort(arr, low, p-1)
		arr = quickSort(arr, p+1, high)
	}

	return arr
}

func partition(arr []AddressActivity, low, high int) ([]AddressActivity, int) {
	pivot := arr[high].Activity
	i := low

	for j := low; j < high; j++ {
		if arr[j].Activity > pivot {
			arr[i], arr[j] = arr[j], arr[i]
			i++
		}
	}

	arr[i], arr[high] = arr[high], arr[i]

	return arr, i
}
