# DAPme Registry

A static DAP registry ([source](https://avparker.github.io/dap-registry-static/)).

# "Registering" a DAP
1. Activate hermit
2. run `just register <your_handle>`
3. add whatever money addresses you want (including the "urn:" prefix)

> [!WARNING]
> Currently overwrites any pre-existing DAP.

# How it works
* `dapme.xyz`'s DID document is hosted at `https://dapme.xyz/.well-known/did.json` and located [here](./registry/.well-known/did.json)
* creates a `did:web` for all handles e.g. `did:web:dapme.xyz:moegrammer` by saving the auto-generated did document to `registry/<handle>/did.json`
* saves the registration to `registry/daps/<handle>`