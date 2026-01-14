package nostr

import (
	"context"
	"fmt"
	"time"

	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip19"
)

// Client represents a Nostr client
type Client struct {
	secretKey string
	relays    []string
	pool      *nostr.SimplePool
	relayPool map[string]*nostr.Relay
	cancel    context.CancelFunc
}

// Config holds Nostr client configuration
type Config struct {
	SecretKey string // nsec or hex
	Relays    []string
}

// NewClient creates a new Nostr client
func NewClient(cfg Config) (*Client, error) {
	// Normalize secret key if in nsec format
	secretKey := cfg.SecretKey
	if len(secretKey) > 0 && secretKey[:4] == "nsec" {
		_, hex, err := nip19.Decode(secretKey)
		if err != nil {
			return nil, fmt.Errorf("failed to decode nsec: %w", err)
		}
		secretKey = hex.(string)
	}

	return &Client{
		secretKey: secretKey,
		relays:    cfg.Relays,
		relayPool: make(map[string]*nostr.Relay),
	}, nil
}

// Connect connects to configured relays
func (c *Client) Connect() error {
	ctx, cancel := context.WithCancel(context.Background())
	c.cancel = cancel

	c.pool = nostr.NewSimplePool(ctx)

	for _, url := range c.relays {
		relay, err := nostr.RelayConnect(ctx, url)
		if err != nil {
			fmt.Printf("Warning: failed to connect to %s: %v\n", url, err)
			continue
		}
		c.relayPool[url] = relay
		fmt.Printf("Connected to relay: %s\n", url)
	}

	if len(c.relayPool) == 0 {
		return fmt.Errorf("failed to connect to any relays")
	}

	return nil
}

// Disconnect disconnects from all relays
func (c *Client) Disconnect() {
	if c.cancel != nil {
		c.cancel()
	}
	for url, relay := range c.relayPool {
		relay.Close()
		fmt.Printf("Disconnected from: %s\n", url)
	}
	c.relayPool = make(map[string]*nostr.Relay)
}

// Publish publishes a text note event to relays
func (c *Client) Publish(content string) error {
	if c.secretKey == "" {
		return fmt.Errorf("secret key not configured")
	}

	// Create event
	ev := nostr.Event{
		Kind:      1, // Text note
		Content:   content,
		CreatedAt: nostr.Timestamp(time.Now().Unix()),
	}

	// Sign event
	if err := ev.Sign(c.secretKey); err != nil {
		return fmt.Errorf("failed to sign event: %w", err)
	}

	// Publish to all connected relays
	ctx := context.Background()
	var lastErr error
	for url, relay := range c.relayPool {
		if err := relay.Publish(ctx, ev); err != nil {
			fmt.Printf("Failed to publish to %s: %v\n", url, err)
			lastErr = err
			continue
		}
		fmt.Printf("Published to %s\n", url)
	}

	return lastErr
}

// GetPublicKey returns the public key (npub) from the secret key
func (c *Client) GetPublicKey() (string, error) {
	if c.secretKey == "" {
		return "", fmt.Errorf("secret key not configured")
	}

	// Create a dummy event to extract public key
	ev := &nostr.Event{}
	if err := ev.Sign(c.secretKey); err != nil {
		return "", fmt.Errorf("failed to get public key: %w", err)
	}

	npub, err := nip19.EncodePublicKey(ev.PubKey)
	if err != nil {
		return "", fmt.Errorf("failed to encode npub: %w", err)
	}

	return npub, nil
}
