package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerateWallet(t *testing.T) {
	w := generateWallet("cosmos")
	require.Equal(t, w.Address[:7], "cosmos1", "Incorrect bech32 prefix")
	require.Equal(t, len(w.Address), 45, "Incorrect privkey length")
	require.Equal(t, len(w.Pubkey), 33, "Incorrect pubkey length")
	require.Equal(t, len(w.Privkey), 32, "Incorrect privkey length")
}

func TestGenerateWalletWithPrefix(t *testing.T) {
	w := generateWallet("osmo")
	require.Equal(t, w.Address[:5], "osmo1", "Incorrect bech32 prefix")
	require.Equal(t, len(w.Address), 43, "Incorrect privkey length")
	require.Equal(t, len(w.Pubkey), 33, "Incorrect pubkey length")
	require.Equal(t, len(w.Privkey), 32, "Incorrect privkey length")
}

func TestStartsWith(t *testing.T) {
	m := matcher{StartsWith: "aaaa"}
	require.True(t, m.Match("cosmos1aaaaqztg6eu45nlljp0wp947juded46aln83kr"))
	require.False(t, m.Match("cosmos1aaa9qztg6eu45nlljp0wp947juded46aln83kr"))
}

func TestEndsWith(t *testing.T) {
	m := matcher{EndsWith: "8888"}
	require.True(t, m.Match("cosmos14sy657pp6tgclhgqnl3dkwzwu3ewt4cf3f8888"))
	require.False(t, m.Match("cosmos14sy657pp6tgclhgqnl3dkwzwu3ewt4cf3ff888"))
}

func TestContains(t *testing.T) {
	m := matcher{Contains: "k2k2k"}
	require.True(t, m.Match("cosmos1s6rlmknaj3swdd7hua6s852sk2k2k409a3z9f9"))
	require.False(t, m.Match("cosmos14sy657pp6tgclhgqnl3dkwzwu3ewt4cf3ff888"))
}

func TestLetters(t *testing.T) {
	m := matcher{Letters: 38}
	require.True(t, m.Match("cosmos1gcjsgsglhacarlumkjzywedykkvkuvrzqlnlxd"))
	require.False(t, m.Match("cosmos1gcjsgsglhacarlumkjzywedykkvkuvrzqlnlx8"))
}

func TestDigits(t *testing.T) {
	m := matcher{Digits: 26}
	require.True(t, m.Match("cosmos1j666m3qz66t786s48t540536465p56zrve5893"))
	require.False(t, m.Match("cosmos1j666m3qz66t786s48t540536465p56zrve589z"))
}

func TestFindMatchingWalletConcurrent(t *testing.T) {
	goroutineCount := 5
	lastChars := "zz"
	m := matcher{EndsWith: lastChars}
	w := findMatchingWalletConcurrent(m, goroutineCount)
	require.Equal(t, w.Address[len(w.Address)-len(lastChars):], lastChars, "Incorrect address suffix")
}
