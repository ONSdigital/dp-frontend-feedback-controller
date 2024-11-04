# dp-frontend-feedback-controller

To allow users to provide feedback about the ONS website

## Getting started

- Run `make debug`
- Clone `dp-design-system`
    - Run `npm run dev`

## Dependencies

- dp-design-system - Serves CSS and JS for the page
- No further dependencies other than those defined in `go.mod`

## Configuration

| Environment variable           | Default                         | Description                                                                                                        |
| ------------------------------ | ------------------------------- | ------------------------------------------------------------------------------------------------------------------ |
| API_ROUTER_URL                 | <http://localhost:23200/v1>     | The URL of the [dp-api-router](https://github.com/ONSdigital/dp-api-router)                                        |
| BIND_ADDR                      | localhost:25200                 | The host and port to bind to                                                                                       |
| CENSUS_TOPIC_ID                | 4445                            | The census topic id                                                                                                |
| DEBUG                          | false                           | Enable debug mode                                                                                                  |
| ENABLE_CENSUS_TOPIC_SUBSECTION | false                           | Enable census topic subsection                                                                                     |
| ENABLE_NEW_NAVBAR              | false                           | Enable new navigation bar                                                                                          |
| GRACEFUL_SHUTDOWN_TIMEOUT      | 5s                              | The graceful shutdown timeout in seconds (`time.Duration` format)                                                  |
| HEALTHCHECK_INTERVAL           | 30s                             | Time between self-healthchecks (`time.Duration` format)                                                            |
| HEALTHCHECK_CRITICAL_TIMEOUT   | 90s                             | Time to wait until an unhealthy dependent propagates its state to make this app unhealthy (`time.Duration` format) |
| MAIL_HOST                      | ""                              | The host for the mail server.                                                                                      |
| MAIL_PORT                      | ""                              | The port for the mail server.                                                                                      |
| MAIL_USER                      | ""                              | A user on the mail server.                                                                                         |
| MAIL_PASSWORD                  | ""                              | The password for the mail server user.                                                                             |
| FEEDBACK_TO                    | ""                              | Receiver email address for feedback.                                                                               |
| FEEDBACK_FROM                  | ""                              | Sender email address for feedback.                                                                                 |
| IS_PUBLISHING_MODE             | false                           |                                                                                                                    |
| PATTERN_LIBRARY_ASSETS_PATH    | ""                              | Pattern library location                                                                                           |
| SERVICE_AUTH_TOKEN             | ""                              | Service authorisation token                                                                                        |
| SITE_DOMAIN                    | localhost                       |                                                                                                                    |
| SUPPORTED_LANGUAGES            | []string{"en", "cy"}            | Supported languages                                                                                                |
| OTEL_EXPORTER_OTLP_ENDPOINT    | localhost:4317                  | Endpoint for OpenTelemetry service                                                                                 |
| OTEL_SERVICE_NAME              | dp-frontend-feedback-controller | Label of service for OpenTelemetry service                                                                         |
| OTEL_BATCH_TIMEOUT             | 5s                              | Timeout for OpenTelemetry                                                                                          |
| OTEL_ENABLED                   | false                           | Feature flag to enable OpenTelemetry

## Contributing

See [CONTRIBUTING](CONTRIBUTING.md) for details.

## License

Copyright © 2023, Office for National Statistics (<https://www.ons.gov.uk>)

Released under MIT license, see [LICENSE](LICENSE.md) for details.
