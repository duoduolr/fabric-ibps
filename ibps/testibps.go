package main

import (
	"fmt"
	//"strconv"
	//"errors"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	
	//"sync"

	b "github.com/hyperledger/fabric/bccsp"
	//"github.com/hyperledger/fabric/bccsp/factory"
	 "github.com/hyperledger/fabric/bccsp/sw"
	 
	// "crypto/ecdsa"
	//"crypto/elliptic"
	"crypto/rand"
	//"crypto/sha256"
	//"crypto/x509"
	//"math/big"
	"crypto/rsa"
	"bytes"
	//"crypto"
	"crypto/sha1"

	//"github.com/hyperledger/fabric/bccsp/utils"
	 
)


func TestRSAKey(stub shim.ChaincodeStubInterface) (pb.Response){
	
	fmt.Println("--------------------TestRSAKey start------- \n")
	
	//lowLevelKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	lowLevelKey, err := rsa.GenerateKey(rand.Reader, 512)

	if err != nil {
		
	}
	
	fmt.Println("--------------------lowLevelKey.PublicKey.N------- v%\n", lowLevelKey.PublicKey.N)
	fmt.Println("--------------------lowLevelKey.PublicKey.E------- v%\n", lowLevelKey.PublicKey.E)
	fmt.Println("--------------------lowLevelKey.PrivateKey------- v%\n", lowLevelKey.D)
	fmt.Println("--------------------lowLevelKey.Primes------- v%\n", lowLevelKey.Primes)
	fmt.Println("--------------------lowLevelKey.Precomputed------- v%\n", lowLevelKey.Precomputed)

//************************golang encrypted&decrypted method test************************
/*	
	pub := &lowLevelKey.PublicKey		//public key
	m := big.NewInt(42)     //message before encrypted
	c := rsa.Encrypt(new(big.Int), pub, m)		//message after encrypted
	fmt.Println("--------------------Before encrypted m------- v%\n", m)
	fmt.Println("--------------------After encrypted c------- v%\n", c)

	m2, err1 := rsa.Decrypt(nil, lowLevelKey, c)  //message decrypt
	
	fmt.Println("--------------------After decrypt m2------- v%\n", m2)
	if err1 != nil {
		fmt.Println("error while decrypting: %s\n", err1)
		
	}
	if m.Cmp(m2) != 0 {
		fmt.Println("got:%v, want:%v (%+v)\n", m2, m, lowLevelKey)
	}
*/
//************************golang EncryptOAEP & DecryptOAEP method test************************
	//var key []byte = []byte{236,208,100,126,108,4,52,191,165,51,65,176,1,106,77,105,16,213,210,37,153,206,210,34,72,152,175,77,144,100,64,62}
	var key []byte = []byte{236,208,100,126,108,4,52,191,165,51,65,176,1,106,77,105}
	var seed []byte = []byte{236,208,100,126,108,4,208,100,126,108,4,52,208,100,126,108,4,52,106,77} //seed size must be size of hash method(if sha1,then size is 20)
	sha1 := sha1.New()
	public := lowLevelKey.PublicKey	
	randomSource := bytes.NewReader(seed)
	
	fmt.Println("--------------------hash.Size------- \n", sha1.Size())
	fmt.Println("--------------------sha1------- \n",sha1)
	fmt.Println("--------------------public key------- \n", public)
	fmt.Println("--------------------randomSource------- \n", randomSource)
	msgEncrypted, err1 := rsa.EncryptOAEP(sha1, randomSource, &public, key, nil)
	if err1 != nil {
		fmt.Println("#%d,%d error1: %s", err1)
	}
	fmt.Println("--------------------message before encryptOAEP ------- \n",key)
	fmt.Println("--------------------message after encryptOAEP ------- \n", msgEncrypted)		
	
	
	private := lowLevelKey
	msgDecrypted, err2 := rsa.DecryptOAEP(sha1, nil, private, msgEncrypted, nil)	
	if err2 != nil {
		fmt.Println("#%d,%d error2: %s", err2)
	}
	
	fmt.Println("--------------------message after decryptOAEP ------- \n", msgDecrypted)
//************************golang EncryptOAEP & DecryptOAEP method test************************

	
	return shim.Success(nil)
	
}

var bccspInst b.BCCSP
//var AESKGen sw.aesKeyGenerator
//var AESPrivateK sw.aesPrivateKey

func AES256KeyGenForTest(stub shim.ChaincodeStubInterface) (pb.Response){
	
	fmt.Println("--------------------AES256KeyGenForTest start------- \n")

	kg := &sw.AesKeyGenerator{Length: 32}
	k, err := kg.KeyGen(&b.AES256KeyGenOpts{Temporary: true})
	aesK, _ := k.(*sw.AesPrivateKey)
	if err != nil {
		//return fmt.Printf("AES256KeyGenForTest error: KeyGen returned %s\n", err)
	}
	
	fmt.Println("--------------------KeyGen returned------- \n", aesK.PrivKey)
	
	
	var ptext = []byte("a 16 byte messag")

	fmt.Println("---------------------Before encrypting --ptext is--: '%s'  ----and key is-- %v \n", ptext, aesK.PrivKey)
	encrypted, encErr := sw.AesCBCEncrypt(aesK.PrivKey, ptext)
	
	fmt.Println("---------------------After encrypting --encrypted is--: '%s' \n", encrypted)
	if encErr != nil {
		fmt.Println("Error encrypting '%s': %v", ptext, encErr)
	}

	decrypted, decErr := sw.AesCBCDecrypt(aesK.PrivKey, encrypted)
	fmt.Println("---------------------After decrypting --decrypted is--: '%s' \n", decrypted)
	if decErr != nil {
		fmt.Println("Error decrypting '%s': %v", ptext, decErr)
	}

	return shim.Success(nil)
}



