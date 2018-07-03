package main

import (
	"encoding/json"
	"errors"
//	"strconv"
//	"bytes"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
    pb "github.com/hyperledger/fabric/protos/peer"
	//"github.com/hyperledger/fabric/bccsp/sw"
)

// ============================================================================================================================
//  Get Contract - get a Contract info from ledger
//
//  Inputs - string
//  MessageIdentification    
//       "id0000"
// ============================================================================================================================
func getContractInfo(stub shim.ChaincodeStubInterface, id string) (CIPS, error) {
	var contract CIPS
	contractAsBytes, err := stub.GetState(id)                  //getState retreives a key/value from the ledger
	if err != nil {                                          //this seems to always succeed, even if key didn't exist
		return contract, errors.New("Failed to find contract - " + id)
	}
	json.Unmarshal(contractAsBytes, &contract)                   //un stringify it aka JSON.parse()

	if contract.MessageIdentification != id {                                     //test if marble is actually here or just nil
		return contract, errors.New("Contract does not exist - " + id)
	}

	return contract, nil
}


// ============================================================================================================================
// Get history of asset
//
// Shows Off GetHistoryForKey() - reading complete history of a key/value
//
// Inputs - Array of strings
//  0
//  MessageIdentification
//  "id0000"
// ============================================================================================================================
func getHistory(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	type AuditHistory struct {
		TxId    string     `json:"txId"`
		Value   CIPS       `json:"value"`
	}
	var history []AuditHistory;
	var contract CIPS

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	contractId := args[0]
	fmt.Printf("- start getHistoryForContract: %s\n", contractId)

	// Get History
	resultsIterator, err := stub.GetHistoryForKey(contractId)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	for resultsIterator.HasNext() {
		historyData, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		fmt.Printf("- historyData: %+v\n", historyData)
		fmt.Printf("- historyData.Value: %s\n", historyData.Value)
		var tx AuditHistory
		tx.TxId = historyData.TxId                     //copy transaction id over
		json.Unmarshal(historyData.Value, &contract)     //un stringify it aka JSON.parse()
		fmt.Printf("- contract111111111: %+v\n", contract)
		if historyData.Value == nil {                  //contract has been deleted
//			fmt.Printf("- 111111111111111111111111111111111\n")
			var emptyContract CIPS
			tx.Value = emptyContract                 //copy nil contract
		} else {
//			fmt.Printf("- 222222222222222222222222222222222\n")
			json.Unmarshal(historyData.Value, &contract) //un stringify it aka JSON.parse()
			tx.Value = contract                      //copy contract over
			fmt.Printf("- contract22222222: %+v\n", contract)
		}
		history = append(history, tx)              //add this tx to the list
	}
	fmt.Printf("- getHistoryForContract returning:\n%s", history)

	//change to array of bytes
	historyAsBytes, _ := json.Marshal(history)     //convert to array of bytes
	return shim.Success(historyAsBytes)
}


// ============================================================================================================================
//  Get a Contract state - get a Contract state from ledger
//
//  Inputs - string
//  MessageIdentification    
//       "id0000"
// ============================================================================================================================
func getCtrctStateById(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	
	contractAsBytes, err := stub.GetState(args[0])                  //getState retreives a key/value from the ledger
	if err != nil {                                          //this seems to always succeed, even if key didn't exist
		return shim.Error(err.Error())
	}
	
	return shim.Success(contractAsBytes)
}



// ============================================================================================================================
//  Get a remittance state by  InstructingParty/Contractstatus
//
//  Inputs - string
//		1,						2					
//  InstructingParty,		Contractstatus		  
//      "A"						"init"				   
// ============================================================================================================================
func getCtrctStateByPingS(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	type Everything struct {
		remittance  []CIPS  `json:"remittance"`
	}
	var everything Everything
	
	
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	
	
	InstructingParty := args[0]
	Contractstatus := args[1]

	fmt.Println("- start getCtrctStateByPingS ", InstructingParty, Contractstatus)

	// Query the 'InstructingParty~Contractstatus~MessageIdentification' index by 'InstructingParty~Contractstatus'
	// This will execute a key range query on all keys starting with 'InstructingParty'

	//PingSResultsIterator, err := stub.GetStateByPartialCompositeKey("InstructingParty~ContractStatus~MessageIdentification", []string{InstructingParty,Contractstatus})
	PingSResultsIterator, err := stub.GetStateByPartialCompositeKey("InstructingParty~ContractStatus~MessageIdentification", []string{InstructingParty,Contractstatus})
	if err != nil {
		return shim.Error(err.Error())
	}
	defer PingSResultsIterator.Close()
	
	var i int
	for i = 0; PingSResultsIterator.HasNext(); i++ {
		tempPingSResultsIterator,temperr := PingSResultsIterator.Next()
		fmt.Println("- -------------PingSResultsIterator.Next()-------------: %+v\n%s\n", tempPingSResultsIterator,temperr)
		if temperr != nil {
			return shim.Error(err.Error())
		}
		queryValAsBytes := tempPingSResultsIterator.Value
		fmt.Println("- -------------tempPingSResultsIterator.Value-------------:\n", queryValAsBytes)
		var tmpremittance CIPS
		json.Unmarshal(queryValAsBytes, &tmpremittance)                  //un stringify it aka JSON.parse()
		everything.remittance = append(everything.remittance, tmpremittance)
	}
	
	fmt.Println("- ------------- remittance array --------- ", everything.remittance)
	
	
	//change to array of bytes
	everythingAsBytes, _ := json.Marshal(everything.remittance)              //convert to array of bytes
	return shim.Success(everythingAsBytes)
}


