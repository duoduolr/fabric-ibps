/*

//todo

*/

package main

import (
	"encoding/json"
	"fmt"
//	"strconv"
//	"strings"
	//"crypto/aes"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	//b "github.com/hyperledger/fabric/bccsp"
	//"github.com/hyperledger/fabric/bccsp/sw"
	//"math/big"
	"crypto/rsa"
	"crypto/rand"
)

// ============================================================================================================================
// Create Contracts - create a new contracts, store into chaincode state
//
// Shows off building a key's JSON value manually
//
// Inputs - Array of strings
//      	0,   						 1, 			  2,      			 3,				     4,				 5,					 6,					  7
//   MessageIdentification,		InstructingParty,	InstructedParty,	QueristAccount,		QueristName,	ReplierAccount, 	ReplierName,		ContractStatus	
//   "id0000", 						"A",                 "B", 				"q0000", 		  "Jack",			"r0000",			"Bob",    			"init"
// ============================================================================================================================
func Create_Contract(stub shim.ChaincodeStubInterface, args []string) (pb.Response) {
	var err error
	//var bccspInst b.BCCSP
	fmt.Println("starting Create_Contract")

	if len(args) != 8 {
		return shim.Error("Incorrect number of arguments. Expecting 8")
	}

	//input sanitation
	err = sanitize_arguments(args)
	if err != nil {
		return shim.Error(err.Error())
	}
	
	/*
	//-----------------------------test gen AES key start-----------
	//kg := &sw.AesKeyGenerator{Length: 32}
	kg := &sw.AesKeyGenerator{Length: 16}
	k, err := kg.KeyGen(&b.AES256KeyGenOpts{Temporary: true})
	aesK, _ := k.(*sw.AesPrivateKey)
	if err != nil {
		//return fmt.Printf("AES256KeyGenForTest error: KeyGen returned %s\n", err)
	}
	
	fmt.Println("--------------------KeyGen returned------- \n", aesK.PrivKey)
	
	//var key []byte = []byte{236,208,100,126,108,4,52,191,165,51,65,176,1,106,77,105,16,213,210,37,153,206,210,34,72,152,175,77,144,100,64,62}
	key := aesK.PrivKey
	//------------------------------test gen AES key end-----------
    */

	MessageIdentification := args[0]
	InstructingParty := args[1]
	InstructedParty := args[2]
	QueristAccount := args[3]
	QueristName := args[4]
	ReplierAccount := args[5]
	ReplierName := args[6]
	ContractStatus := args[7]
	RejectReason := ""

	//todo - check if contract id already exists


    /*
	//encrypt some text of contract
	var ptext = []byte(QueristAccount)
	Mod :=len(ptext)%aes.BlockSize
	for i := 0; i < aes.BlockSize - Mod; i++ {
		ptext = append(ptext,0)
	}
	fmt.Println("---------------------aes.BlockSize is--: ", aes.BlockSize," \n")
	fmt.Println("---------------------Before encrypting --ptext.length is--: ", len(QueristAccount)," \n")
	fmt.Println("---------------------Before encrypting --ptext is--: '%s'  ----and key is-- %v \n", ptext, key)
	encryptedQueristAccount, encQueristAccountErr := sw.AesCBCEncrypt(key, ptext)
	
	fmt.Println("---------------------After encrypting --encrypted QueristAccount is--: '%s' \n", encryptedQueristAccount)
	if encQueristAccountErr != nil {
		fmt.Println("Error encrypting '%s': %v", ptext, encQueristAccountErr)
	}
	//encrypt some text of contract end
    */
	AESkey := GenAESKey(16)

	encryptedQueristAccount := AESEncryptString(QueristAccount,AESkey)
	encryptedQueristName := AESEncryptString(QueristName,AESkey)
	encryptedReplierAccount := AESEncryptString(ReplierAccount,AESkey)
	encryptedReplierName := AESEncryptString(ReplierName,AESkey)

	//RSApub :=new(rsa.PublicKey)
	//pubN := big.NewInt(int64(9017))
	//RSApub.N = pubN
	//RSApub.E = 65537
	
	lowLevelKey, err := rsa.GenerateKey(rand.Reader, 512)

	if err != nil {
		
	}
	
	fmt.Println("--------------------lowLevelKey.PublicKey.N------- v%\n", lowLevelKey.PublicKey.N)
	fmt.Println("--------------------lowLevelKey.PublicKey.E------- v%\n", lowLevelKey.PublicKey.E)
	fmt.Println("--------------------lowLevelKey.PrivateKey------- v%\n", lowLevelKey.D)
	
	
	encryptedRSApub := RSAEncryptWithPub(AESkey,lowLevelKey.PublicKey)
	fmt.Println("---------------------encryptedRSApub is--: ", encryptedRSApub," \n")
	//build the contract json string manually
	str := `{
		"MessageIdentification": "` + MessageIdentification + `", 
		"InstructingParty": "` + InstructingParty + `", 
		"InstructedParty": "` + InstructedParty + `", 
		"QueristAccount": "` + encryptedQueristAccount + `",
		"QueristName": "` + encryptedQueristName + `",
		"ReplierAccount": "` + encryptedReplierAccount + `",
		"ReplierName": "` + encryptedReplierName +`",
		"ContractStatus": "` + ContractStatus +`",
		"RejectReason": "` + RejectReason +`"
	}`
	
	//------------------------------test Encrypt with AES key start-----------
	/*var ptext = []byte(str)
	fmt.Println("---------------------Before encrypting --ptext.length is--: ", len(str)," \n")
	//fmt.Println("---------------------Before encrypting --ptext is--: '%s'  ----and key is-- %v \n", ptext, aesK.PrivKey)
	fmt.Println("---------------------Before encrypting --ptext is--: '%s'  ----and key is-- %v \n", ptext, key)
	//encrypted, encErr := sw.AesCBCEncrypt(aesK.PrivKey, ptext)
	encrypted, encErr := sw.AesCBCEncrypt(key, ptext)
	
	fmt.Println("---------------------After encrypting --encrypted is--: '%s' \n", encrypted)
	if encErr != nil {
		fmt.Println("Error encrypting '%s': %v", ptext, encErr)
	}*/
	//------------------------------test Encrypt with AES key end-----------
	
	//err = stub.PutState(MessageIdentification, encrypted)
	err = stub.PutState(MessageIdentification, []byte(str))                         //store marble with MessageIdentification as key
	if err != nil {
		return shim.Error(err.Error())
	}


    // CompositeKey: "InstructingParty~QueristAccount~QueristName~MessageIdentification"
	//   
	indexName := "InstructingParty~QueristAccount~QueristName~MessageIdentification"
	IPQAQNIndexKey, err := stub.CreateCompositeKey(indexName, []string{InstructingParty,QueristAccount,QueristName,MessageIdentification})
	fmt.Println("-------------IPQAQNIndexKey-------------: %s\n",IPQAQNIndexKey)
	if err != nil {
		return shim.Error(err.Error())
	}

	var err1 error
	err1 = stub.PutState(IPQAQNIndexKey,[]byte(str))
	if err1 != nil {
		return shim.Error(err1.Error())
	}

    // CompositeKey: "InstructedParty~ReplierAccount~ReplierName~MessageIdentification"
	//   
	indexName2 := "InstructedParty~ReplierAccount~ReplierName~MessageIdentification"
	IPRARNIndexKey, err := stub.CreateCompositeKey(indexName2, []string{InstructedParty,ReplierAccount,ReplierName,MessageIdentification})
	fmt.Println("-------------IPRARNIndexKey-------------: %s\n",IPRARNIndexKey)
	if err != nil {
		return shim.Error(err.Error())
	}

	var err2 error
	err2 = stub.PutState(IPRARNIndexKey,[]byte(str))
	if err2 != nil {
		return shim.Error(err1.Error())
	}


	fmt.Println("- end Create_Contract")
	return shim.Success(nil)
}

