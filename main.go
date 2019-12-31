package main

import (
	"encoding/hex"
	"fmt"
	"os"
	"runtime"
	"strings"

	flag "github.com/spf13/pflag"

	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/libs/bech32"
)

type matcher struct {
	StartsWith string
	EndsWith   string
	Contains   string
}

func (m *matcher) Match(candidate string) bool {
	if !strings.HasPrefix(candidate, m.StartsWith) {
		return false
	}
	if !strings.HasSuffix(candidate, m.EndsWith) {
		return false
	}
	if !strings.Contains(candidate, m.Contains) {
		return false
	}
	return true
}

type wallet struct {
	Address string
	Pubkey  [33]byte
	Privkey [32]byte
}

func (w wallet) String() string {
	return "Address:\t" + w.Address + "\n" +
		"Public key:\t" + hex.EncodeToString(w.Pubkey[:]) + "\n" +
		"Private key:\t" + hex.EncodeToString(w.Privkey[:])
}

func generateWallet() wallet {
	var privkey secp256k1.PrivKeySecp256k1 = secp256k1.GenPrivKey()
	var pubkey secp256k1.PubKeySecp256k1 = privkey.PubKey().(secp256k1.PubKeySecp256k1)
	bech32Addr, err := bech32.ConvertAndEncode("cosmos", pubkey.Address())
	if err != nil {
		panic(err)
	}

	return wallet{bech32Addr, pubkey, privkey}
}

func findMatchingWallet(m matcher) wallet {
	for {
		w := generateWallet()
		trimmedAdress := strings.TrimPrefix(w.Address, "cosmos1")
		if m.Match(trimmedAdress) {
			return w
		}
	}
}

func findMatchingWallets(ch chan wallet, m matcher) {
	for {
		ch <- findMatchingWallet(m)
	}
}

func findMatchingWalletMultiProcess(m matcher) wallet {
	ch := make(chan wallet)
	for i := 0; i < runtime.NumCPU(); i++ {
		go findMatchingWallets(ch, m)
	}
	return <-ch
}

// This is alphanumeric chars minus chars "1", "b", "i", "o" (case insensitive)
const bech32chars = "acdefghjklmnpqrstuvwxyzACDEFGHJKLMNPQRSTUVWXYZ023456789"

func bech32Only(s string) bool {
	for _, char := range s {
		if !strings.Contains(bech32chars, string(char)) {
			return false
		}
	}
	return true
}

func main() {
	var mustContain string
	var mustStartWith string
	var mustEndWith string
	flag.StringVarP(&mustContain, "contains", "c", "", "A string that the address must contain")
	flag.StringVarP(&mustStartWith, "startswith", "s", "", "A string that the address must start with")
	flag.StringVarP(&mustEndWith, "endswith", "e", "", "A string that the address must end with")
	flag.Parse()

	mustContain = strings.ToLower(mustContain)
	mustStartWith = strings.ToLower(mustStartWith)
	mustEndWith = strings.ToLower(mustEndWith)
	if !bech32Only(mustContain) || !bech32Only(mustStartWith) || !bech32Only(mustEndWith) {
		fmt.Println("ERROR: A provided matcher flag contains bech32 incompatible chars")
		os.Exit(1)
	}
	m := matcher{
		StartsWith: mustStartWith,
		EndsWith:   mustEndWith,
		Contains:   mustContain,
	}

	matchingWallet := findMatchingWalletMultiProcess(m)
	fmt.Println("A match found:")
	fmt.Println(matchingWallet)
}
