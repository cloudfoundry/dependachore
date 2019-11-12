package kms

import (
	"context"
	"fmt"

	cloudkms "cloud.google.com/go/kms/apiv1"
	kmspb "google.golang.org/genproto/googleapis/cloud/kms/v1"
)

type Client struct {
	project  string
	location string
	keyRing  string
	key      string
}

func NewClient(project, location, keyRing, key string) *Client {
	return &Client{
		project:  project,
		location: location,
		keyRing:  keyRing,
		key:      key,
	}
}

func (c *Client) Encrypt(plainText []byte) ([]byte, error) {
	ctx := context.Background()
	client, err := cloudkms.NewKeyManagementClient(ctx)
	if err != nil {
		return []byte{}, fmt.Errorf("cloudkms.NewKeyManagementClient: %v", err)
	}

	req := &kmspb.EncryptRequest{
		Name:      c.keyPath(),
		Plaintext: plainText,
	}

	resp, err := client.Encrypt(ctx, req)
	if err != nil {
		return []byte{}, fmt.Errorf("client.Encrypt: %v", err)
	}

	return resp.Ciphertext, nil
}

func (c *Client) Decrypt(cipherText []byte) ([]byte, error) {
	ctx := context.Background()
	client, err := cloudkms.NewKeyManagementClient(ctx)
	if err != nil {
		return []byte{}, fmt.Errorf("cloudkms.NewKeyManagementClient: %v", err)
	}

	req := &kmspb.DecryptRequest{
		Name:       c.keyPath(),
		Ciphertext: cipherText,
	}

	resp, err := client.Decrypt(ctx, req)
	if err != nil {
		return []byte{}, fmt.Errorf("client.Decrypt: %v", err)
	}

	return resp.Plaintext, nil
}

func (c *Client) keyPath() string {
	return fmt.Sprintf("projects/%s/locations/%s/keyRings/%s/cryptoKeys/%s", c.project, c.location, c.keyRing, c.key)
}
