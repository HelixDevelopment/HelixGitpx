FROM golang:1.23.4-alpine AS build
WORKDIR /src
COPY . .
RUN GOTOOLCHAIN=go1.23.4 go build -o /out/billing-service ./cmd/billing-service

FROM gcr.io/distroless/static:nonroot
COPY --from=build /out/billing-service /billing-service
USER nonroot:nonroot
ENTRYPOINT ["/billing-service"]
