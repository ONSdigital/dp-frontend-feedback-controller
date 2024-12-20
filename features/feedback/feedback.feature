Feature: Feedback
    Scenario: When I navigate to the feedback page
        Given the feedback controller is running
        When I navigate to "/feedback"
        And the page should have the following content
        """
            {
                "#main h1": "Feedback"
            }
        """

    Scenario: When I click the submit button without filling the form
        Given the feedback controller is running
        When I navigate to "/feedback"
        Then element ".ons-btn" should be visible
        When I click the ".ons-btn" element
        Then element "#error-summary-title" should be visible
        Then element ".ons-panel--error" should be visible
