#/usr/bin/sh

/usr/local/go/bin/go run . &

while inotifywait -e modify .; do
  PID=$(pidof valettecomputing.ch)
  echo "killing $PID"
  kill $PID
  echo "running again $(date -Is)"
  /usr/local/go/bin/go run . &
done