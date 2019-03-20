ARG GOLANG_VERSION
FROM golang:${GOLANG_VERSION} as build

WORKDIR /opt/statsdebug

COPY . /opt/statsdebug

RUN go build --ldflags '-extldflags "-static"' -tags netgo -installsuffix netgo .

FROM scratch

COPY --from=build /opt/statsdebug/statsdebug /

ENTRYPOINT ["/statsdebug"]
