#!/bin/bash
set -exo pipefail

cp -r /config/* /etc/icinga2/
exec "$@"
