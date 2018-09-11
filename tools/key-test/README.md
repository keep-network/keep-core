Key Test
===============

Code copied from github.com/ethereum/go-ethereum/cmd/ethkey
 with lots of additions and fixes.

Fixed command line so that options that are take are actually used.

Fixed so that this is a complete set of operations on keys.

Fixed so that error messages reflect what went wrong.

Added check that passwords are not in list of 500 million known to be pwned
passwords.

Fixed so that you can update a password for a key file.


