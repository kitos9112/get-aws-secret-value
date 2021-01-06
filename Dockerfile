FROM golang:1.15.6-alpine AS build
RUN apk add --no-cache make jq
WORKDIR /go/src/github.com/kitos9112/get-aws-secret-value.git/
COPY . .
ENV CGO_ENABLED=0
RUN make binaries/linux_x86_64/aws-get-secret-value && mv binaries/linux_x86_64/aws-get-secret-value /app

FROM alpine:3.12
RUN apk add --no-cache ca-certificates
COPY --from=build /app /bin/aws-get-secret-value
CMD [ "aws-get-secret-value" ]
