#!/usr/bin/env bash

THISDIR=$(dirname $0)

cat <<EOM | gofmt > $THISDIR/bundled_asc_inject.go
// Code generated during build process, along with asc.json. DO NOT EDIT.
package contracts
var ascJson []byte
func init() {
        ascJson = []byte{
        $(cat $THISDIR/asc-v1_1.json | hexdump -v -e '1/1 "0x%02X, "' | fmt)
        }
}
EOM
