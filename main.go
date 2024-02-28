package main

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

type AddressActivity struct {
	Address  string
	Activity int
}

func main() {
	_ = godotenv.Load()

	toBlockHex, err := getLastBlockNumber()
	if err != nil {
		log.Error().Err(err).Msg("")
		return
	}

	toBlock, err := strconv.ParseUint(toBlockHex[2:], 16, 64)
	if err != nil {
		log.Error().Err(err).Msg("")
		return
	}

	fromBlockHex := "0x" + strconv.FormatUint(toBlock-100, 16)

	topAddresses, err := processTransferEvents(fromBlockHex, toBlockHex)
	if err != nil {
		log.Error().Err(err).Msg("")
		return
	}

	result := "Top addresses: "

	for i, topAddress := range topAddresses {
		if i > 0 {
			result += ", "
		}

		result += strconv.Itoa(i+1) + ") " + topAddress.Address + ": " + strconv.Itoa(topAddress.Activity)
	}

	log.Info().Msg(result)
}

func getLastBlockNumber() (string, error) {
	response, err := sendRequest("eth_blockNumber", "")
	if err != nil {
		return "", err
	}

	result, ok := response.Result.(string)
	if !ok {
		return "", errors.New("Can't convert result to string")
	}

	log.Info().Msg("Last block number: " + result)

	return result, nil
}

func processTransferEvents(fromBlockHex, toBlockHex string) ([]AddressActivity, error) {
	transferTopic := "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef" // keccak256("Transfer(address,address,uint256)")
	params := fmt.Sprintf(`{"fromBlock": "%s", "toBlock": "%s", "topics": ["%s"]}`, fromBlockHex, toBlockHex, transferTopic)

	response, err := sendRequest("eth_getLogs", params)
	if err != nil {
		return nil, err
	}

	events, ok := response.Result.([]interface{})
	if !ok {
		return nil, errors.New("Can't convert result to []interface{}")
	}

	activities := make(map[string]int)

	for _, event := range events {
		event, ok := event.(map[string]interface{})
		if !ok {
			return nil, errors.New("Can't convert event log to map[string]interface{}")
		}

		topics, ok := event["topics"].([]interface{})
		if !ok {
			return nil, errors.New("Can't convert topics to []interface{}")
		}

		if len(topics) != 3 {
			continue
		}

		fromAddressTopic, ok := topics[1].(string)
		if !ok {
			return nil, errors.New("Can't convert topics[1] to string")
		}

		toAddressTopic, ok := topics[2].(string)
		if !ok {
			return nil, errors.New("Can't convert topics[2] to string")
		}

		fromAddress := "0x" + fromAddressTopic[26:]
		toAddress := "0x" + toAddressTopic[26:]

		activities[fromAddress]++
		activities[toAddress]++
	}

	objects := make([]AddressActivity, len(activities))

	for address, activity := range activities {
		// Use this if block for skipping zero address activity
		// if address == "0x0000000000000000000000000000000000000000" {
		// 	continue
		// }

		objects = append(objects, AddressActivity{Address: address, Activity: activity})
	}

	topAddresses := quickSort(objects, 0, len(objects)-1)[:5]

	return topAddresses, nil
}
