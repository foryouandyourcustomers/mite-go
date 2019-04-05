# Mite-go

A [Mite](https://mite.yo.lk/en/) time tracking command line interface.

## Setup

1. Grab a release from https://github.com/leanovate/mite-go/releases for your operating system and unpack it into your
$PATH (or %PATH% on windows).
2. Setup `mite-go` to use your API key by:
   1. visiting https://<your account name here>.mite.yo.lk/myself and note down the API key
   2. executing the following commands
   `mite-go config api.key=<your API key here>`
   `mite-go config api.url=https://<your account name here>.mite.yo.lk`
3. Optional: set a default project & service by:
   1. retrieving the desired project & service id by executing `mite-go projects` and `mite-go services` respectively
   2. configuring those id's as default by executing `mite-go config projectId=<the project id>` and `mite-go config serviceId=<the service id>`
4. Optional: mite-go allows you to define often used project & service combinations as activities. You can configure them by:
   1. think of a good name for the activity
   2. run `mite-go config activity.<your activity name here>.projectId=<the project id>`
   3. run `mite-go config activity.<your activity name here>.serviceId=<the service id>`
   4. the activity names can be used in the `entries create` and `entries edit` sub-commands

# Usage

Supported sub-commands:

| command  | functionality                                     |
|----------|---------------------------------------------------|
| config   | sets or reads a config property                   |
| entries  | lists & adds time entries                         |
| help     | Help about any command                            |
| projects | list & adds projects                              |
| services | list & adds services                              |
| tracker  | starts, stops and shows the status of the tracker |

For an up-to-date usage check `mite-go help`.
