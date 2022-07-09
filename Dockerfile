FROM docker.io/library/golang:1.13-alpine AS build
WORKDIR /src
ENV CGO_ENABLED=0
COPY . .
RUN go build -ldflags="-s -w" -trimpath -o /out/devproxy2

FROM scratch AS bin
COPY --from=build /out/devproxy2 /bin/devproxy2
COPY ./docker.toml /etc/devproxy.toml
EXPOSE 8111/tcp
ENTRYPOINT ["/bin/devproxy2"]
