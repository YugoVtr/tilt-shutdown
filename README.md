# Graceful Shutdown and TILT

## [TILT](https://tilt.dev/)
```bash
$ tilt up
```
Open `0.0.0.0:8000` on your browser

## Graceful Shutdown
In a graceful shutdown, the server stops accepting new requests while continuing to process the already received ones. When a server receives a signal to shut down gracefully, it performs the following steps:

1. The server stops accepting new connections.
2. The server waits for all the ongoing requests to finish.
3. If there are long-running requests that exceed a certain timeout, the server may choose to forcibly close those connections.
4. Once all active requests are completed, the server shuts down.
