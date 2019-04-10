Feature: config
  In order to use mite in a sane manner
  As a mite user
  I need to be able to set defaults in a config

  Scenario: Write & read api key
    Given an empty config file called ".mite.toml"
    When I execute "-c .mite.toml config api.key=foo"
    Then "-c .mite.toml config api.key" should return "foo"
  Scenario: Write & read api url
    Given an empty config file called ".mite.toml"
    When I execute "-c .mite.toml config api.url=http://foo.invalid"
    Then "-c .mite.toml config api.url" should return "http://foo.invalid"
  Scenario: Write & read default projectId
    Given an empty config file called ".mite.toml"
    When I execute "-c .mite.toml config projectId=12345"
    Then "-c .mite.toml config projectId" should return "12345"
  Scenario: Write & read default serviceId
    Given an empty config file called ".mite.toml"
    When I execute "-c .mite.toml config serviceId=54321"
    Then "-c .mite.toml config serviceId" should return "54321"
