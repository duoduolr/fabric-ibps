package main

import (
	//"encoding/json"
	"fmt"
//	"strconv"
//	"strings"
	"crypto/aes"
	//"github.com/hyperledger/fabric/core/chaincode/shim"
	//pb "github.com/hyperledger/fabric/protos/peer"
	b "github.com/hyperledger/fabric/bccsp"
	"github.com/hyperledger/fabric/bccsp/sw"
	"math/big"
	"crypto/rsa"
	"bytes"
	"crypto/sha1"
)


type ibpsPublicKey struct {
	N *big.Int // modulus
	E int      // public exponent
}

func GenAESKey(len int) (key []byte) {

//-----------------------------test gen AES key start-----------
	//kg := &sw.AesKeyGenerator{Length: 32}
	kg := &sw.AesKeyGenerator{Length: len}
	k, err := kg.KeyGen(&b.AES256KeyGenOpts{Temporary: true})
	aesK, _ := k.(*sw.AesPrivateKey)
	if err != nil {
		//return fmt.Printf("AES256KeyGenForTest error: KeyGen returned %s\n", err)
	}
	
	fmt.Println("--------------------KeyGen returned------- \n", aesK.PrivKey)
	
	//var key []byte = []byte{236,208,100,126,108,4,52,191,165,51,65,176,1,106,77,105,16,213,210,37,153,206,210,34,72,152,175,77,144,100,64,62}
	key = aesK.PrivKey
	
	return key

}


//----------------------
func AESEncryptString(msgBefore string, key []byte) (string){
	
	//encrypt some text of contract
	var ptext = []byte(msgBefore)
	Mod :=len(ptext)%aes.BlockSize
	for i := 0; i < aes.BlockSize - Mod; i++ {
		ptext = append(ptext,0)
	}
	fmt.Println("---------------------aes.BlockSize is--: ", aes.BlockSize," \n")
	fmt.Println("---------------------Before encrypting --ptext.length is--: ", len(msgBefore)," \n")
	fmt.Println("---------------------Before encrypting --ptext is--: '%s'  ----and key is-- %v \n", ptext, key)
	msgAfterbyte, encmsgAfterErr := sw.AesCBCEncrypt(key, ptext)
	
	fmt.Println("---------------------After encrypting --encrypted msg is--: '%s' \n", msgAfterbyte)
	if encmsgAfterErr != nil {
		fmt.Println("Error encrypting '%s': %v", ptext, encmsgAfterErr)
	}
	
	return string(msgAfterbyte)
	//encrypt some text of contract end

}

//
func RSAEncryptWithPub(msgin []byte, pub rsa.PublicKey) (msgout string) {
	
	var seed []byte = []byte{236,208,100,126,108,4,208,100,126,108,4,52,208,100,126,108,4,52,106,77} //seed size must be size of hash method(if sha1,then size is 20)
	sha1 := sha1.New()
	public := pub	
	randomSource := bytes.NewReader(seed)
	
	fmt.Println("--------------------hash.Size------- \n", sha1.Size())
	fmt.Println("--------------------sha1------- \n",sha1)
	fmt.Println("--------------------public key------- \n", public)
	fmt.Println("--------------------randomSource------- \n", randomSource)
	msgEncrypted, err1 := rsa.EncryptOAEP(sha1, randomSource, &public, msgin, nil)
	if err1 != nil {
		fmt.Println("#%d,%d error1: %s", err1)
	}
	fmt.Println("--------------------message before encryptOAEP ------- \n",msgin)
	fmt.Println("--------------------message after encryptOAEP ------- \n", msgEncrypted)
	
	return string(msgEncrypted)

}


//---------------