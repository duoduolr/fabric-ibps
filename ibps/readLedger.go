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
func getContractInfo(stub shim.ChaincodeStubInterface, id string) (Contract, error) {
	var contract Contract
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
		Value   Contract   `json:"value"`
	}
	var history []AuditHistory;
	var contract Contract

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
			var emptyContract Contract
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
// Get everything we need (owners + marbles + companies)
//
// Inputs - none
//
// ============================================================================================================================
func read_everything(stub shim.ChaincodeStubInterface) pb.Response {
	type Everything struct {
		Contracts  []Contract  `json:"contracts"`
	}
	var everything Everything

	// ---- Get All Contracts ---- //
	resultsIterator, err := stub.GetStateByRange("id0", "id9999999999999999999")
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()
	
	for resultsIterator.HasNext() {
		aKeyValue, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		queryKeyAsStr := aKeyValue.Key
		queryValAsBytes := aKeyValue.Value
		fmt.Println("on contract id - ", queryKeyAsStr)
		var contract Contract
		json.Unmarshal(queryValAsBytes, &contract)                  //un stringify it aka JSON.parse()
		everything.Contracts = append(everything.Contracts, contract)   //add this contract to the list
	}
	fmt.Println("contract array - ", everything.Contracts)

	//change to array of bytes
	everythingAsBytes, _ := json.Marshal(everything)              //convert to array of bytes
	return shim.Success(everythingAsBytes)
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
	
	//------------------------------------------test decrypted-------------------------------------------------------------------------------
	/*var key []byte = []byte{236,208,100,126,108,4,52,191,165,51,65,176,1,106,77,105,16,213,210,37,153,206,210,34,72,152,175,77,144,100,64,62}
	decrypted, decErr := sw.AesCBCDecrypt(key, contractAsBytes)
	fmt.Println("---------------------After decrypting --decrypted is--: '%s' \n", decrypted)
	if decErr != nil {
		fmt.Println("Error decrypting '%s': %v", contractAsBytes, decErr)
	}*/
	//------------------------------------------test decrypted-------------------------------------------------------------------------------
	
	return shim.Success(contractAsBytes)
}


// ============================================================================================================================
//  Get a Contract state by  InstructingParty/QueristAccount/QueristName
//
//  Inputs - string
//		1,						2,						3
//  InstructingParty,		QueristAccount,			QueristName    
//      "A"						"q0000"				   "Jack"
// ============================================================================================================================
func getCtrctStateByIPQAQN(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	type Everything struct {
		contracts  []Contract  `json:"contracts"`
	}
	var everything Everything
	
	
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}
	
	
	InstructingParty := args[0]
	QueristAccount := args[1]
	QueristName := args[2]
	fmt.Println("- start getCtrctStateByIPQAQN ", InstructingParty, QueristAccount,QueristName)

	// Query the 'InstructingParty~QueristAccount~QueristName~MessageIdentification' index by 'InstructingParty~QueristAccount~QueristName'
	// This will execute a key range query on all keys starting with 'InstructingParty'

	IPQAResultsIterator, err := stub.GetStateByPartialCompositeKey("InstructingParty~QueristAccount~QueristName~MessageIdentification", []string{InstructingParty,QueristAccount,QueristName})
	if err != nil {
		return shim.Error(err.Error())
	}
	defer IPQAResultsIterator.Close()
	
	var i int
	for i = 0; IPQAResultsIterator.HasNext(); i++ {
		tempIPQAResultsIterator,temperr := IPQAResultsIterator.Next()
		fmt.Println("- -------------IPQAResultsIterator.Next()-------------: %+v\n%s\n", tempIPQAResultsIterator,temperr)
		if temperr != nil {
			return shim.Error(err.Error())
		}
		queryValAsBytes := tempIPQAResultsIterator.Value
		fmt.Println("- -------------tempIPQAResultsIterator.Value-------------:\n", queryValAsBytes)
		var tmpcontract Contract
		json.Unmarshal(queryValAsBytes, &tmpcontract)                  //un stringify it aka JSON.parse()
		everything.contracts = append(everything.contracts, tmpcontract)
	}
	
	fmt.Println("- ---------------contract array --------- ", everything.contracts)
	
	
	//change to array of bytes
	everythingAsBytes, _ := json.Marshal(everything.contracts)              //convert to array of bytes
	return shim.Success(everythingAsBytes)
}



// ============================================================================================================================
//  Get a Contract state by InstructedParty/ReplierAccount/ReplierName
//
//  Inputs - string
//		1,						2,						3
//  InstructedParty,		ReplierAccount,			ReplierName    
//      "B"						"r0000"				   "Bob"
// ============================================================================================================================
func getCtrctStateByIPRARN(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	type Everything struct {
		contracts  []Contract  `json:"contracts"`
	}
	var everything Everything
	
	
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}
	
	
	InstructedParty := args[0]
	ReplierAccount := args[1]
	ReplierName := args[2]
	fmt.Println("- start getCtrctStateByIPRARN ", InstructedParty,ReplierAccount,ReplierName)

	// Query the 'InstructedParty~ReplierAccount~ReplierName~MessageIdentification' index by 'InstructedParty~ReplierAccount~ReplierName'
	// This will execute a key range query on all keys starting with 'InstructedParty'

	IPRARNResultsIterator, err := stub.GetStateByPartialCompositeKey("InstructedParty~ReplierAccount~ReplierName~MessageIdentification", []string{InstructedParty,ReplierAccount,ReplierName})
	if err != nil {
		return shim.Error(err.Error())
	}
	defer IPRARNResultsIterator.Close()
	
	var i int
	for i = 0; IPRARNResultsIterator.HasNext(); i++ {
		tempIPRARNResultsIterator,temperr := IPRARNResultsIterator.Next()
		fmt.Println("- -------------IPRARNResultsIterator.Next()-------------: %+v\n%s\n", tempIPRARNResultsIterator,temperr)
		if temperr != nil {
			return shim.Error(err.Error())
		}
		queryValAsBytes := tempIPRARNResultsIterator.Value
		fmt.Println("- -------------IPRARNResultsIterator.Value-------------:\n", queryValAsBytes)
		var tmpcontract Contract
		json.Unmarshal(queryValAsBytes, &tmpcontract)                  //un stringify it aka JSON.parse()
		everything.contracts = append(everything.contracts, tmpcontract)
	}
	
	fmt.Println("- ---------------contract array --------- ", everything.contracts)
	
	
	//change to array of bytes
	everythingAsBytes, _ := json.Marshal(everything.contracts)              //convert to array of bytes
	return shim.Success(everythingAsBytes)
}


// ============================================================================================================================
// Read - read a generic variable from ledger
//
// Shows Off GetState() - reading a key/value from the ledger
//
// Inputs - Array of strings
//  0
//  key
//  "abc"
// 
// Returns - string
// ============================================================================================================================
func read(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var key, jsonResp string
	var err error
	fmt.Println("starting read")

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting key of the var to query")
	}

	// input sanitation
	err = sanitize_arguments(args)
	if err != nil {
		return shim.Error(err.Error())
	}

	key = args[0]
	valAsbytes, err := stub.GetState(key)           //get the var from ledger
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return shim.Error(jsonResp)
	}

	fmt.Println("- end read")
	return shim.Success(valAsbytes)                  //send it onward
}
