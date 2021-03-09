FROM golang:1.16-alpine as go_build
WORKDIR /build
COPY . .
RUN go run mage.go -v
RUN ls

FROM node:12-alpine as node_build
WORKDIR /build
COPY package.json yarn.lock ./
RUN yarn --frozen-lockfile
COPY . .
RUN yarn build
RUN ls

FROM grafana/grafana:7.4.3
ENV GF_PLUGINS_ALLOW_LOADING_UNSIGNED_PLUGINS=google-analytics
COPY --from=go_build /build/dist /var/lib/grafana/plugins
COPY --from=node_build /build/dist /var/lib/grafana/plugins
