package eth_sign

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func Signature(message string, privateKeyStr string) (string, error) {
	hash := crypto.Keccak256Hash([]byte(message))
	privateKey, err := crypto.HexToECDSA(privateKeyStr)
	if err != nil {
		return "", err
	}
	signature, err := crypto.Sign(hash.Bytes(), privateKey)
	if err != nil {
		return "", err
	}
	return hexutil.Encode(signature), nil
}

func VerifySignature(message, privateKeyStr, signature string) (bool, error){
	signatureBytes, err := hexutil.Decode(signature)
	if err != nil {
		return false, err
	}
	privateKey, err := crypto.HexToECDSA(privateKeyStr)
	if err != nil {
		return false, err
	}
	hash := crypto.Keccak256Hash([]byte(message))
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return false, errors.New("error casting public key to ECDSA")
	}
	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	// remove recovery ID
	signatureNoRecoverID := signatureBytes[:len(signatureBytes)-1]
	verified := crypto.VerifySignature(publicKeyBytes, hash.Bytes(), signatureNoRecoverID)
	if !verified {
		return false, errors.New("verify signature error")
	}
	return true, nil
}

// PersonalSignature
// Sign calculates an Ethereum ECDSA signature for:
// keccak256("\x19Ethereum Signed Message:\n" + len(message) + message))
func PersonalSignature(message string, privateKeyStr string) (string, error) {
	fullMessage := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(message), message)
	hash := crypto.Keccak256Hash([]byte(fullMessage))
	privateKey, err := crypto.HexToECDSA(privateKeyStr)
	if err != nil {
		return "", err
	}
	signature, err := crypto.Sign(hash.Bytes(), privateKey)
	if err != nil {
		return "", err
	}
	signature[64] += 27
	return hexutil.Encode(signature), nil
}