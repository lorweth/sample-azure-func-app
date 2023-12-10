# Simple app for learning Azure

## Deployment
- Open infrastructure folder `cd deploy/infrastructure`
- Run `echo .env.example > .env.dev`
- Login with azure cli `az login`
- Set the account with the Azure CLI `az account set --subscription "<subscription-id>"`
- Create a Service Principal `az ad sp create-for-rbac --role="Contributor" --scopes="/subscriptions/<SUBSCRIPTION_ID>"`
- Fill the `.env.prod` configurations.
- Run `make ready` and `make deploy`.

## References
- [Azure Functions Tips: Grouping Functions into Function Apps](https://marcduiker.dev/articles/azure-functions-grouping-functions-in-function-apps/)
- [Azure build with terraform](https://developer.hashicorp.com/terraform/tutorials/azure-get-started/azure-build)
- [Terraform register azurerm_linux_function_app](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/linux_function_app#example-usage)
- [Deploy Azure function with terraform](https://www.maxivanov.io/deploy-azure-functions-with-terraform/)

## Issues
- [`func` Command got error on Ubuntu container with Mac M1/M2](http://issamben.com/running-azure-function-as-docker-container-on-an-m1-m2/)
- [Deploying multiple function under same azure function app not working](https://stackoverflow.com/questions/57079549/deploying-multiple-function-under-same-azure-function-app-not-working)
