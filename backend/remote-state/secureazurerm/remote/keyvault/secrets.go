package keyvault

import (
	"context"
	"fmt"

	KV "github.com/Azure/azure-sdk-for-go/services/keyvault/2016-10-01/keyvault"
)

// getSecretVersion gets the secret version from ID.
func getSecretVersion(ID string) string {
	i := len(ID) - 1
	for ID[i] != '/' {
		i--
	}
	return ID[i+1 : len(ID)]
}

// InsertSecret inserts a secret into the key vault. Returns the version.
func (k *KeyVault) InsertSecret(ctx context.Context, name string, value string) (string, error) {
	vault, err := k.vaultClient.Get(ctx, k.resourceGroupName, k.vaultName)
	if err != nil {
		return "", fmt.Errorf("error getting key vault: %s", err)
	}

	contentType := "text/plain;charset=UTF-8"
	bundle, err := k.keyClient.SetSecret(ctx, *vault.Properties.VaultURI, name, KV.SecretSetParameters{Value: &value, ContentType: &contentType})
	if err != nil {
		return "", fmt.Errorf("error inserting secret: %s", err)
	}
	return getSecretVersion(*bundle.ID), nil
}

// GetSecret gets the secret named name from the key vault.
func (k *KeyVault) GetSecret(ctx context.Context, name string, version string) (string, error) {
	vault, err := k.vaultClient.Get(ctx, k.resourceGroupName, k.vaultName)
	if err != nil {
		return "", fmt.Errorf("error getting key vault: %s", err)
	}

	bundle, err := k.keyClient.GetSecret(ctx, *vault.Properties.VaultURI, name, version)
	if err != nil {
		return "", fmt.Errorf("error getting secret: %s", err)
	}
	return *bundle.Value, nil
}