// ============================================================================================================================
// Create Contracts - create a new contracts, store into chaincode state
//
// Shows off building a key's JSON value manually
//
// Inputs - Array of strings
//      	0,   						 1, 			  2,      			 3,				     4,				 5,					 6,					  7
//   MessageIdentification,		InstructingParty,	InstructedParty,	QueristAccount,		QueristName,	ReplierAccount, 	ReplierName,		ContractStatus	
//   "id0000", 						"A",                 "B", 				"q0000", 		  "Jack",			"r0000",			"Bob",    			"init"
// ============================================================================================================================
/*func Create_ContractForManage(stub shim.ChaincodeStubInterface, args []string) (pb.Response) {
	var err error
	fmt.Println("starting Create_Contract")

	if len(args) != 8 {
		return shim.Error("Incorrect number of arguments. Expecting 8")
	}

	//input sanitation
	err = sanitize_arguments(args)
	if err != nil {
		return shim.Error(err.Error())
	}

	MessageIdentification := args[0]
	InstructingParty := args[1]
	InstructedParty := args[2]
	QueristAccount := args[3]
	QueristName := args[4]
	ReplierAccount := args[5]
	ReplierName := args[6]
	ContractStatus := args[7]
	RejectReason := ""

	//todo - check if contract id already exists


	//build the contract json string manually
	str := `{
		"MessageIdentification": "` + MessageIdentification + `", 
		"InstructingParty": "` + InstructingParty + `", 
		"InstructedParty": "` + InstructedParty + `", 
		"QueristAccount": "` + QueristAccount + `",
		"QueristName": "` + QueristName + `",
		"ReplierAccount": "` + ReplierAccount + `",
		"ReplierName": "` + ReplierName +`",
		"ContractStatus": "` + ContractStatus +`",
		"RejectReason": "` + RejectReason +`"
	}`
	err = stub.PutState(MessageIdentification, []byte(str))                         //store marble with MessageIdentification as key
	if err != nil {
		return shim.Error(err.Error())
	}


    // CompositeKey: "InstructingParty~QueristAccount~QueristName~MessageIdentification"
	//   
	indexName := "InstructingParty~QueristAccount~QueristName~MessageIdentification"
	IPQAQNIndexKey, err := stub.CreateCompositeKey(indexName, []string{InstructingParty,QueristAccount,QueristName,MessageIdentification})
	fmt.Println("-------------IPQAQNIndexKey-------------: %s\n",IPQAQNIndexKey)
	if err != nil {
		return shim.Error(err.Error())
	}

	var err1 error
	err1 = stub.PutState(IPQAQNIndexKey,[]byte(str))
	if err1 != nil {
		return shim.Error(err1.Error())
	}

    // CompositeKey: "InstructedParty~ReplierAccount~ReplierName~MessageIdentification"
	//   
	indexName2 := "InstructedParty~ReplierAccount~ReplierName~MessageIdentification"
	IPRARNIndexKey, err := stub.CreateCompositeKey(indexName2, []string{InstructedParty,ReplierAccount,ReplierName,MessageIdentification})
	fmt.Println("-------------IPRARNIndexKey-------------: %s\n",IPRARNIndexKey)
	if err != nil {
		return shim.Error(err.Error())
	}

	var err2 error
	err2 = stub.PutState(IPRARNIndexKey,[]byte(str))
	if err2 != nil {
		return shim.Error(err1.Error())
	}


	fmt.Println("- end Create_Contract")
	return shim.Success(nil)
}*/
// ============================================================================================================================
// deleteContract() - remove a contract from state and from contract index
// 
// Shows Off DelState() - "removing"" a key/value from the ledger
//
// Inputs - Array of strings			
//             0,				   			 1, 								2	
//     MessageIdentification, 	 		QueristAccount,					ReplierAccount	
//      	"id0000",						"q0000",						"r0000"
// ============================================================================================================================
func deleteContract(stub shim.ChaincodeStubInterface, args []string) (pb.Response) {
	fmt.Println("starting deleteContract")

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	// input sanitation
	err := sanitize_arguments(args)
	if err != nil {
		return shim.Error(err.Error())
	}

	id := args[0]
	queristAccount := args[1]
	replierAccount := args[2]
	
	// get the contract
	contract, err := getContractInfo(stub, id)
	if err != nil{
		fmt.Println("Failed to find contract by MessageIdentification " + id)
		return shim.Error(err.Error())
	}


    // check querist account
	if contract.QueristAccount != queristAccount{
		return shim.Error("The querist account '" + queristAccount + "' is not available, please type the right querist account and try again! '" )
	}
	
	// check replier account
	if contract.ReplierAccount != replierAccount{
		return shim.Error("The replier account '" + replierAccount + "' is not available, please type the right replier account and try again! '" )
	}
	
	// remove the contract
	err = stub.DelState(id)                                                 //remove the key from chaincode state
	if err != nil {
		return shim.Error("Failed to delete state")
	}

	fmt.Println("- end deleteContract")
	return shim.Success(nil)
	
}	

