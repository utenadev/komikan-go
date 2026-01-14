package main

import (
	"fmt"

	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip19"
)

func main() {
	// Generate new key pair
	privKey := nostr.GeneratePrivateKey()
	pubKey, _ := nostr.GetPublicKey(privKey)

	// Encode to nsec and npub
	nsec, _ := nip19.EncodePrivateKey(privKey)
	npub, _ := nip19.EncodePublicKey(pubKey)

	fmt.Println("=== Nostr Key Pair Generated ===")
	fmt.Printf("Secret Key (hex): %s\n", privKey)
	fmt.Printf("Public Key (hex): %s\n", pubKey)
	fmt.Println()
	fmt.Printf("nsec (for config.yaml): %s\n", nsec)
	fmt.Printf("npub (your public address): %s\n", npub)
	fmt.Println()
	fmt.Println("Add this to config.yaml:")
	fmt.Printf("  nostr:\n")
	fmt.Printf("    secret_key: \"%s\"\n", nsec)
}
