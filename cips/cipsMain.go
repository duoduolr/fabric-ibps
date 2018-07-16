/*

//todo

*/

package main

import (
	"fmt"
	"strconv"
	"errors"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

// ============================================================================================================================
// Asset Definitions - The ledger will store CIPS
// ============================================================================================================================

// ----- CIPS ----- //
type CIPS struct {

	MessageIdentification       string        `json:"MessageIdentification"`  	//业务ID InstructingParty+YYYYMMDD+8numbers
	InstructingParty      	    string        `json:"InstructingParty"`        	//发起行ID
	DebtorAccount               string        `json:"DebtorAccount"`           	//汇款人账户
	DebtorName            	    string        `json:"DebtorName"`        		//汇款人姓名
	InstructedParty       	    string        `json:"InstructedParty"`     		//接受行ID
	CreditorAccount      	    string        `json:"CreditorAccount"`		//收款人账户
	CreditorName        	    string        `json:"CreditorName"`			//收款人姓名
	Amount                      string        `json:"Amount"`			//金额
	ContractStatus              string        `json:"ContractStatus"`	        //状态 init-confirmed-cleared-remitted-received
	InstructingKey              string        `json:"InstructingKey"`               //发起行加密后的对称加密KEY 
	InstructedKey               string        `json:"InstructedKey"`                //接受行加密后的对称加密KEY
        ClearHouseKey               string        `json:"ClearHouseKey"`                //清算中心加密后的对称加密KEY
}


// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode - %s", err)
	}
}


// ============================================================================================================================
// Init - initialize the chaincode 
//
// CIPS does not require initialization, so let's run a simple test instead.
//
// Shows off PutState() and how to pass an input argument to chaincode.
// Shows off GetFunctionAndParameters() and GetStringArgs()
// Shows off GetTxID() to get the transaction ID of the proposal
//
// Inputs - Array of strings
//  ["314"]
// 
// Returns - shim.Success or error
// ============================================================================================================================
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("CIPS Is Starting Up")
	funcName, args := stub.GetFunctionAndParameters()
	var number int
	var err error
	txId := stub.GetTxID()
	
	fmt.Println("Init() is running")
	fmt.Println("Transaction ID:", txId)
	fmt.Println("  GetFunctionAndParameters() function:", funcName)
	fmt.Println("  GetFunctionAndParameters() args count:", len(args))
	fmt.Println("  GetFunctionAndParameters() args found:", args)

	// expecting 1 arg for instantiate or upgrade
	if len(args) == 1 {
		fmt.Println("  GetFunctionAndParameters() arg[0] length", len(args[0]))

		// expecting arg[0] to be length 0 for upgrade
		if len(args[0]) == 0 {
			fmt.Println("  Uh oh, args[0] is empty...")
		} else {
			fmt.Println("  Great news everyone, args[0] is not empty")

			// convert numeric string to integer
			number, err = strconv.Atoi(args[0])
			if err != nil {
				return shim.Error("Expecting a numeric string argument to Init() for instantiate")
			}

			// this is a very simple test. let's write to the ledger and error out on any errors
			// it's handy to read this right away to verify network is healthy if it wrote the correct value
			err = stub.PutState("selftest", []byte(strconv.Itoa(number)))
			if err != nil {
				return shim.Error(err.Error())                  //self-test fail
			}
		}
	}

	// showing the alternative argument shim function
	alt := stub.GetStringArgs()
	fmt.Println("  GetStringArgs() args count:", len(alt))
	fmt.Println("  GetStringArgs() args found:", alt)

	// store compatible marbles application version
	err = stub.PutState("CIPS_ui", []byte("0.0.1"))
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("Ready for action")                          //self-test pass
	return shim.Success(nil)
}


// ============================================================================================================================
// Invoke - Our entry point for Invocations
// ============================================================================================================================
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println(" ")
	fmt.Println("starting invoke, for - " + function)

	// Handle different functions
	if function == "init" {                    //initialize the chaincode state, used as reset
		return t.Init(stub)
	} else if function == "Create_Remittance" {      //create a new Remittance
		return Create_Remittance(stub, args)
	} else if function == "modifyStatus"{        //modify a Remittance state 
		return modifyStatus(stub, args)
	} else if function == "getCtrctStateById"{        //get a Remittance state by ID
		return getCtrctStateById(stub, args)		
	} else if function == "getCtrctStateByPingS"{        //get a Remittance state by InstructingParty~ContractStatus
		return getCtrctStateByPingS(stub, args)
	} else if function == "getCtrctStateByPedS"{        //get a Remittance state by InstructedParty~ContractStatus
		return getCtrctStateByPedS(stub, args)
	} else if function == "getCtrctStateByStatus"{        //get a Remittance state by ContractStatus
		return getCtrctStateByStatus(stub, args)
	} else if function == "getHistory"{        //get a Remittance state by ContractStatus
		return getHistory(stub, args)
	} else if function == "read" {    //generic read ledger
		return read(stub, args)
	} else if function == "write" {      //generic writes to ledger
		return write(stub, args)
	}

	// error out
	fmt.Println("Received unknown invoke function name - " + function)
	return shim.Error("Received unknown invoke function name - '" + function + "'")
}

// ========================================================
// Input Sanitation - dumb input checking, look for empty strings
// ========================================================
func sanitize_arguments(strs []string) error{
	for i, val:= range strs {
		if len(val) <= 0 {
			return errors.New("Argument " + strconv.Itoa(i) + " must be a non-empty string")
		}
		if len(val) > 320 {
			return errors.New("Argument " + strconv.Itoa(i) + " must be <= 320 characters")
		}
	}
	return nil
}




