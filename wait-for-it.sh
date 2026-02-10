#!/usr/bin/env bash
# wait-for-it.sh — espera por host:port e então executa o comando

set -e

if [ $# -lt 1 ]; then
  echo "Uso: $0 host:port [-- comando args...]"
  exit 1
fi

hostport="$1"
shift

# separa host e porta
host="${hostport%:*}"
port="${hostport##*:}"

until nc -z "$host" "$port"; do
  >&2 echo "⏳ $host:$port ainda não está disponível — aguardando..."
  sleep 1
done

>&2 echo "✅ $host:$port está no ar — executando comando"
exec "$@"
