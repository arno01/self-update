#!/bin/bash
# Ref. https://github.com/gchudnov/docker-tools/blob/master/signals/background/program.sh

[[ ! -z "$VERBOSE" ]] && set -x

pid=0

# SIGUSR1 -handler
my_handler() {
  ls -1 /etc/service/ | xargs -n1 sv 1
}

# SIGTERM -handler

term_handler() {
  # If runsvdir receives a TERM signal, it exits with 0 immediately.
  # If runsvdir receives a HUP signal, it sends a TERM signal to each runsv(8) process it is monitoring and then exits with 111.
  # http://smarden.org/runit/runsvdir.8.html
  #
  # runsvdir isn't waiting for the services to gracefully shutdown on SIGHUP,
  # hence gracefully stop the services and only then send SIGHUP to runsvdir.
  ls -1 /etc/service/ | xargs -n1 sv stop

  if [ $pid -ne 0 ]; then
    kill -SIGHUP "$pid"
    wait "$pid"
  fi
  exit 143; # 128 + 15 -- SIGTERM
}

# SIGINT -handler
int_handler() {
  ls -1 /etc/service/ | xargs -n1 sv interrupt

  # Most of software would terminate on interrupt.
  # However, since there is no good way to really know
  # what software would do, let's assume it wants to exit,
  # giving it a little extra time for graceful shutdown.

  # Apply runit's SVWAIT environment variable override if present.
  sleep ${SVWAIT:-7}

  # runsvdir does not handle SIGINT, so do not relay it over to runsvdir.
  exit 130; # 128 + 2 -- SIGINT
}

# setup handlers
# on callback, kill the last background process, which is `tail -f /dev/null` and execute the specified handler
trap 'kill ${!}; my_handler' SIGUSR1
trap 'kill ${!}; term_handler' SIGTERM
trap 'kill ${!}; int_handler' SIGINT

# run application
runsvdir -P /etc/service &
pid="$!"

# wait forever
while true
do
  tail -f /dev/null & wait ${!}
done
