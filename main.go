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
	Letters    int
	Digits     int
}

func (m *matcher) Match(candidate string) bool {
	candidate = strings.TrimPrefix(candidate, "cosmos1")
	if !strings.HasPrefix(candidate, m.StartsWith) {
		return false
	}
	if !strings.HasSuffix(candidate, m.EndsWith) {
		return false
	}
	if !strings.Contains(candidate, m.Contains) {
		return false
	}
	if countUnionChars(candidate, bech32digits) < m.Digits {
		return false
	}
	if countUnionChars(candidate, bech32letters) < m.Letters {
		return false
	}
	return true
}

func (m *matcher) ValidationErrors() []string {
	var errs []string
	if !bech32Only(m.Contains) || !bech32Only(m.StartsWith) || !bech32Only(m.EndsWith) {
		errs = append(errs, "ERROR: A provided matcher contains bech32 incompatible characters")
	}
	if len(m.Contains) > 38 || len(m.StartsWith) > 38 || len(m.EndsWith) > 38 {
		errs = append(errs, "ERROR: A provided matcher is too long. Must be max 38 characters.")
	}
	if m.Digits < 0 || m.Letters < 0 {
		errs = append(errs, "ERROR: Can't require negative amount of characters")
	}
	if m.Digits+m.Letters > 38 {
		errs = append(errs, "ERROR: Can't require more than 38 characters")
	}
	return errs
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

func findMatchingWallets(ch chan wallet, quit chan struct{}, m matcher) {
	for {
		select {
		case <-quit:
			return
		default:
			w := generateWallet()
			if m.Match(w.Address) {
				ch <- w
			}
		}
	}
}

func findMatchingWalletConcurrent(m matcher, goroutines int) wallet {
	ch := make(chan wallet)
	quit := make(chan struct{})
	defer close(quit)

	for i := 0; i < goroutines; i++ {
		go findMatchingWallets(ch, quit, m)
	}
	return <-ch
}

const bech32digits = "023456789"
const bech32letters = "acdefghjklmnpqrstuvwxyzACDEFGHJKLMNPQRSTUVWXYZ"

// This is alphanumeric chars minus chars "1", "b", "i", "o" (case insensitive)
const bech32chars = bech32digits + bech32letters

func bech32Only(s string) bool {
	return countUnionChars(s, bech32chars) == len(s)
}

func countUnionChars(s string, letterSet string) int {
	count := 0
	for _, char := range s {
		if strings.Contains(letterSet, string(char)) {
			count++
		}
	}
	return count
}

func main() {
	var walletsToFind *int = flag.IntP("count", "n", 1, "Amount of matching wallets to find")

	var mustContain *string = flag.StringP("contains", "c", "", "A string that the address must contain")
	var mustStartWith *string = flag.StringP("startswith", "s", "", "A string that the address must start with")
	var mustEndWith *string = flag.StringP("endswith", "e", "", "A string that the address must end with")
	var letters *int = flag.IntP("letters", "l", 0, "Amount of letters (a-z) that the address must contain")
	var digits *int = flag.IntP("digits", "d", 0, "Amount of digits (0-9) that the address must contain")
	flag.Parse()

	if *walletsToFind < 1 {
		fmt.Println("ERROR: The number of wallets to generate must be 1 or more")
		os.Exit(1)
	}

	m := matcher{
		StartsWith: strings.ToLower(*mustStartWith),
		EndsWith:   strings.ToLower(*mustEndWith),
		Contains:   strings.ToLower(*mustContain),
		Letters:    *letters,
		Digits:     *digits,
	}
	matcherValidationErrs := m.ValidationErrors()
	if len(matcherValidationErrs) > 0 {
		for i := 0; i < len(matcherValidationErrs); i++ {
			fmt.Println(matcherValidationErrs[i])
		}
		os.Exit(1)
	}

	var matchingWallet wallet
	for i := 0; i < *walletsToFind; i++ {
		matchingWallet = findMatchingWalletConcurrent(m, runtime.NumCPU())
		fmt.Printf(":::: Matching wallet %d/%d found ::::\n", i+1, *walletsToFind)
		fmt.Println(matchingWallet)
	}
}
