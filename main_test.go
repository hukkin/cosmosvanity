package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerateWallet(t *testing.T) {
	w := generateWallet()
	require.Equal(t, w.Address[:7], "cosmos1", "Incorrect bech32 prefix")
	require.Equal(t, len(w.Address), 45, "Incorrect privkey length")
	require.Equal(t, len(w.Pubkey), 33, "Incorrect pubkey length")
	require.Equal(t, len(w.Privkey), 32, "Incorrect privkey length")
}

func TestFindMatchingWallet(t *testing.T) {
	firstChars := "aa"
	m := matcher{StartsWith: firstChars}
	w := findMatchingWallet(m)
	require.Equal(t, w.Address[:7+len(firstChars)], "cosmos1"+firstChars, "Incorrect address prefix")
}

func TestFindMatchingWalletMultiProcess(t *testing.T) {
	lastChars := "zz"
	m := matcher{EndsWith: lastChars}
	w := findMatchingWalletMultiProcess(m)
	require.Equal(t, w.Address[len(w.Address)-len(lastChars):], lastChars, "Incorrect address suffix")
}
