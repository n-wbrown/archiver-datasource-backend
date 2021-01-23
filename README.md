# Grafana Data Source Backend Plugin for Epics Archiver


## Getting started

### Installation 

1. Download a pre-built release from the releases page.

2. Unzip the contents in the Grafana plugins folder.

3. This package is unsigned requiring it to be specially listed in the Grafana `configur.ini` file in order to be run. Add `https://github.com/n-wbrown/archiver-datasource-backend` to the `allow_loading_unsigned_plugins` parameter in the `[plugins]` section.

4. Restart Grafana.

### Configuring the plugin

1. Log into Grafana's web interface as an administrator. Go to the "Configuration / Data Sources" page and click the "Add data source" button. Locate `archiver-datasource-backend` near the bottom of the list. 

2. To configure the plugin, add the url of the archiver. This will likely take the form `https://[hostname]/retrieval` where `hostname` will be unique to your case. 

## Getting started for developers

A data source backend plugin consists of both frontend and backend components.

### Frontend

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

### Backend

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
