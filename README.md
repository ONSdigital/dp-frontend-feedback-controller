# dp-frontend-feedback-controller

To allow users to provide feedback about the ONS website

## Getting started

* Run `make debug`

## Dependencies

* No further dependencies other than those defined in `go.mod`

## Configuration

| Environment variable           | Default   | Description
| ------------------------------ | --------- | -----------
| BIND_ADDR                      | localhost:25200    | The host and port to bind to
| GRACEFUL_SHUTDOWN_TIMEOUT      | 5s        | The graceful shutdown timeout in seconds (`time.Duration` format)
| HEALTHCHECK_INTERVAL           | 30s       | Time between self-healthchecks (`time.Duration` format)
| HEALTHCHECK_CRITICAL_TIMEOUT   | 90s       | Time to wait until an unhealthy dependent propagates its state to make this app unhealthy (`time.Duration` format)
| RENDERER_URL                   | http://localhost:20010  | The URL of [dp-frontend-renderer](https://www.github.com/ONSdigital/dp-frontend-renderer).
| MAIL_HOST                      | ""                      | The host for the mail server.
| MAIL_PORT                      | ""                      | The port for the mail server.
| MAIL_USER                      | ""                      | A user on the mail server.
| MAIL_PASSWORD                  | ""                      | The password for the mail server user.
| FEEDBACK_TO                    | ""                      | Receiver email address for feedback.
| FEEDBACK_FROM                  | ""                      | Sender email address for feedback.
| DEBUG                          | false                        | Enable debug mode
| API_ROUTER_URL                 | http://localhost:23200/v1    | The URL of the [dp-api-router](https://github.com/ONSdigital/dp-api-router)
| SITE_DOMAIN                    | localhost                    |
| PATTERN_LIBRARY_ASSETS_PATH    | ""                           | Pattern library location
| SUPPORTED_LANGUAGES            | [2]string{"en", "cy"}        | Supported languages

## Contributing

See [CONTRIBUTING](CONTRIBUTING.md) for details.

## License

Copyright © 2021, Office for National Statistics (https://www.ons.gov.uk)

Released under MIT license, see [LICENSE](LICENSE.md) for details.

