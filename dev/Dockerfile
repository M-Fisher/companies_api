# CI
FROM golang:1.19 AS ci
ENV PROJECTDIR /usr/local/go/src/M-Fisher/companies-service/companies-service
WORKDIR ${PROJECTDIR}
COPY go.mod go.sum ${PROJECTDIR}/
RUN  go mod download
COPY . ${PROJECTDIR}/
RUN make setup

# Binary build
FROM ci AS build
RUN make build

# Production image
FROM scratch AS bin
COPY --from=build /usr/local/go/src/M-Fisher/companies-service/companies-service/output /bin/output
ENTRYPOINT ["/bin/output"]