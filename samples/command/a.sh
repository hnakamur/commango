#!/bin/sh
echo hello, stdout.
sleep 1
echo hi, stderr 1>&2
sleep 1
echo hello, stdout 2.
echo hello, stdout 3.
sleep 1
echo hello, stdout 4.
sleep 1
echo hi, stderr 2 1>&2
sleep 1
echo hello, stdout 5.
sleep 1
echo hello, stdout 6.
exit 1
