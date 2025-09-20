#/usr/bin/sh

/usr/local/go/bin/go run . &

while inotifywait -e modify .; do
  PID=$(pidof valettecomputing.ch)
  echo "killing $PID"
  kill $PID
  echo "running again"
  /usr/local/go/bin/go run . &
  echo "running $(pidof valettecomputing.ch)"
done