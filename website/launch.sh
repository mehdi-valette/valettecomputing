#/usr/bin/sh

cd src && /usr/local/go/bin/go run ./cmd/valettesoftware.go --port 8080 &

while inotifywait -r -e modify .; do
  PID=$(pidof valettesoftware)
  echo "killing $PID"
  kill $PID
  echo "running again $(date -Is)"
  cd src && /usr/local/go/bin/go run ./cmd/valettesoftware.go --port 8080 &
done