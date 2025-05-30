package backend

import (
	"fmt"
	"strings"

	"github.com/dvcrn/chainenv/logger"
	"github.com/dvcrn/go-1password-cli/op"
)

type OnePasswordBackend struct {
	client    *op.Client
	vault     *op.Vault
	vaultName string

	logger *logger.Logger
}

func NewOnePasswordBackend(vaultName string, opts ...BackendOption) *OnePasswordBackend {
	options := newBackendOpts(opts...)
	logger := options.logger

	return &OnePasswordBackend{
		client:    op.NewOpClient(),
		logger:    logger,
		vaultName: vaultName,
	}
}

func (o *OnePasswordBackend) ensureVaultExists() error {
	vaults, err := o.client.Vaults()
	if err != nil {
		o.logger.Err("couldn't get vaults: %s", err.Error())
		return err
	}

	var vault *op.Vault
	for _, v := range vaults {
		if v.Name == o.vaultName {
			vault = v
			o.logger.Debug("Using existing 1Password vault: ID: %s, Name: %s, ContentVersion: %d\n", vault.ID, vault.Name, vault.ContentVersion)
			break
		}
	}

	if vault == nil {
		// Create a new vault
		var err error
		vault, err = o.client.CreateVault(o.vaultName,
			op.WithVaultDescription("Created by chainenv"),
			op.WithVaultIcon("treasure-chest"),
		)

		if err != nil {
			o.logger.Err("Error creating new 1Password vault: %s", err.Error())
			return err
		}

		o.logger.Info("Created new 1Password vault: ID: %s, Name: %s\n", vault.ID, vault.Name)
	}

	o.vault = vault
	return nil
}

func (o *OnePasswordBackend) GetPassword(account string) (string, error) {
	if err := o.ensureVaultExists(); err != nil {
		return "", fmt.Errorf("error ensuring vault exists: %v", err)
	}

	value, err := o.client.ReadItemField(o.vault.ID, account, "password")
	if err != nil {
		if strings.Contains(err.Error(), "isn't an item") {
			return "", fmt.Errorf("the item '%s' does not exist in the vault", account)
		}

		return "", fmt.Errorf("error retrieving password from 1Password: %v", err)
	}

	return value, nil
}

func (o *OnePasswordBackend) SetPassword(account, password string, update bool) error {
	if err := o.ensureVaultExists(); err != nil {
		return fmt.Errorf("error ensuring vault exists: %v", err)
	}

	vaultItem, _ := o.client.VaultItem(account, o.vault.ID)

	// If updating, first check if the item exists
	if update {
		o.logger.Debug("Running in update mode")

		if vaultItem == nil {
			return fmt.Errorf("item not found for update")
		}

		// If the item exists, update it
		editedItem, err := o.client.EditItemField(o.vault.ID, vaultItem.ID, op.Assignment{Name: "password", Value: password})
		if err != nil {
			return fmt.Errorf("error updating item in 1Password: %v", err)
		}

		o.logger.Debug("Updated item: %s, value: %s\n", editedItem.Title, password)

		return nil
	}

	if vaultItem != nil {
		return fmt.Errorf("item already exists. use 'update' to update.")
	}

	item, err := o.client.CreateItem(o.vault.ID, "password", account,
		op.WithItemTags([]string{"chainenv"}),
		op.WithItemAssignments([]op.Assignment{
			{Name: "password", Value: password},
			{Name: "notes", Value: fmt.Sprintf("This item was generated with `chainenv`. Access it with \n```\nchainenv get %s\n```", account)},
		}),
	)
	if err != nil {
		return fmt.Errorf("error creating item in 1Password: %v", err)
	}

	o.logger.Debug("Created item: %s\n", item.Title)

	return nil
}

func (o *OnePasswordBackend) List() ([]string, error) {
	if err := o.ensureVaultExists(); err != nil {
		return nil, fmt.Errorf("error ensuring vault exists: %v", err)
	}

	items, err := o.client.ItemsByVault(o.vault.ID, op.WithTags([]string{"chainenv"}))
	if err != nil {
		return nil, fmt.Errorf("error listing items in 1Password: %v", err)
	}

	var accounts []string
	for _, item := range items {
		accounts = append(accounts, item.Title)
	}

	return accounts, nil
}
