Feature: Feedback

    Scenario: When I navigate to the feedback page
      Given the feedback controller is running
      When I navigate to "/feedback"
      Then the page should have the following content
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
        And element ".ons-panel--error" should be visible

    Scenario: When I submit the form selecting whole website with no feedback
        Given the feedback controller is running
        When I navigate to "/feedback"
        Then I click the "#whole-site" element
        When I click the ".ons-btn" element
        Then element ".ons-panel--error" should be visible
        And the page should have the following content
        """
            {
                "#main h2": "There is a problem with this page"
            }
        """

    Scenario: When I submit the form selecting specific page with invalid url
        Given the feedback controller is running
        When I navigate to "/feedback"
        Then I click the "#specific-page" element
        Then I fill in input element "#page-url-field" with value "https://some-url.net"
        Then I fill in input element "#description-field" with value "good and useful website"
        When I click the ".ons-btn" element
        Then element ".ons-panel--error" should be visible
        And the page should have the following content
        """
            {
                "#main h2": "There is a problem with this page"
            }
        """

    Scenario: When I submit the form selecting specific page with no url
        Given the feedback controller is running
        When I navigate to "/feedback"
        Then I click the "#specific-page" element
        Then I fill in input element "#description-field" with value "good and useful website"
        When I click the ".ons-btn" element
        Then element ".ons-panel--error" should be visible
        And the page should have the following content
        """
            {
                "#main h2": "There is a problem with this page"
            }
        """

    Scenario: When I submit the form selecting specific page with no url, no feedback
        Given the feedback controller is running
        When I navigate to "/feedback"
        Then I click the "#specific-page" element
        When I click the ".ons-btn" element
        Then element ".ons-panel--error" should be visible
        And the page should have the following content
        """
            {
                "#main h2": "There are 2 problems with this page"
            }
        """

    Scenario: When I submit the form with feedback
        Given the feedback controller is running
        When I navigate to "/feedback"
        And I click the "#whole-site" element
        When I fill in input element "#description-field" with value "good and useful website"
        When I click the ".ons-btn" element
        Then I navigate to "/feedback/thanks"
        And the page should have the following content
        """
            {
                "#main .ons-panel__body": "Thank you\nYour feedback will help us to improve the website. We are unable to respond to all enquiries. If your matter is urgent, please contact us.",
                "#main .ons-js-submit-btn": "Done"
            }
        """
