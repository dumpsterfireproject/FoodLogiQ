Feature: Authentication

  Background: Start server
    Given the server is started

  Scenario: No Bearer Token
    When I perform a GET request with no token
    Then the response should have http status 401

  Scenario Outline: Validate Bearer Token
    When I perform a GET request with bearer token <token>
    Then the response should have http status <status>
Examples:
    | token                            | status |
    | 74edf612f393b4eb01fbc2c29dd96671 | 200    |
    | d88b4b1e77c70ba780b56032db1c259b | 200    |
    | 00edf612f393b4eb01fbc2c29dd96671 | 403    |
