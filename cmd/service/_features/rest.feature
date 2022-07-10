Feature: Event Rest API

  Background: Start server
    Given the server is started
    And the seed data has been inserted

  Scenario Outline: Get Single Event

    As an API consumer, I want a REST API to retrieve a specific event 
    If the user authentication fails, return an authentication failure status code 
    If the ID provided is not found, return an appropriate failure status code 
    If the ID provided is found, but was not created by the user who is making the REST API call, return the same failure as if the ID was not  found 
    If the event of the ID has been deleted, return the same failure as if the ID was not found 
    If the ID provided is found and the user making the REST API call also created the event, then return the event. 

    When I perform a GET for ID <id> with token <token>
    Then the response should have http status <status>
#    And the response should have <count> events

Examples:
| id                       | token                            | status | count | description |
| 0063c3a5e4232e4cd0274ac2 | 74edf612f393b4eb01fbc2c29dd96671 | 200    | 1     | valid       |
| 1063c3a5e4232e4cd0274ac2 | 74edf612f393b4eb01fbc2c29dd96671 | 404    | 1     | not owner   |
| 0263c3a5e4232e4cd0274ac2 | 74edf612f393b4eb01fbc2c29dd96671 | 404    | 1     | deleted     |
