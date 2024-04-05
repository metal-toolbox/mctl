### mctl

mctl is a CLI utility to interact with the metal-toolbox ecosystem of services.

### Getting started

0. Install the latest available version using `go install github.com/metal-toolbox/mctl@latest`.
1. Create a configuration file as `.mctl.yml`, for sample configuration files checkout [samples/mctl.yml](https://github.com/metal-toolbox/mctl/blob/main/samples).
2. Export `MCTLCONFIG=~/.mctl.yml`.

### Actions

For the updated list of all commands available, check out the [CLI docs](https://github.com/metal-toolbox/mctl/tree/main/docs/mctl.md)

- Get component information on a server - `mctl get component --server-id <>`
- List available firmware - `mctl list firmware`
- List firmware sets - `mctl list firmware-set`
- Retrieve information about a firmware - `mctl get firmware --id <>`
- Install a firmware set on a server - `mctl install firmware-set --server <>`
- Import firmware, firmware-set from file - `mctl create firmware-set  --from-file samples/fw-set.json`, where the JSON file contents is the output of `mctl list firmware-set`
