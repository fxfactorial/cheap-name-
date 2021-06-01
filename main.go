package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
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

func permutate(input string, parent string) []string {
	if len(input) == 1 {
		return []string{parent + input}
	}

	var permutations []string
	for i := 0; i < len(input); i++ {
		restOfInput := input[0:i] + input[i+1:]
		curChar := input[i : i+1]
		permutations = append(permutations, permutate(restOfInput, parent+curChar)...)
	}
	return permutations
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
	t0 := time.Now()
	fmt.Println("generating all the permutations ")

	all := permutate("abcdefghijk", "")
	fmt.Println("took", time.Since(t0), "to generate all permutations")

	length := len(all)
	chopped := length / 4

	fmt.Println("kicking off 4 threads")

	for i := 0; i < 4; i++ {
		fmt.Println("range", i*chopped, i*chopped+chopped)
		subrange := all[i*chopped : i*chopped+chopped]
		go func() {
			for _, a := range subrange {
				combined := string(a) + "(" + sig + ")"
				b := crypto.Keccak256([]byte(combined))[:4]

				if bytes.Equal(wanted, b) {
					fmt.Println("FOUND exactly all zeros after",
						time.Since(started), "signature should be", combined,
					)
					os.Exit(0)
					return
				}

				if p == false && bytes.Compare(b, atMost) == -1 {
					fmt.Println("this is good enough - can do ctrl-c now",
						"use this as your signature",
						combined, "found after",
						time.Since(started),
						hexutil.Encode(b),
					)
				}
			}
		}()
	}

	select {}
}
