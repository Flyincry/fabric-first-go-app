#!/bin/bash
#
# Copyright IBM Corp. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#

# This script is designed to be run in the cli container as the
# second step of the EYFN tutorial. It joins the org3 peers to the
# channel previously setup in the BYFN tutorial and install the
# chaincode as version 2.0 on peer0.org3.
#

CHANNEL_NAME="$1"
DELAY="$2"
TIMEOUT="$3"
VERBOSE="$4"
ORG_NUM="$5"
: ${CHANNEL_NAME:="mychannel"}
: ${DELAY:="3"}
: ${TIMEOUT:="10"}
: ${VERBOSE:="false"}
: ${ORG_NUM:="3"}
COUNTER=1
MAX_RETRY=5

# import environment variables
. scripts/envVar.sh

# joinChannel ORG
joinChannel() {
  ORG=$1
  local rc=1
  local COUNTER=1
  ## Sometimes Join takes time, hence retry
  while [ $rc -ne 0 -a $COUNTER -lt $MAX_RETRY ] ; do
    sleep $DELAY
    set -x
    peer channel join -b $BLOCKFILE >&log.txt
    res=$?
    { set +x; } 2>/dev/null
    let rc=$res
    COUNTER=$(expr $COUNTER + 1)
	done
	cat log.txt
	verifyResult $res "After $MAX_RETRY attempts, peer0.org${ORG} has failed to join channel '$CHANNEL_NAME' "
}

setAnchorPeer() {
  ORG=$1
  scripts/setAnchorPeer.sh $ORG $CHANNEL_NAME
}

setGlobalsCLI $ORG_NUM
BLOCKFILE="${CHANNEL_NAME}.block"

echo "Fetching channel config block from orderer..."
set -x
peer channel fetch 0 $BLOCKFILE -o orderer.example.com:7050 --ordererTLSHostnameOverride orderer.example.com -c $CHANNEL_NAME --tls --cafile "$ORDERER_CA" >&log.txt
res=$?
{ set +x; } 2>/dev/null
cat log.txt
verifyResult $res "Fetching config block from orderer has failed"

infoln "Joining org$ORG_NUM peer to the channel..."
joinChannel $ORG_NUM

infoln "Setting anchor peer for org$ORG_NUM..."
setAnchorPeer $ORG_NUM

successln "Channel '$CHANNEL_NAME' joined"
successln "Org$ORG_NUM peer successfully added to network"