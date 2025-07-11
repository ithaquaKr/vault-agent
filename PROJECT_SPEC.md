# Project spec

> This project help you manage your HashiCorp Vault Cluster in declarative way.

This project work in two mode:

- Agent Mode:
  - Vault-Agent auto sync your configurations to Vault.
- CLI Mode:
  - Manual operations with Vault using built-in commands.

Two type of configurations:

- AgentConfigure
- VaultInitialConfigure
- VaultDeclarativeData (Roles, Secrets, Policy, Entity,..)
