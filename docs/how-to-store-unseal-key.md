# How to store Unseal key

- The best practices for securely to storing Vault unseal key and token key depends on organization's specific security requirements.
- Need to define a `Thread model`

- If store unseal key locally on the disk of Vault server => Easy to maintain, or run script but anyone has access --> Scan decrypt Vault.
- Using online Key Management system can has it own problem like leaking permission to get key,...

## Common solutions

- Using KV service like KMS
- alibabakms
- alibabaoss
- awskms
- azurek
- dev
- file
- gckms
- gcs
- hsm
- k8s
- multi
- oci
- ocikms
- s3
- vault

## A private solutions

- Send a encrypt unseal key to maintainer via internal email system.
