#!/bin/bash

# Extract the address
grep GenRequestID: $1 | awk '{print $2}' > ,addr1

# Add teh checksum to the address
../bin/cs.js $( cat ,addr1 ) >,addr2

# put in a var so we can substitute it in the m4 template below
addr=$(cat ,addr2)

# generate the .m4 file
cat >addr.m4 <<XXxx
m4_changequote([[[,]]])m4_dnl
m4_define([[[m4_comment]]],[[[]]])m4_dnl
m4_comment([[[
	@title This is the address of the RequestID generation contract.   This file is generated in ./gen-addr.sh.  Do not edit.
	@author Philip Schlump
]]])m4_dnl
m4_define([[[GEN_REQUEST_ID_ADDR]]],[[[$addr]]])m4_dnl
XXxx


