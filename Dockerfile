FROM golang

WORKDIR /work

COPY go.mod go.sum ./

RUN ["go", "mod", "download"]

COPY . .

RUN ["go", "build", "main.go"]

FROM gcr.io/distroless/base-debian11

COPY --from=0 /work/main /

USER nonroot:nonroot

ENTRYPOINT ["/main"]
