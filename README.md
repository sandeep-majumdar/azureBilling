# AzureBilling

NOTE: IN AUGUST THIS CODE HAS BEEN REFACTORED INTO MODULES WITH A CLI TO ALLOW FOR ADDITIONAL CAPABILITY TO BE ADDED

PLEASE REVIEW THE RELEASE NOTES BELOW FOR CHANGES

## Overview

Please see the relevant readme file

- Original Azure billfile processing functionality is available in the [Billing README](billing/README.md)

- Rightsizing for Azure using the billfile: an experimental capability still under development, see [Rightsizing Readme](rightsizing/README.md)

NOTE: Rightsizing requires AZURE_AUTH_LOCATION to be configured to an appropriate SPN file with clientid and secret.

## RELEASE NOTES

The repo has been refactored and fitted with a cli to separate the functionality. The following cli commands are possible:

### Billfile processing

```bash
go run cmd/main.go azure billing -c ./config.json
```

### Rightsizing functionality

```bash
go run cmd/main.go azure rightsizing vm -c ./config.json 
```

### Arbitrarily run an az cli command

Useful, and can be used to test scripts used in the rightsizing recommendations.

```bash
# e.g. subscriptionID = 05ab872c-fa1a-42e7-9cf1-eafcdes9971d
go run cmd/main.go azure cli -l cli/examples/login.txt -f cli/examples/azureMonitor.txt -s <subscriptionID>
```
