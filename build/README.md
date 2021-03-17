## Getting started

A data source backend plugin consists of both frontend and backend components.

## Use Dockerfile

```BASH
docker build . -t blackcowmoo/grafana-ga-ds
docker run --rm -p 3000:3000  --name=blackcowmoo-grafana-ga-ds blackcowmoo/grafana-ga-ds
```

### Apply Change

```BASH
docker build . -t blackcowmoo/grafana-ga-ds
docker run --rm -p 3000:3000  --name=blackcowmoo-grafana-ga-ds blackcowmoo/grafana-ga-ds
```

## Local Build

### Frontend

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

### Run

```BASH
docker run --rm -p 3000:3000 -v "$(pwd)":/var/lib/grafana/plugins --name=grafana -e "GF_PLUGINS_ALLOW_LOADING_UNSIGNED_PLUGINS=google-analytics" grafana/grafana
```

### Apply Change

```BASH
docker restart grafana
```
