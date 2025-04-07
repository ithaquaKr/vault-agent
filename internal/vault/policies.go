package vault

import (
	"fmt"
	"log/slog"

	"github.com/hashicorp/vault-client-go/schema"
)

func (v *vaultServer) SyncPolicy() error {
	managedPolicies, err := initPolicies(v.config.Policies)
	if err != nil {
		return fmt.Errorf("error while initialing policies config: %s", err)
	}
	if err := v.addManagedPolicies(managedPolicies); err != nil {
		return fmt.Errorf("error while adding policies: %s", err)
	}
	if err := v.removeUnmanagedPolicies(managedPolicies); err != nil {
		return fmt.Errorf("error while deleting unmanaged policies: %s", err)
	}

	return nil
}

type policy struct {
	Name           string `mapstructure:"name"`
	Rules          string `mapstructure:"rules"`
	RulesFormatted string
}

func initPolicies(policiesConfig []policy) ([]policy, error) {
	// 	for index, policy := range policiesConfig {
	// 		for k, v := range mounts {
	// 			policy.Rules = strings.ReplaceAll(policy.Rules, fmt.Sprintf("__accessor__%s", strings.TrimRight(k, "/")), v.Accessor)
	// 		}
	// 		//
	// 		// Format HCL polices.
	// 		rulesFormatted, err := hclPrinter.Format([]byte(policy.Rules))
	// 		if err != nil {
	// 			// Check if rules parse (HCL or JSON).
	// 			if _, err := hcl.Parse(policy.Rules); err != nil {
	// 				return nil, fmt.Errorf("error parsing %s policy rules: %s", policy.Name, err)
	// 			}
	//
	// 			// Policies are parsable but couldn't be HCL formatted (most likely JSON).
	// 			rulesFormatted = []byte(policy.Rules)
	// 			slog.Debug(fmt.Sprintf("error HCL-formatting %s policy rules (ignore if rules are JSON-formatted): %s",
	// 				policy.Name, err.Error()))
	// 		}
	// 		policiesConfig[index].RulesFormatted = string(rulesFormatted)
	// 	}
	//
	return policiesConfig, nil
}

// addManagedPolicies add defined policies to Vault
func (v *vaultServer) addManagedPolicies(managedPolicies []policy) error {
	for _, policy := range managedPolicies {
		slog.Info(fmt.Sprintf("adding policy %s", policy.Name))
		reqSchema := schema.PoliciesWriteAclPolicyRequest{
			Policy: policy.RulesFormatted,
		}
		if _, err := v.cl.System.PoliciesWriteAclPolicy(v.ctx, policy.Name, reqSchema); err != nil {
			return fmt.Errorf("error putting %s policy into vault: %s", policy.Name, err)
		}
	}

	return nil
}

// getExistingPolicies get all policies that are already in Vault
func (v *vaultServer) getExistingPolicies() (map[string]bool, error) {
	existingPolicies := make(map[string]bool)

	resp, err := v.cl.System.PoliciesListAclPolicies(v.ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to list existing policies: %s", err)
	}
	existingPoliciesList := resp.Data.Policies

	for _, existingPolicy := range existingPoliciesList {
		existingPolicies[existingPolicy] = true
	}

	return existingPolicies, nil
}

// getUnanagedPolicies gets unmanaged policies by comparing what's already in Vault and what's in the externalConfig.
func (v *vaultServer) getUnanagedPolicies(managedPolicies []policy) map[string]bool {
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
func (v *vaultServer) removeUnmanagedPolicies(managedPolicies []policy) error {
	// TODO: Check if has configure purge unmanaged config
	// if !v.config.PurgeUnmanagedConfig.Enabled || v.externalConfig.PurgeUnmanagedConfig.Exclude.Policies {
	// 	slog.Debug("purge config is disabled, no unmanaged policies will be removed")
	// 	return nil
	// }

	unmanagedPolicies := v.getUnanagedPolicies(managedPolicies)
	for policyName := range unmanagedPolicies {
		slog.Info(fmt.Sprintf("removing policy %s", policyName))
		if _, err := v.cl.System.PoliciesDeleteAclPolicy(v.ctx, policyName); err != nil {
			return fmt.Errorf("error deleting %s policy from vault: %s", policyName, err)
		}
	}
	return nil
}
