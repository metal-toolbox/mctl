### mctl

mctl is a CLI utility to interact with Server service.

### Getting started

Create a `.mctl.yml` file with the below contents and place it in your home
directory.

```
serverservice_endpoint: <>
oidc_issuer_endpoint: <>
oidc_audience: <>
oidc_client_id: <>
```

### Run queries

- First authenticate with `mctl auth`
- Run queries `mctl list firmware`
