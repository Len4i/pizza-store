### Logger
Using 2 loggers:
1. standard slog in json format for the app logs
2. slog based logger for http requests with middleware that enriches log with http request data 