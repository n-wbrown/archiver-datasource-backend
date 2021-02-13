# Grafana Data Source Backend Plugin for Epics Archiver

## Introduction

`archiver-datasource-backend` is a [Grafana](https://grafana.com/) plugin to use data from the [Epics Archiver](https://slacmshankar.github.io/epicsarchiver_docs/) in Grafana dashboards and alerts.

This plugin adds a back-end query system to the front-end plugin [archiverappliance-datasource](https://github.com/sasaki77/archiverappliance-datasource), created by [Shinya Sasaki](https://github.com/sasaki77). The front-end plugin is included as a git submodule and may not reflect its latest iteration.

This plugin was built using the [grafana-sdk-go-plugin](https://github.com/grafana/grafana-plugin-sdk-go)


## Build Status
| Build: | Status: | 
| :---: | :---: | 
| Latest release | ![alt_text](https://github.com/n-wbrown/archiver-datasource-backend/workflows/publish-release.yml/badge.svg) |
| Master Branch | ![alt_text](https://github.com/n-wbrown/archiver-datasource-backend/workflows/build-and-test.yml/badge.svg?branch=master)

## Features

### Feature Status

| Feature: | Complete: |
| :---: | :---: | 
| Scalar PV support | Complete |
| Alerts for scalar data | Complete | 
| Archiver operators | Complete | 
| Functions for scalar data | Complete | 
| Waveform PV support | Incomplete |
| Functions for waveform data | Incomplete |

### Testing New Features

If you're unsure if a new alert is working properly, you can test it using the `Test rule` button found at the bottom of the Alert tab. If the response includes a message like `state:"no_data"` the alert may be misconfigured, dependent upon an incomplete feature, or merits a bug report. 

## Getting Started: For Users

### Installation 

1. Download a pre-built release from the releases page.

2. Unzip the contents in the Grafana plugins folder.

3. This plugin is unsigned. It must be specially listed by name in the Grafana `configure.ini` file to allow Grafana to use it. Add `https://github.com/n-wbrown/archiver-datasource-backend` to the `allow_loading_unsigned_plugins` parameter in the `[plugins]` section.

4. Restart Grafana.

### Configuring The Plugin

1. Log into Grafana's web interface as an administrator. Go to the "Configuration / Data Sources" page and click the "Add data source" button. Locate `archiver-datasource-backend` near the bottom of the list. 

2. To configure the plugin, add the url of the archiver. This will likely take the form `https://[hostname]/retrieval` where `hostname` will be unique to your case. 

## Getting Started: For Developers

A data source backend plugin consists of both frontend and backend components.

### Front-end Development

Node version v12.19.0 is recommended. If you're new to the Node.js ecosystem, [Node Version Manager](https://github.com/nvm-sh/nvm) is a good place to start for managing different Node.js installations and environments. 


1. Install dependencies
```BASH
yarn install
```

2. Build plugin in development mode or run in watch mode
```BASH
yarn dev
```
or
```BASH
yarn watch
```
3. Build plugin in production mode
```BASH
yarn build
```

### Back-end Development

1. Update [Grafana plugin SDK for Go](https://grafana.com/docs/grafana/latest/developers/plugins/backend/grafana-plugin-sdk-for-go/) dependency to the latest minor version:

```bash
go get -u github.com/grafana/grafana-plugin-sdk-go
```

2. Build backend plugin binaries for Linux, Windows and Darwin:
```BASH
mage -v
```

3. List all available Mage targets for additional commands:
```BASH
mage -l
```
