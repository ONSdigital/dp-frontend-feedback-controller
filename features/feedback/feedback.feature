Feature: Feedback
    Scenario: GET /feedback view feedback page
        Given the feedback controller is running
        When I navigate to "/feedback"
        And element "ons-fieldset__legend" should be visible
        And element "ons-radios__items" should be visible
        And element "ons-field" should be visible

    Scenario: GET /feedback view feedback page
        Given the feedback controller is running
        When I navigate to "/feedback"
        And element "ons-btn" should be visible
        When I click the "ons-btn" element
        Then element "ons-panel__body" should be visible

