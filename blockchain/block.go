package blockchain

import (
	"bytes"
	"log"
	"crypto/sha256"
	"github.com/mr-tron/base58/base58"
	"golang.org/x/crypto/ed25519"
	cbor "github.com/whyrusleeping/cbor/go"
	"fmt"
)

type Block struct {
	Height        int			`json:"height"`
	Hash		  []byte		`json:"hash"`
	Transactions  []Transaction	`json:"transactions"`
	PreviousBlock []byte		`json:"previous_block"`
	Difficulty int				`json:"difficulty"`
	Nonce int					`json:"nounce"`
}

const COINBASE_AMOUNT = 25

func GenerateGenesisBlock(publicKey ed25519.PublicKey, privateKey ed25519.PrivateKey) (Block, error) {
	coinbase := GenerateCoinbase(publicKey, COINBASE_AMOUNT)
	err := coinbase.Sign(privateKey, 0)
	if err != nil {
		return Block{}, err
	}

	// change this to a static nonce once mining algorithm is implemented
	block := Block{0, []byte{}, []Transaction{coinbase}, []byte{}, 3, 1}
	hash, err := block.GetHash()
	if err != nil {
		return Block{}, err
	}

	block.Hash = hash
	return block, err
}

func (b *Block) GetCBOR() (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	enc := cbor.NewEncoder(buf)
	err := enc.Encode(b)

	if err != nil {
		log.Fatal("Error decoding ", err)
		return new(bytes.Buffer), err
	}

	return buf, err
}

func (b *Block) GetHash() ([]byte, error) {
	b.Hash = []byte{}
	block, err := b.GetCBOR()
	if err != nil {
		log.Fatal("Error encoding to CBOR %s", err)
		return []byte{}, err
	}
	hasher := sha256.New()
	hasher.Write(block.Bytes())
	return hasher.Sum(nil), err
}

func (b *Block) GetBase58Hash() (string, error) {
	hash, err := b.GetHash()
	if err != nil {
		return "", err
	}
	return base58.Encode(hash), err
}

func HashMatchesDifficulty(hash []byte, difficulty int) bool {
	for i, n := range(hash) {
		prefix := fmt.Sprintf("%b", n)
		if prefix != "0" && i < difficulty {
			return false
		}

		if prefix == "0" && i == difficulty {
			break
		}
	}
	return true
}
