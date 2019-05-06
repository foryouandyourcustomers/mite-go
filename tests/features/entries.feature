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
  Scenario: create entries
    Given  A local mock server is setup for the http method "POST" and path "/time_entries.json" which expects a body of:
      """
      {
        "time_entry": {
          "date_at": "2015-09-12",
          "minutes": 185,
          "note": "foo",
          "service_id": 243,
          "project_id": 123
        }
      }
      """
    And The mock server returns the following if the expectation is met:
      """
      {
        "time_entry": {
          "id": 52324,
          "minutes": 185,
          "date_at": "2015-09-12",
          "note": "foo",
          "billable": true,
          "locked": false,
          "revenue": null,
          "hourly_rate": 0,
          "user_id": 211,
          "user_name": "Fridolin Frei",
          "project_id": 123,
          "project_name": "Mite",
          "service_id": 243,
          "service_name": "Dokumentation",
          "created_at": "2015-09-13T18:54:45+02:00",
          "updated_at": "2015-09-13T18:54:45+02:00"
        }
      }
      """
    And Mite is setup to connect to this mock server
    Then "-c .mite.toml entries create -D 2015-09-12 -d 185m -p 123 -s 243 -n foo" should return the following:
      """
      id     notes  date        time  project  service
      --     -----  ----        ----  -------  -------
      52324  foo    2015-09-12  3h5m  Mite     Dokumentation
      """
  Scenario: edit entries
    Given  A local mock server is setup for the http method "GET" and path "/time_entries/52324.json" which expects a body of:
      """
      """
    And The mock server returns the following if the expectation is met:
      """
      {
        "time_entry": {
          "id": 52324,
          "minutes": 185,
          "date_at": "2015-09-12",
          "note": "foo",
          "billable": true,
          "locked": false,
          "revenue": null,
          "hourly_rate": 0,
          "user_id": 211,
          "user_name": "Fridolin Frei",
          "project_id": 123,
          "project_name": "Mite",
          "service_id": 243,
          "service_name": "Dokumentation",
          "created_at": "2015-09-13T18:54:45+02:00",
          "updated_at": "2015-09-13T18:54:45+02:00"
        }
      }
      """
    And A local mock server is setup for the http method "PATCH" and path "/time_entries/52324.json" which expects a body of:
      """
      {
        "time_entry": {
          "date_at": "2015-09-12",
          "minutes": 200,
          "note": "bar",
          "service_id": 243,
          "project_id": 123
        }
      }
      """
    And The mock server returns the following if the expectation is met:
      """
      """
    And Mite is setup to connect to this mock server
    Then "-c .mite.toml entries edit -i 52324 -D 2015-09-12 -d 200m -p 123 -s 243 -n bar" should return the following:
      """
      id     notes  date        time  project  service
      --     -----  ----        ----  -------  -------
      52324  foo    2015-09-12  3h5m  Mite     Dokumentation
      """
  Scenario: delete entries
    Given  A local mock server is setup for the http method "DELETE" and path "/time_entries/123.json" which returns:
      """
      """
    And Mite is setup to connect to this mock server
    Then "-c .mite.toml entries delete -i123" should return the following:
      """
      """


