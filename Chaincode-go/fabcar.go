
package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing sensor
type SmartContract struct {
	contractapi.Contract
}

// Sensor describes basic details of a sensor
type Sensor struct {
	Owner  string `json:"owner"`
	From   string `json:"from"`
	Time  string `json:"time"`
	Hash  string `json:"hash"`
	Address string `json:"address"`
}

// QueryResult structure used for handling result of query
type QueryResult struct {
	Key    string `json:"Key"`
	Record *Sensor
}

// InitLedger adds a base set of sensors to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	sensors := []Sensor{
		Sensor{Owner: "UserA", From: "CNC", Time: "1623057190444", Hash: "9218e8b37fabdc50bb6eca8597ffce22", Address: "swdbucket"},
		Sensor{Owner: "UserA", From: "CNC", Time: "1623057190444", Hash: "fe2b8c1ba153119466f9380b84553d87", Address: "swdbucket"},
		Sensor{Owner: "UserA", From: "CNC", Time: "1623057190444", Hash: "1e80de3f054f2d26b03f8c1c4de3f6b1", Address: "swdbucket"},
		Sensor{Owner: "UserA", From: "CNC", Time: "1623057190444", Hash: "3fc9480c6a6072853d478f4c04686936", Address: "swdbucket"},
		Sensor{Owner: "UserA", From: "CNC", Time: "1623057190444", Hash: "7730ca003aa8dd4e210685b4bbf95d56", Address: "swdbucket"},
		Sensor{Owner: "UserA", From: "CNC", Time: "1623057190444", Hash: "9218e8b37fabdc50bb6eca8597ffce22", Address: "swdbucket"},
		Sensor{Owner: "UserA", From: "CNC", Time: "1623057190444", Hash: "fe2b8c1ba153119466f9380b84553d87", Address: "swdbucket"},
		Sensor{Owner: "UserA", From: "CNC", Time: "1623057190444", Hash: "1e80de3f054f2d26b03f8c1c4de3f6b1", Address: "swdbucket"},
		Sensor{Owner: "UserA", From: "CNC", Time: "1623057190444", Hash: "3fc9480c6a6072853d478f4c04686936", Address: "swdbucket"},
		Sensor{Owner: "UserA", From: "CNC", Time: "1623057190444", Hash: "7730ca003aa8dd4e210685b4bbf95d56", Address: "swdbucket"},
	}

	for i, sensor := range sensors {
		sensorAsBytes, _ := json.Marshal(sensor)
		err := ctx.GetStub().PutState("SENSOR"+strconv.Itoa(i), sensorAsBytes)

		if err != nil {
			return fmt.Errorf("Failed to put to world state. %s", err.Error())
		}
	}

	return nil
}

// CreateSensor adds a new sensor to the world state with given details
func (s *SmartContract) CreateSensor(ctx contractapi.TransactionContextInterface, sensorNumber string, owner string, from string, time string, hash string, address string) error {
	sensor := Sensor{
		Owner:   owner,
		From:  from,
		Time: time,
		Hash:  hash,
		Address: address,
	}

	sensorAsBytes, _ := json.Marshal(sensor)

	return ctx.GetStub().PutState(sensorNumber, sensorAsBytes)
}

// QuerySensor returns the Sensor stored in the world state with given id
func (s *SmartContract) QuerySensor(ctx contractapi.TransactionContextInterface, sensorNumber string) (*Sensor, error) {
	sensorAsBytes, err := ctx.GetStub().GetState(sensorNumber)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if sensorAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", sensorNumber)
	}

	sensor := new(Sensor)
	_ = json.Unmarshal(sensorAsBytes, sensor)

	return sensor, nil
}

// QueryAllSensors returns all Sensors found in world state
func (s *SmartContract) QueryAllSensors(ctx contractapi.TransactionContextInterface) ([]QueryResult, error) {
	startKey := ""
	endKey := ""

	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	results := []QueryResult{}

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		sensor := new(Sensor)
		_ = json.Unmarshal(queryResponse.Value, sensor)

		queryResult := QueryResult{Key: queryResponse.Key, Record: sensor}
		results = append(results, queryResult)
	}

	return results, nil
}

// ChangeSensorOwner updates the owner field of Sensor with given id in world state
func (s *SmartContract) ChangeSensorOwner(ctx contractapi.TransactionContextInterface, sensorNumber string, newOwner string) error {
	sensor, err := s.QuerySensor(ctx, sensorNumber)

	if err != nil {
		return err
	}

	sensor.Owner = newOwner

	sensorAsBytes, _ := json.Marshal(sensor)

	return ctx.GetStub().PutState(sensorNumber, sensorAsBytes)
}

// ChangeSensorStatus updates the status of Sensor with given id in world state
func (s *SmartContract) ChangeSensorStatus(ctx contractapi.TransactionContextInterface, sensorNumber string,  newTime string, newHash string) error {
	sensor, err := s.QuerySensor(ctx, sensorNumber)

	if err != nil {
		return err
	}

	sensor.Time = newTime

	sensor.Hash = newHash

	sensorAsBytes, _ := json.Marshal(sensor)

	return ctx.GetStub().PutState(sensorNumber, sensorAsBytes)
}

func main() {

	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Error create fabsensor chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting fabsensor chaincode: %s", err.Error())
	}
}
