# syntax=docker/dockerfile:1

##
## Build the application from source
##

FROM golang:latest AS build

WORKDIR /app

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o ./netflixify ./cmd/api/main.go

#
# Run the tests in the container
#

FROM build
RUN go test -v ./tests

##
## Deploy the application binary into a lean image
##

# final image
# https://github.com/chromedp/docker-headless-shell#using-as-a-base-image
FROM chromedp/headless-shell:latest

RUN export DEBIAN_FRONTEND=noninteractive \
    && apt-get update \
    && apt-get install -y --no-install-recommends \
    ca-certificates dumb-init fonts-noto fonts-noto-cjk \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/

COPY --from=build /app/netflixify /usr/local/bin

ENV URL https://netflixify.onrender.com
ENV BOT_TOKEN=6836197587:AAHdUrYAAINPDN1V4OhS4aAFDLDTzxE_eqU

# ENV PORT 3000
EXPOSE $PORT

ENTRYPOINT ["dumb-init", "--"]
CMD [ "netflixify" ]