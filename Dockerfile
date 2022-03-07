
FROM golang:1.17-buster AS build

# Update aptitude with new repo
RUN apt-get update

WORKDIR /
RUN apt-get install -y git

RUN git clone https://github.com/bankole7782/sites115.git

WORKDIR /sites115
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/sites115d ./sites115d
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/sites115 .
RUN chmod +x bin/sites115d
RUN chmod +x bin/sites115
COPY . src
RUN /sites115/bin/sites115 rso src


FROM gcr.io/distroless/base-debian11
WORKDIR /
COPY --from=build /sites115/bin/sites115d sites115d
COPY --from=build /sites115/src/out site
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/sites115d", "site"]
