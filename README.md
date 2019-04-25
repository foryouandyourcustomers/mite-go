# Mite

A [mite](https://mite.yo.lk/en/) time tracking command line interface.

## Status

[![Build Status](https://travis-ci.org/leanovate/mite-go.svg?branch=master)](https://travis-ci.org/leanovate/mite-go)

## Setup

1. Grab a release from https://github.com/leanovate/mite-go/releases for your operating system and unpack it into your
$PATH (or %PATH% on windows).
2. Make sure that the `mite` command is executable by executing `mite version` in your shell
3. Setup `mite` to use your API key by:
   1. visiting https://"your account name here".mite.yo.lk/myself and note down the API key
   2. executing the following commands
   `mite config api.key="your API key here"`
   `mite config api.url=https://"your account name here".mite.yo.lk`
4. Optional: set a default project & service by:
   1. retrieving the desired project & service id by executing `mite projects` and `mite services` respectively
   2. configuring those id's as default by executing `mite config projectId="the project id"` and `mite config serviceId="the service id"`
5. Optional: mite allows you to define often used project & service combinations as activities. You can configure them by:
   1. think of a good name for the activity
   2. run `mite config activity."your activity name here".projectId="the project id"`
   3. run `mite config activity."your activity name here".serviceId=<the service id"`
   4. the activity names can be used in the `entries create` and `entries edit` sub-commands
6. Optional: set a project & service for your vacation tracking by:
   1. retrieving the desired project & service id by executing `mite projects` and `mite services` respectively
   2. configuring those id's as default by executing `mite config vacation.projectId="the project id"` and `mite config vacation.serviceId="the service id"`

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

For an up-to-date usage check `mite help`.
