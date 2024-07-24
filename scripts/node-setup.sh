#!/bin/bash
set -exo pipefail

test -e /var/lib/icinga2/certs/ca.crt || icinga2 node setup --master --disable-confd --accept-{config,commands}
exec "$@"
