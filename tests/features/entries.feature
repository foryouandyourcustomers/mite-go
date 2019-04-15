Feature: entries
  In order to use mite in a sane manner
  As a mite user
  I need to be able to add, remove, edit and list time entries

  Scenario: list entries
    Given  A local mock server is setup for the http method "GET" and path "/time_entries.json" which returns:
      """
      [
        {
          "time_entry": {
            "id": 36159117,
            "minutes": 15,
            "date_at": "2015-10-16",
            "note": "Feedback einarbeiten",
            "billable": true,
            "locked": false,
            "revenue": null,
            "hourly_rate": 0,
            "user_id": 211,
            "user_name": "Fridolin Frei",
            "project_id": 88309,
            "project_name": "API v2",
            "customer_id": 3213,
            "customer_name": "König",
            "service_id": 12984,
            "service_name": "Entwurf",
            "created_at": "2015-10-16T12:39:00+02:00",
            "updated_at": "2015-10-16T12:39:00+02:00"
          }
        }
      ]
      """
    And Mite is setup to connect to this mock server
    Then "-c .mite.toml entries" should return the following:
      """
      id        notes                 date        time  project  service
      --        -----                 ----        ----  -------  -------
      36159117  Feedback einarbeiten  2015-10-16  15m   API v2   Entwurf

      """
