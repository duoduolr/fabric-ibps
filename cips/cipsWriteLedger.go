/*

//todo

*/

package main

import (
	"encoding/json"
	"fmt"
//	"strconv"
//	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// ============================================================================================================================
// Create Remittance - create a new remittance, store into chaincode state
//
// Shows off building a key's JSON value manually
//
// Inputs - Array of strings
//      	0,  	    		1, 	              2,      		   3,       	4,		         5,		  		 6,			  7          		8
//   MessageIdentification, InstructingParty, DebtorAccount,    DebtorName, InstructedParty, CreditorAccount, CreditorName, 	Amount	, Contractstatus
//   "id0001",   				 "A",             "q0000", 		 "Jack",     	 "B",             "r0001",        "Rose", 		"10000" , 	"init"
// ============================================================================================================================
func Create_Remittance(stub shim.ChaincodeStubInterface, args []string) (pb.Response) {
	var err error
	fmt.Println("starting Create_Contract")

	if len(args) != 9 {
		return shim.Error("Incorrect number of arguments. Expecting 9")
	}

	//input sanitation
	err = sanitize_arguments(args)
	if err != nil {
		return shim.Error(err.Error())
	}

	MessageIdentification := args[0]
	InstructingParty := args[1]
	DebtorAccount := args[2]
	DebtorName := args[3]
	InstructedParty := args[4]
	CreditorAccount := args[5]
	CreditorName := args[6]
	Amount := args[7]
	ContractStatus := args[8]
	
	//todo - check if contract id already exists
	contract, err := getContractInfo(stub, MessageIdentification)
	if err == nil {
		fmt.Println("This contract already exists - " + MessageIdentification)
		fmt.Println(contract)
		return shim.Error("This contract already exists - " + MessageIdentification)  // stop creating contract by this id exists
	}

	//build the contract json string manually
	str := `{
		"MessageIdentification": "` + MessageIdentification + `", 
		"InstructingParty": "` + InstructingParty + `", 
		"DebtorAccount": "` + DebtorAccount + `", 
		"DebtorName": "` + DebtorName + `",
		"InstructedParty": "` + InstructedParty + `",
		"CreditorAccount": "` + CreditorAccount + `",
		"CreditorName": "` + CreditorName +`",
		"Amount": "` + Amount +`",
		"ContractStatus": "` + ContractStatus +`"
	}`
	
	
	//store CIPS with MessageIdentification as key
	err = stub.PutState(MessageIdentification, []byte(str))                         //store CIPS with MessageIdentification as key
	if err != nil {
		return shim.Error(err.Error())
	}


    // CompositeKey: "InstructingParty~ContractStatus~MessageIdentification"
	// store CIPS with CompositeKey as key
	indexName := "InstructingParty~ContractStatus~MessageIdentification"
	PingSIDIndexKey, err1 := stub.CreateCompositeKey(indexName, []string{InstructingParty,ContractStatus,MessageIdentification})
	fmt.Println("-------------PingSIDIndexKey-------------: %s\n",PingSIDIndexKey)
	if err1 != nil {
		return shim.Error(err1.Error())
	}

	var err2 error
	err2 = stub.PutState(PingSIDIndexKey, []byte(str)) 
	if err2 != nil {
		return shim.Error(err2.Error())
	}
	
	
    // CompositeKey: "InstructedParty~ContractStatus~MessageIdentification"
	// store CIPS with CompositeKey as key  
	indexName2 := "InstructedParty~ContractStatus~MessageIdentification"
	PedSIDIndexKey, err3 := stub.CreateCompositeKey(indexName2, []string{InstructedParty,ContractStatus,MessageIdentification})
	fmt.Println("-------------PedSIDIndexKey-------------: %s\n",PedSIDIndexKey)
	if err3 != nil {
		return shim.Error(err.Error())
	}

	var err4 error
	err4 = stub.PutState(PedSIDIndexKey,[]byte(str))
	if err4 != nil {
		return shim.Error(err4.Error())
	}
	
	// 
	
	// CompositeKey: "ContractStatus~MessageIdentification"
	// store CIPS with CompositeKey as key  
	indexName3 := "ContractStatus~MessageIdentification"
	SIDIndexKey, err5 := stub.CreateCompositeKey(indexName3, []string{ContractStatus,MessageIdentification})
	fmt.Println("-------------SIDIndexKey-------------: %s\n",SIDIndexKey)
	if err5 != nil {
		return shim.Error(err5.Error())
	}
	
	var err6 error
	err6 = stub.PutState(SIDIndexKey,[]byte(str))
	if err6 != nil {
		return shim.Error(err6.Error())
	}


	fmt.Println("- end Create_Contract")
	return shim.Success(nil)
}


// ============================================================================================================================
// modifyStatus() - modify status of contract
// 
// Inputs - Array of strings
//  0,		  							1      				
// 	MessageIdentification  , 	 ContractStatus
// 	"id0000" , 						 "confirmed	"
// ============================================================================================================================
func modifyStatus(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var jsonResp string
	var err error
	var cID string
	var cStatus string
	fmt.Println("starting modifyStatus")
	
	indexName := "InstructingParty~ContractStatus~MessageIdentification"
	indexName2 := "InstructedParty~ContractStatus~MessageIdentification"
	indexName3 := "ContractStatus~MessageIdentification"

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2. key of the variable and value to set")
	}

	// input sanitation
	err = sanitize_arguments(args)
	if err != nil {
		return shim.Error(err.Error())
	}

	cID = args[0]
	cStatus = args[1]
	
	//todo - check if contract id exists
	contract, err := getContractInfo(stub, cID)
	if err != nil {
		fmt.Println("This contract doesn't exist - " + cID)
		fmt.Println(contract)
		return shim.Error("This contract doesn't exist - " + cID)  // stop creating contract by this id exists
	}
	
	valAsbytes, err := stub.GetState(cID)           //get the var from ledger
	
	
	
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + cID + "\"}"
		return shim.Error(jsonResp)
	}
	
	var tempVal CIPS
	json.Unmarshal(valAsbytes, &tempVal)
	
	fmt.Println("--------tempVal_before---------:\n", tempVal)
	
	
	//check if the input status valid
	
	if cStatus == "confirmed" && tempVal.ContractStatus != "init"{
		return shim.Error("You can't modify the status!!!! The current status is NOT 'init'!!!")
	}else if cStatus == "cleared" && tempVal.ContractStatus != "confirmed"{
		return shim.Error("You can't modify the status!!!! The current status is NOT 'confirmed'!!!")
	}else if cStatus == "remitted" && tempVal.ContractStatus != "cleared"{
		return shim.Error("You can't modify the status!!!! The current status is NOT 'cleared'!!!")
	}else if cStatus == "received" && tempVal.ContractStatus != "remitted"{
		return shim.Error("You can't modify the status!!!! The current status is NOT 'remitted'!!!")
	}else if cStatus == "init"{
		return shim.Error("You can't modify the status to init!!!!")
	}else if tempVal.ContractStatus == cStatus{
		return shim.Error("The input status and the current status are same!!!! Please check your input status!!!")
	}else {
		PingSIDIndexKey, err := stub.CreateCompositeKey(indexName, []string{tempVal.InstructingParty,tempVal.ContractStatus,tempVal.MessageIdentification})
		fmt.Println("-------------PingSIDIndexKey-------------: %s\n",PingSIDIndexKey)
		if err != nil {
			return shim.Error(err.Error())
		}
		
		// remove the Remittance that key is "InstructingParty~ContractStatus~MessageIdentification", because  Contract status will be modified
		delerr1 := stub.DelState(PingSIDIndexKey)                                                 //remove the key from chaincode state
		if delerr1 != nil {
			return shim.Error("Failed to delete Remittance that key is InstructingParty~ContractStatus~MessageIdentification")
		}
		
		PedSIDIndexKey, err := stub.CreateCompositeKey(indexName2, []string{tempVal.InstructedParty,tempVal.ContractStatus,tempVal.MessageIdentification})
		fmt.Println("-------------PedSIDIndexKey-------------: %s\n",PedSIDIndexKey)
		if err != nil {
			return shim.Error(err.Error())
		}
	
		// remove the Remittance that key is "InstructedParty~ContractStatus~MessageIdentification", because  Contract status will be modified
		delerr2 := stub.DelState(PedSIDIndexKey)                                                 //remove the key from chaincode state
		if delerr2 != nil {
			return shim.Error("Failed to delete Remittance that key is InstructedParty~ContractStatus~MessageIdentification")
		}
		
	
		SIDIndexKey, err := stub.CreateCompositeKey(indexName3, []string{tempVal.ContractStatus,tempVal.MessageIdentification})
		fmt.Println("-------------SIDIndexKey-------------: %s\n",SIDIndexKey)
		if err != nil {
			return shim.Error(err.Error())
		}
				
		
		// remove the Remittance that key is "ContractStatus~MessageIdentification", because  Contract status will be modified
		delerr3 := stub.DelState(SIDIndexKey)                                                 //remove the key from chaincode state
		if delerr3 != nil {
			return shim.Error("Failed to delete Remittance that key is ContractStatus~MessageIdentification")
		}
	
		tempVal.ContractStatus = cStatus
		fmt.Println("--------tempVal_after---------:\n", tempVal)
	}
	
	
	
	//change to array of bytes
	tempvalAsbytes, _ := json.Marshal(tempVal)              //convert to array of bytes
	
	
	err = stub.PutState(cID, tempvalAsbytes)                         //store marble with MessageIdentification as key
	if err != nil {
		return shim.Error(err.Error())
	}
	
	
	// CompositeKey: "InstructingParty~ContractStatus~MessageIdentification"
	// store CIPS with CompositeKey as key
	PingSIDIndexKey, err := stub.CreateCompositeKey(indexName, []string{tempVal.InstructingParty,tempVal.ContractStatus,tempVal.MessageIdentification})
	fmt.Println("-----modified--------PingSIDIndexKey-------------: %s\n",PingSIDIndexKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	
	

	var err1 error
	err1 = stub.PutState(PingSIDIndexKey,tempvalAsbytes)
	if err1 != nil {
		return shim.Error(err1.Error())
	}

    // CompositeKey: "InstructedParty~ContractStatus~MessageIdentification"
	// store CIPS with CompositeKey as key  
	PedSIDIndexKey, err := stub.CreateCompositeKey(indexName2, []string{tempVal.InstructedParty,tempVal.ContractStatus,tempVal.MessageIdentification})
	fmt.Println("-----modified--------PedSIDIndexKey-------------: %s\n",PedSIDIndexKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	var err2 error
	err2 = stub.PutState(PedSIDIndexKey,tempvalAsbytes)
	if err2 != nil {
		return shim.Error(err2.Error())
	}
	
	// CompositeKey: "ContractStatus~MessageIdentification"
	// store CIPS with CompositeKey as key  
	
	SIDIndexKey, err := stub.CreateCompositeKey(indexName3, []string{tempVal.ContractStatus,tempVal.MessageIdentification})
	fmt.Println("-----modified--------SIDIndexKey-------------: %s\n",SIDIndexKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	
	var err3 error
	err3 = stub.PutState(SIDIndexKey,tempvalAsbytes)
	if err3 != nil {
		return shim.Error(err3.Error())
	}
	
	//todo  delete old key-value "ContractStatus~MessageIdentification" 
	
	
	fmt.Println("- end modifyStatus")
	return shim.Success(tempvalAsbytes)                  //send it onward
}

	
