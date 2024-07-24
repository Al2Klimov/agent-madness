#!/bin/bash
set -exo pipefail

while ! [ -e /ca/var/lib/icinga2/certs/ca.crt ]
do sleep 1
done

test -e /var/lib/icinga2/ca || ln -vs {/ca,}/var/lib/icinga2/ca
exec "$@"
