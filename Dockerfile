FROM golang:1.13-alpine AS build
RUN apk add binutils && rm -r /var/cache/
WORKDIR /src
ENV CGO_ENABLED=0
COPY . .
RUN go build -o /out/devproxy2 && strip /out/devproxy2

FROM scratch AS bin
COPY --from=build /out/devproxy2 /bin/devproxy2
COPY ./docker.toml /etc/devproxy.toml
EXPOSE 8111/tcp
ENTRYPOINT ["/bin/devproxy2"]
