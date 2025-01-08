ARG BUILD_IMAGE=go-build

FROM $BUILD_IMAGE AS build
ARG SERVICE_NAME
WORKDIR /app
COPY . .
RUN GOOS=linux go build -o ${SERVICE_NAME} -a ./microservices/${SERVICE_NAME}/cmd/app

FROM alpine:latest
ARG SERVICE_NAME
WORKDIR /app
COPY --from=build /app/$SERVICE_NAME .
COPY deploy/config.yml .
COPY ./microservices/${SERVICE_NAME}/docs docs

ENV SERVICE_NAME=${SERVICE_NAME}
ENTRYPOINT ["sh", "-c", "./$SERVICE_NAME"]
