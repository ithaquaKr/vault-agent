package vault

import (
	"fmt"
	"log/slog"

	hclPrinter "github.com/hashicorp/hcl/v2/hclwrite"

	"github.com/ithaquaKr/vault-agent/pkg/config"
)

func (v *vaultManager) SyncPolicy() error {
	managedPolicies, err := v.initPolicies(v.data.Policies)
	if err != nil {
		return fmt.Errorf("error while initialing policies config: %w", err)
	}
	if err := v.addManagedPolicies(managedPolicies); err != nil {
		return fmt.Errorf("error while adding policies: %w", err)
	}
	if err := v.removeUnmanagedPolicies(managedPolicies); err != nil {
		return fmt.Errorf("error while deleting unmanaged policies: %w", err)
	}

	return nil
}

type policy struct {
	Name           string `mapstructure:"name"`
	Rules          string `mapstructure:"rules"`
	RulesFormatted string
}

func (v *vaultManager) initPolicies(policiesConfig []config.Policy) ([]policy, error) {
	var p []policy
	for _, policyConfig := range policiesConfig {
		// TODO: Templating policies with mounts
		// for k, v := range mounts {
		// 	policy.Rules = strings.ReplaceAll(policy.Rules, fmt.Sprintf("__accessor__%s", strings.TrimRight(k, "/")), v.Accessor)
		// }

		// Format HCL polices.
		rulesFormatted := hclPrinter.Format([]byte(policyConfig.Rules))
		p = append(p, policy{
			Name:           policyConfig.Name,
			Rules:          policyConfig.Rules,
			RulesFormatted: string(rulesFormatted),
		})
	}

	return p, nil
}

// addManagedPolicies add defined policies to Vault
func (v *vaultManager) addManagedPolicies(managedPolicies []policy) error {
	for _, policy := range managedPolicies {
		slog.Info(fmt.Sprintf("adding policy %s", policy.Name))
		if err := v.cl.Sys().PutPolicy(policy.Name, policy.RulesFormatted); err != nil {
			return fmt.Errorf("error putting %s policy into vault: %w", policy.Name, err)
		}
	}

	return nil
}

// getExistingPolicies get all policies that are already in Vault
func (v *vaultManager) getExistingPolicies() (map[string]bool, error) {
	existingPolicies := make(map[string]bool)

	existingPoliciesList, err := v.cl.Sys().ListPolicies()
	if err != nil {
		return nil, fmt.Errorf("unable to list existing policies: %w", err)
	}

	for _, existingPolicy := range existingPoliciesList {
		existingPolicies[existingPolicy] = true
	}

	return existingPolicies, nil
}

// getUnanagedPolicies gets unmanaged policies by comparing what's already in Vault and what's in the externalConfig.
func (v *vaultManager) getUnanagedPolicies(managedPolicies []policy) map[string]bool {
	policies, _ := v.getExistingPolicies()

	// Vault doesn't allow to remove default or root policies.
	delete(policies, "root")
	delete(policies, "default")

	// Remove managed polices form the items since the reset will be removed.
	for _, managedPolicy := range managedPolicies {
		delete(policies, managedPolicy.Name)
	}

	return policies
}

// removeUnmanagedPolicies remove the unmanaged policies in Vault
func (v *vaultManager) removeUnmanagedPolicies(managedPolicies []policy) error {
	// TODO: Check if has configure purge unmanaged config
	// if !v.config.PurgeUnmanagedConfig.Enabled || v.externalConfig.PurgeUnmanagedConfig.Exclude.Policies {
	// 	slog.Debug("purge config is disabled, no unmanaged policies will be removed")
	// 	return nil
	// }

	unmanagedPolicies := v.getUnanagedPolicies(managedPolicies)
	for policyName := range unmanagedPolicies {
		slog.Info(fmt.Sprintf("removing policy %s", policyName))
		if err := v.cl.Sys().DeletePolicy(policyName); err != nil {
			return fmt.Errorf("error deleting %s policy from vault: %w", policyName, err)
		}
	}
	return nil
}
