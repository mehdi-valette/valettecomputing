#/usr/bin/sh

/usr/local/go/bin/go run . --port 8080 &

while inotifywait -r -e modify .; do
  PID=$(pidof valette.software)
  echo "killing $PID"
  kill $PID
  echo "running again $(date -Is)"
  /usr/local/go/bin/go run . --port 8080 &
done