// ============================================================================================================================
// write() - genric write variable into ledger
// 
// Shows Off PutState() - writting a key/value into the ledger
//
// Inputs - Array of strings
//    0   ,    1
//   key  ,  value
//  "abc" , "test"
// ============================================================================================================================
func write(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var key, value string
	var err error
	fmt.Println("starting write")

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2. key of the variable and value to set")
	}

	// input sanitation
	err = sanitize_arguments(args)
	if err != nil {
		return shim.Error(err.Error())
	}

	key = args[0]                                   //rename for funsies
	value = args[1]
	err = stub.PutState(key, []byte(value))         //write the variable into the ledger
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end write")
	return shim.Success(nil)
}


// ============================================================================================================================
// modifyStatus() - modify status of contract
// 
// Inputs - Array of strings
//  0,		  							1      				
// 	MessageIdentification  , 	 ContractStatus
// 	"id0000" , 						 "confirm"
// ============================================================================================================================
func modifyStatus(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var jsonResp string
	var err error
	var cID string
	var cStatus string
	fmt.Println("starting modifyStatus")

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
	
	valAsbytes, err := stub.GetState(cID)           //get the var from ledger
	
	
	
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + cID + "\"}"
		return shim.Error(jsonResp)
	}
	
	var tempVal Contract
	json.Unmarshal(valAsbytes, &tempVal)
	
	fmt.Println("--------tempVal_before---------:\n", tempVal)
	
	tempVal.ContractStatus = cStatus
	
	fmt.Println("--------tempVal_after---------:\n", tempVal)
	
	
	
	//change to array of bytes
	tempvalAsbytes, _ := json.Marshal(tempVal)              //convert to array of bytes
	
	
	err = stub.PutState(cID, tempvalAsbytes)                         //store marble with MessageIdentification as key
	if err != nil {
		return shim.Error(err.Error())
	}
	
	
	// CompositeKey: "InstructingParty~QueristAccount~QueristName~MessageIdentification"
	//   
	indexName := "InstructingParty~QueristAccount~QueristName~MessageIdentification"
	IPQAQNIndexKey, err := stub.CreateCompositeKey(indexName, []string{tempVal.InstructingParty,tempVal.QueristAccount,tempVal.QueristName,tempVal.MessageIdentification})
	fmt.Println("-------------IPQAQNIndexKey-------------: %s\n",IPQAQNIndexKey)
	if err != nil {
		return shim.Error(err.Error())
	}

	var err1 error
	err1 = stub.PutState(IPQAQNIndexKey,tempvalAsbytes)
	if err1 != nil {
		return shim.Error(err1.Error())
	}

    // CompositeKey: "InstructedParty~ReplierAccount~ReplierName~MessageIdentification"
	//   
	indexName2 := "InstructedParty~ReplierAccount~ReplierName~MessageIdentification"
	IPRARNIndexKey, err := stub.CreateCompositeKey(indexName2, []string{tempVal.InstructedParty,tempVal.ReplierAccount,tempVal.ReplierName,tempVal.MessageIdentification})
	fmt.Println("-------------IPRARNIndexKey-------------: %s\n",IPRARNIndexKey)
	if err != nil {
		return shim.Error(err.Error())
	}

	var err2 error
	err2 = stub.PutState(IPRARNIndexKey,tempvalAsbytes)
	if err2 != nil {
		return shim.Error(err1.Error())
	}
	
	fmt.Println("- end modifyStatus")
	return shim.Success(tempvalAsbytes)                  //send it onward
}

	