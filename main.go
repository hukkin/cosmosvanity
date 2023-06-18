package main

import (
	"encoding/hex"
	"fmt"
	"os"
	"runtime"
	"strings"

	flag "github.com/spf13/pflag"

	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

type matcher struct {
	StartsWith string
	EndsWith   string
	Contains   string
	Prefix     string
	Letters    int
	Digits     int
}

func (m matcher) Match(candidate string) bool {
	var builder strings.Builder
	builder.WriteString(m.Prefix)
	builder.WriteString("1") // build prefix string
	candidate = strings.TrimPrefix(candidate, builder.String())
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

func (m matcher) ValidationErrors() []string {
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
	Pubkey  []byte
	Privkey []byte
}

func (w wallet) String() string {
	return "Address:\t" + w.Address + "\n" +
		"Public key:\t" + hex.EncodeToString(w.Pubkey) + "\n" +
		"Private key:\t" + hex.EncodeToString(w.Privkey)
}

func generateWallet(prefix string) wallet {
	var privkey secp256k1.PrivKey = secp256k1.GenPrivKey()
	var pubkey secp256k1.PubKey = privkey.PubKey().(secp256k1.PubKey)
	bech32Addr, err := bech32.ConvertAndEncode(prefix, pubkey.Address())
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
			w := generateWallet(m.Prefix)
			if m.Match(w.Address) {
				// Do a non-blocking write instead of simple `ch <- w` to prevent
				// blocking when it's time to quit and ch is full.
				select {
				case ch <- w:
				default:
				}
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
	var walletsToFind = flag.IntP("count", "n", 1, "Amount of matching wallets to find")
	var cpuCount = flag.Int("cpus", runtime.NumCPU(), "Amount of CPU cores to use")

	var mustContain = flag.StringP("contains", "c", "", "A string that the address must contain")
	var mustStartWith = flag.StringP("startswith", "s", "", "A string that the address must start with")
	var mustEndWith = flag.StringP("endswith", "e", "", "A string that the address must end with")
	var letters = flag.IntP("letters", "l", 0, "Amount of letters (a-z) that the address must contain")
	var digits = flag.IntP("digits", "d", 0, "Amount of digits (0-9) that the address must contain")
	var bechPrefix = flag.StringP("prefix", "p", "cosmos", "The bech32 prefix the address should have")

	flag.Parse()

	if *walletsToFind < 1 {
		fmt.Println("ERROR: The number of wallets to generate must be 1 or more")
		os.Exit(1)
	}
	if *cpuCount < 1 {
		fmt.Println("ERROR: Must use at least 1 CPU core")
		os.Exit(1)
	}

	m := matcher{
		StartsWith: strings.ToLower(*mustStartWith),
		EndsWith:   strings.ToLower(*mustEndWith),
		Contains:   strings.ToLower(*mustContain),
		Prefix:     strings.ToLower(*bechPrefix),
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
		matchingWallet = findMatchingWalletConcurrent(m, *cpuCount)
		fmt.Printf(":::: Matching wallet %d/%d found ::::\n", i+1, *walletsToFind)
		fmt.Println(matchingWallet)
	}
}
