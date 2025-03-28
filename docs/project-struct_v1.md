# GitHub CLI project layout

At a high level, these areas make up the project:

- [`cmd/`](../cmd) - `main` packages for building binaries such as the `gh` executable
- [`pkg/`](../pkg) - most other packages, including the implementation for individual gh commands
- [`docs/`](../docs) - documentation for maintainers and contributors
- [`script/`](../script) - build and release scripts
- [`internal/`](../internal) - Go packages highly specific to our needs and thus internal
- [`go.mod`](../go.mod) - external Go dependencies for this project, automatically fetched by Go at build time

Some auxiliary Go packages are at the top level of the project for historical reasons:

- [`api/`](../api) - main utilities for making requests to the Vault API
  <!-- - [`context/`](../context) - DEPRECATED: use only for referencing git remotes -->
  <!-- - [`git/`](../git) - utilities to gather information from a local git repository -->
  <!-- - [`test/`](../test) - DEPRECATED: do not use -->
  <!-- - [`utils/`](../utils) - DEPRECATED: use only for printing table output -->
