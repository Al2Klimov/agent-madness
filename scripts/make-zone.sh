#!/bin/bash
set -exo pipefail

curl -fsSLX POST "http://mkzones:8080/v1?name=$(hostname)"
exec "$@"
