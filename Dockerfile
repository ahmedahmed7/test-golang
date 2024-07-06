FROM golang:1.22.4 as dev

# Install reflex
RUN go install github.com/cespare/reflex@latest

RUN mkdir -p /app
WORKDIR /app

EXPOSE 5000

HEALTHCHECK --interval=20s --timeout=1m --start-period=20s \
   CMD curl -f --connect-timeout 5 --max-time 10 --retry 5 --retry-delay 0 --retry-max-time 40 --retry-all-errors 'http://localhost:5000/health' || bash -c 'kill -s 15 -1 && (sleep 10; kill -s 9 -1)'

ENTRYPOINT reflex -r '(.go$|go.mod)' --decoration='none' -s -- sh -c 'go run .'