// ============================================================================================================================
//  Get a remittance state by  InstructedParty/Contractstatus
//
//  Inputs - string
//		1,						2					
//  InstructedParty,		Contractstatus		  
//      "B"						"init"				   
// ============================================================================================================================
func getCtrctStateByPedS(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	type Everything struct {
		remittance  []CIPS  `json:"remittance"`
	}
	var everything Everything
	
	
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	
	
	InstructedParty := args[0]
	Contractstatus := args[1]

	fmt.Println("- start getCtrctStateByPedS ", InstructedParty, Contractstatus)

	// Query the 'InstructedParty~Contractstatus~MessageIdentification' index by 'InstructedParty~Contractstatus'
	// This will execute a key range query on all keys starting with 'InstructedParty'

	PedSResultsIterator, err := stub.GetStateByPartialCompositeKey("InstructedParty~ContractStatus~MessageIdentification", []string{InstructedParty,Contractstatus})
	if err != nil {
		return shim.Error(err.Error())
	}
	defer PedSResultsIterator.Close()
	
	var i int
	for i = 0; PedSResultsIterator.HasNext(); i++ {
		tempPedSResultsIterator,temperr := PedSResultsIterator.Next()
		fmt.Println("- -------------PedSResultsIterator.Next()-------------: %+v\n%s\n", tempPedSResultsIterator,temperr)
		if temperr != nil {
			return shim.Error(err.Error())
		}
		queryValAsBytes := tempPedSResultsIterator.Value
		fmt.Println("- -------------tempPedSResultsIterator.Value-------------:\n", queryValAsBytes)
		var tmpremittance CIPS
		json.Unmarshal(queryValAsBytes, &tmpremittance)                  //un stringify it aka JSON.parse()
		everything.remittance = append(everything.remittance, tmpremittance)
	}
	
	fmt.Println("- ------------- remittance array --------- ", everything.remittance)
	
	
	//change to array of bytes
	everythingAsBytes, _ := json.Marshal(everything.remittance)              //convert to array of bytes
	return shim.Success(everythingAsBytes)
}	
	
	
//	============================================================================================================================
//  Get a remittance state by Contractstatus
//
//  Inputs - string
//		1					
//  Contractstatus		  
//      "init"				   
// ============================================================================================================================


func getCtrctStateByStatus(stub shim.ChaincodeStubInterface, args []string) pb.Response{
	type Everything struct {
		remittance  []CIPS  `json:"remittance"`
	}
	var everything Everything
	
	
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	
	
	Contractstatus := args[0]

	fmt.Println("- start getCtrctStateByStatus ",  Contractstatus)

	// Query the 'Contractstatus~MessageIdentification' index by 'Contractstatus'
	// This will execute a key range query on all keys starting with 'Contractstatus'

	SResultsIterator, err := stub.GetStateByPartialCompositeKey("ContractStatus~MessageIdentification", []string{Contractstatus})
	if err != nil {
		return shim.Error(err.Error())
	}
	defer SResultsIterator.Close()
	
	var i int
	for i = 0; SResultsIterator.HasNext(); i++ {
		tempSResultsIterator,temperr := SResultsIterator.Next()
		fmt.Println("- -------------SResultsIterator.Next()-------------: %+v\n%s\n", tempSResultsIterator,temperr)
		if temperr != nil {
			return shim.Error(err.Error())
		}
		queryValAsBytes := tempSResultsIterator.Value
		fmt.Println("- -------------tempSResultsIterator.Value-------------:\n", queryValAsBytes)
		var tmpremittance CIPS
		json.Unmarshal(queryValAsBytes, &tmpremittance)                  //un stringify it aka JSON.parse()
		everything.remittance = append(everything.remittance, tmpremittance)
	}
	
	fmt.Println("- ------------- remittance array --------- ", everything.remittance)
	
	
	//change to array of bytes
	everythingAsBytes, _ := json.Marshal(everything.remittance)              //convert to array of bytes
	return shim.Success(everythingAsBytes)
}