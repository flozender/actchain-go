package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

type SmartContract struct{

}

type Maintenance struct{

	Name string `json:"name"`
	Date string `json:"date"`
	PID string `json:"pid"`
	Equip_rq string `json:"equip_rq"`
	Staff_rq string `json:"staff_rq"`
	Staff_av string `json:"staff_av"`
	Staff_out string `json:"staff_out"`
	Equip_out string `json:"equip_out"`
	Cost string `json:"cost"`
	Inspec string `json:"inspec"`
}

func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response{

	function, args := APIstub.GetFunctionAndParameters()

	if function == "queryAllMaintenance" {
		return s.queryAllMaintenance(APIstub, args)
	} else if function == "initLedger" {
		return s.initLedger(APIstub)
	} else if function == "createMaintenance" {
		return s.createMaintenance(APIstub, args)
	} else if function == "queryMaintenance" {
		return s.queryMaintenance(APIstub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) queryAllMaintenance(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	startKey := "MTN1"
	endKey := "MTN999"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext(){
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- queryAllMaintenance:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	maint := []Maintenance{
		Maintenance{Name: "Gas Leak", Date: "12-01-19", PID: "PLN1249", Equip_rq: "1,4", Staff_rq: "6", Staff_av: "2", Staff_out: "1,2,1", Equip_out: "1,4", Cost: "99280", Inspec: "False"},
		Maintenance{Name: "Pipe Repairs", Date: "23-03-19", PID: "PLN4214", Equip_rq: "2,3", Staff_rq: "8", Staff_av: "7", Staff_out: "1", Equip_out: "2,3", Cost: "82442", Inspec: "True"},
		Maintenance{Name: "Regular Check Up", Date: "18-04-19", PID: "PLN2502", Equip_rq: "3,4", Staff_rq: "1", Staff_av: "1", Staff_out: "0", Equip_out: "0", Cost: "15000", Inspec: "True"},
	}

	i := 0
	for i < len(maint) {
		fmt.Println("i is ", i)
		MaintenanceAsBytes, _ := json.Marshal(maint[i])
		APIstub.PutState("MTN"+strconv.Itoa(i), MaintenanceAsBytes)
		fmt.Println("Added", maint[i])
		i = i + 1
	}

	return shim.Success(nil)
}

func (s *SmartContract) createMaintenance(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 11 {
		return shim.Error("Incorrect number of arguments. Expecting 11")
	}

	var maint = Maintenance{Name: args[1], Date: args[2], PID: args[3], Equip_rq: args[4], Staff_rq: args[5], Staff_av: args[6], Staff_out: args[7], Equip_out: args[8], Cost: args[9], Inspec: args[10]}

	maintAsBytes, _ := json.Marshal(maint)
	APIstub.PutState(args[0], maintAsBytes)

	return shim.Success(nil)
}

func (s *SmartContract) queryMaintenance(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	maintAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(maintAsBytes)
}

func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
