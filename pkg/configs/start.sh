#!/usr/bin/env bash
pushd "${BASH_SOURCE%/*}"

# This is the command that "ggt run <node>" uses
#
# Disables staking
# Validates all subnets that have plugins installed
# Allow connections from anywhere 
# Disable NAT 
# Dont try to connect to anyone bootstrap nodes

cmd="bin/avalanchego \
	--http-host=0.0.0.0 \
	--public-ip=127.0.0.1 \
	--bootstrap-ids= \
	--bootstrap-ips= \
	--sybil-protection-enabled=false \
  --index-enabled=true \
  --api-keystore-enabled=true \
  --api-admin-enabled=true \
  --log-rotater-max-files=1 \
  --log-rotater-max-size=1 \
  --consensus-shutdown-timeout=1s \
	--data-dir={{.DataDir}} \
	--config-file={{.ConfigFile}} \
	--chain-config-dir={{.ChainConfigDir}} \
	--plugin-dir={{.PluginDir}} \
	--vm-aliases-file={{.VMAliasesFile}} \
  --genesis-file={{.AvaGenesisFile}} \
	--chain-aliases-file={{.ChainAliasesFile}}"

if [[ -n "$VERBOSE" ]]; then
  exec $cmd "$@"
else 
  echo "Node running with stdout suppressed. See logs in data/logs"
  exec $cmd "$@" > /dev/null
fi
