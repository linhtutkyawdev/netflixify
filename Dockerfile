# syntax=docker/dockerfile:1

##
## Build the application from source
##

FROM golang:latest AS build-stage

WORKDIR /app

COPY . ./

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o /netflixify ./cmd/api/main.go

RUN cd ./doctron && \
    go mod download && \
    CGO_ENABLED=0 GOOS=linux go build -o /doctron

##
## Run the tests in the container
##

FROM build-stage
RUN go test -v ./tests

##
## Deploy the application binary into a lean image
##


FROM lampnick/runtime:chromium-alpine

WORKDIR /

COPY --from=build-stage /netflixify /netflixify
COPY --from=build-stage /doctron /doctron
COPY --from=build-stage /app/doctron/conf/default.yaml /doctron.yaml
COPY --from=build-stage /app/run.sh /run.sh
RUN chmod +x /run.sh

ENV PORT 3000
EXPOSE $PORT

EXPOSE 8080

CMD [ "./run.sh" ]