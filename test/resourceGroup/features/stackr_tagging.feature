Feature: The correct billing tags are applied at the resource group level
    In order to do billing correctly
    I need to rely on all resources being tagged in a consistent manner

    Scenario: A resource group is created through a terraform module used by Stackr
        Given A resource group is created
        When I check for tags against the resource group
        Then I expect to have at least the following tags present