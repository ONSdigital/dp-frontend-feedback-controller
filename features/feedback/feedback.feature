Feature: Feedback
    Scenario: GET /feedback view feedback page
        Given the feedback controller is running
        When I navigate to "/feedback"
        Then element ".ons-fieldset__legend" should be visible
        Then element ".ons-radios__items" should be visible
        Then element ".ons-field" should be visible

    Scenario: When I click the submit button without filling the form
        Given the feedback controller is running
        When I navigate to "/feedback"
        Then element ".ons-btn" should be visible
        When I click the ".ons-btn" element
        Then element "#error-summary-title" should be visible
        Then element ".ons-panel--error" should be visible
