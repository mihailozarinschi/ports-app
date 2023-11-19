FROM golang:1.20-alpine AS build
WORKDIR /go/src/portsd
COPY . .
RUN CGO_ENABLED=0 go build -o /go/bin/portsd ./cmd/portsd

FROM scratch
COPY --from=build /go/bin/portsd /bin/portsd
ENTRYPOINT ["/bin/portsd"]
