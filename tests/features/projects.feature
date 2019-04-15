Feature: projects
  In order to use mite in a sane manner
  As a mite user
  I need to be able to list projects

  Scenario: list projects
    Given  A local mock server is setup for the http method "GET" and path "/projects.json" which returns:
      """
      [
        {
          "project": {
            "id": 643,
            "name": "Open-Source",
            "note": "valvat, memento et all.",
            "customer_id": 291,
            "customer_name": "Yolk",
            "budget": 0,
            "budget_type": "minutes",
            "hourly_rate": 6000,
            "archived": false,
            "active_hourly_rate": "hourly_rate",
            "hourly_rates_per_service": [
              {
                "service_id": 31272,
                "hourly_rate": 4500
              },
              {
                "service_id": 149228,
                "hourly_rate": 5500
              }
            ],
            "created_at": "2011-08-17T12:06:57+02:00",
            "updated_at": "2015-02-19T10:53:10+01:00"
          }
        }
      ]
      """
    And Mite is setup to connect to this mock server
    Then "-c .mite.toml projects" should return the following:
      """
      id   name         notes
      --   ----         -----
      643  Open-Source  valvat, memento et all.
      """
