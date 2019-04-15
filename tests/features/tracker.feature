Feature: tracker
  In order to use mite in a sane manner
  As a mite user
  I need to be able to start, stop and show the status of my time tracker

  Scenario: time tracker status
    Given  A local mock server is setup for the http method "GET" and path "/tracker.json" which returns:
      """
      {
        "tracker": {
          "tracking_time_entry": {
            "id": 36135321,
            "minutes": 247,
            "since": "2015-10-15T17:05:04+02:00"
          }
        }
      }
      """
    And Mite is setup to connect to this mock server
    Then "-c .mite.toml tracker status" should return the following:
      """
      id        time  state     since
      --        ----  -----     -----
      36135321  4h7m  tracking  2015-10-15 15:05:04 +0000 UTC

      """
