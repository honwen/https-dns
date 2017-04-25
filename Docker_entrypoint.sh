#!/bin/sh

cmdArgs="$*"
if [ -n "$cmdArgs" ]; then
  /opt/google-https-dns --logtostderr -V 4 $cmdArgs
  exit 0
fi

Args=${Args:--T -U -d 8.8.8.8}

cat > /opt/supervisord.conf <<EOF
[supervisord]
nodaemon=true

[program:google-https-dns]
command=/opt/google-https-dns ${Args} --logtostderr -V 3
autorestart=true
redirect_stderr=true
stdout_logfile=/dev/stdout
stdout_logfile_maxbytes=0

EOF

/usr/bin/supervisord -c /opt/supervisord.conf
