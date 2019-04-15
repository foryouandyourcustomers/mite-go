Feature: projects
  In order to use mite in a sane manner
  As a mite user
  I need to be able to list projects

  Scenario: list services
    Given  A local mock server is setup for the http method "GET" and path "/services.json" which returns:
      """
      [
        {
          "service": {
            "id": 38672,
            "name": "Website Konzeption",
            "note": "",
            "hourly_rate": 3300,
            "archived": false,
            "billable": true,
            "created_at": "2009-12-13T12:12:00+01:00",
            "updated_at": "2015-12-13T07:20:04+01:00"
          }
        }
      ]
      """
    And Mite is setup to connect to this mock server
    Then "-c .mite.toml services" should return the following:
      """
      id     name                notes
      --     ----                -----
      38672  Website Konzeption
      """
