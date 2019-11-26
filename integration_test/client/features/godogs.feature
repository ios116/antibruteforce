Feature: antibruteforce
  In order to use grpc api
  As an GRPC user
  I need to be able to manage ip and buckets

  Scenario: check request if limit 10 request per 60s it is max request for login
    Given request set "Admin", "123456", "127.0.0.12"
    When check request 10 times
    Then response should be match "true"

  Scenario: check request if limit > 10 request per 60s
    And check request 10 times
    Then  response should be match "false"

  Scenario: check request if bucket login has been reset
    And reset Bucket "login", "Admin"
    When check request 10 times
    Then response should be match "true"

  Scenario: check request if ip in whitelist
    And add ip "127.0.0.12/32" to list "whitelist"
    When response should be match "true"
    When check request 10 times
    Then response should be match "true"

  Scenario: check request if ip in blacklist
    And remove ip "127.0.0.12/32" from list
    When response should be match "true"
    And add ip "127.0.0.12/32" to list "blacklist"
    When response should be match "true"
    When check request 10 times
    Then response should be match "false"
    And remove ip "127.0.0.12/32" from list



