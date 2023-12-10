provider "azurerm" {
    features {}
}

resource "azurerm_resource_group" "resource_group" {
    name     = "${var.project}-${var.environment}-resource-group"
    location = var.location
}

resource "azurerm_storage_account" "storage_account" {
    name                     = "${var.project}${var.environment}storage"
    resource_group_name      = azurerm_resource_group.resource_group.name
    location                 = azurerm_resource_group.resource_group.location
    account_tier             = "Standard"
    account_replication_type = "LRS"
}

resource "azurerm_application_insights" "application_insights" {
    name                = "${var.project}-${var.environment}-application-insights"
    location            = var.location
    resource_group_name = azurerm_resource_group.resource_group.name
    application_type    = "other"
}

resource "azurerm_service_plan" "service_plan" {
    name                = "${var.project}-${var.environment}-app-service-plan"
    resource_group_name = azurerm_resource_group.resource_group.name
    location            = var.location
    os_type             = "Linux"
    sku_name            = "Y1"
}

resource "azurerm_linux_function_app" "function_app" {
    location            = var.location
    name                = "${var.project}-${var.environment}-function-app"
    resource_group_name = azurerm_resource_group.resource_group.name
    service_plan_id     = azurerm_service_plan.service_plan.id

    storage_account_name = azurerm_storage_account.storage_account.name
    storage_account_access_key = azurerm_storage_account.storage_account.primary_access_key

    app_settings = {
        application_insights_connection_string = azurerm_application_insights.application_insights.connection_string
        application_insights_key = azurerm_application_insights.application_insights.instrumentation_key
    }

    site_config {}
}