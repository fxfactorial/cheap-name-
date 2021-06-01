package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	signature        = flag.String("selector", "", "selector given")
	mustBeZeros      = flag.Bool("all_zeros", false, "force search for zeros")
	found       bool = false
)

// Perm calls f with each permutation of a.
func Perm(a []rune, f func([]rune)) {
	perm(a, f, 0)
}

// Permute the values at index i to len(a)-1.
func perm(a []rune, f func([]rune), i int) {
	if i > len(a) {
		f(a)
		return
	}
	if found {
		return
	}
	perm(a, f, i+1)
	for j := i + 1; j < len(a); j++ {
		a[i], a[j] = a[j], a[i]
		perm(a, f, i+1)
		a[i], a[j] = a[j], a[i]
	}
}

func main() {
	flag.Parse()

	if *signature == "" {
		log.Fatal("no selector given")
	}

	sig := *signature

	atMost := []byte{0x00, 0x00, 0x00, 0xff}
	wanted := []byte{0x00, 0x00, 0x00, 0x00}
	started := time.Now()
	p := *mustBeZeros

	Perm([]rune("abcdefghijklmn"), func(a []rune) {
		combined := string(a) + "(" + sig + ")"
		b := crypto.Keccak256([]byte(combined))[:4]

		if bytes.Equal(wanted, b) {
			fmt.Println("FOUND exactly all zeros after",
				time.Since(started), "signature should be", combined,
			)
			return
		}

		if p == false && bytes.Compare(b, atMost) == -1 {
			fmt.Println("this is good enough - can do ctrl-c now",
				"use this as your signature",
				combined, "found after",
				time.Since(started),
				hexutil.Encode(b),
			)

			found = true
			return
		}
	})

	fmt.Println("odd didnt not find match - please report ticket ")
}
