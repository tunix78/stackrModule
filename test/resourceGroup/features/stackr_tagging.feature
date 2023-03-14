Feature: The correct billing tags are applied at the resource group level
    In order to do billing correctly
    I need to rely on all resources being tagged in a consistent manner

    Scenario: A resource group is created through a terraform module used by Stackr
        Given A resource group is planned via Terraform
        Then I expect to have at least the following tags present
        And I expect the location of the resource group to be one of the following
            | location   |
            | westeurope |
            | eastus     |
            | eastus2    |
