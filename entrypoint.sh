#!/bin/sh

echo Starting server...

if [ "$DEBUG" == "true"  ]; then
  echo "Run into debug mode.. you can exec to pod to check image status..."
  exec sleep 50000
else
  echo "Executing container cmd..."
  exec "$@"
fi
