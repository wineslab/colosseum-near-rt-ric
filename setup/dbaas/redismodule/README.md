# Introduction

This subdirectory provides implementation for the commands which are implemented
as a [Redis modules](https://redis.io/topics/modules-intro).

# Compiling and Unit Tests

To compile, run unit tests and install use the commands:
```
./autogen.sh
./configure
make
make test
make install
```

By default unit tests and valgrind memory checking are enabled.
This requires `cpputest` and `valgrind` as additional dependencies.
Unit test memory checking can be disabled with the `configure` option
`--disable-unit-test-memcheck` and the unit tests can be completely disabled
with the `configure` option `--disable-unit-test`.
For example to compile and install with unit tests completely disabled
one would run the commands:
```
./autogen.sh
./configure --disable-unit-test
make
make install
```

# Commands

## SETIE key value oldvalue [expiration EX seconds|PX milliseconds]

Time complexity: O(1) + O(1)

Checks a String 'key' for 'oldvalue' equality and set key for 'value' with
optional expired.

```
Example:

redis> get mykey
(nil)
redis> setie mykey "Hello again" "Hello"
(nil)

redis> set mykey "Hello"
OK
redis> get mykey
"Hello"
redis> setie mykey "Hello again" "Hello"
"OK"
redis> get mykey
"Hello again"
redis> setie mykey "Hello 2" "Hello"
(nil)
redis> get mykey
"Hello again"
redis> setie mykey "Hello 2" "Hello again" ex 100
"OK"
redis> ttl mykey
(integer) 96
redis> get mykey
"Hello 2"
```

## SETNE key value oldvalue [expiration EX seconds|PX milliseconds]

Time complexity: O(1) + O(1)

Checks a String 'key' for 'oldvalue' not equality and set key for 'value' with optional expired.

Example:

```
redis> get mykey
(nil)
redis> setne mykey "Hello again" "Hello"
"OK"
redis> get mykey
"Hello again"
redis> setne mykey "Hello 2" "Hello again"
(nil)
redis> setne mykey "Hello 2" "Hello"
"OK"
redis> get mykey
"Hello 2"
redis> setne mykey "Hello 3" "Hello" ex 100
"OK"
redis> get mykey
"Hello 3"
redis> ttl mykey
(integer) 93
```

## DELIE key oldvalue

Time complexity: O(1) + O(1)

Checks a String 'key' for 'oldvalue' equality and delete the key.

```
Example:
redis> get mykey
(nil)
redis> set mykey "Hello"
"OK"
redis> get mykey
"Hello"
redis> delie mykey "Hello again"
(integer) 0
redis> get mykey
"Hello"
redis> delie mykey "Hello"
(integer) 1
redis> get mykey
(nil)
```

## DELNE key oldvalue

Time complexity: O(1) + O(1)

Checks a String 'key' for 'oldvalue' not equality and delete the key.

```
Example:
redis> get mykey
(nil)
redis> set mykey "Hello"
"OK"
redis> get mykey
"Hello"
redis> delne mykey "Hello"
(integer) 0
redis> get mykey
"Hello"
redis> delne mykey "Hello again"
(integer) 1
redis> get mykey
(nil)
```

## MSETPUB key value [key value...] channel message

Time complexity: O(N) where N is the number of keys to set + O(N+M) where N is the number of clients subscribed to the receiving channel and M is the total number of subscribed patterns (by any client)

Set the given keys to their respective values and post a message to the given channel

## MSETMPUB number_of_key_value_pairs number_of_channel_message_pairs key value [ key value ... ] channel message [ channel message ... ]

Time complexity: O(N) where N is the number of keys to set + O(N_1+M) [ + O(N_2+M) + ... ] where N_i are the number of clients subscribed to the corresponding receiving channel and M is the total number of subscribed patterns (by any client)

Set the given keys to their respective values and post messages to their respective channels

## SETXXPUB key value channel message [channel message...]

Time complexity: O(1) + O(1) + O(N_1+M) [ + O(N_2+M) + ... ] where N_i are the number of clients subscribed to the receiving channel and M is the total number of subscribed patterns (by any client).

Set key to hold string value if key already exists and post given messages to the corresponding channels if key value was set successfully

## SETNXPUB key value channel message [channel message...]

Time complexity: O(1) + O(1) + O(N_1+M) [ + O(N_2+M) + ... ] where N_i are the number of clients subscribed to the receiving channel and M is the total number of subscribed patterns (by any client).

Set key to hold string value if key does not exist and post given messages to the corresponding channels if key value was set successfully

## SETIEPUB key value oldvalue channel message [channel message...]

Time complexity: O(1) + O(1) + O(1) + O(N_1+M) [ + O(N_2+M) + ... ] where N_i are the number of clients subscribed to the receiving channel and M is the total number of subscribed patterns (by any client).

If the string corresponding to 'key' is equal to 'oldvalue' then set key for 'value' and post given messages to the corresponding channels if key value was set successfully

## SETNEPUB key value oldvalue channel message [channel message...]

Time complexity: O(1) + O(1) + O(1) + O(N_1+M) [ + O(N_2+M) + ... ] where N_i are the number of clients subscribed to the receiving channel and M is the total number of subscribed patterns (by any client).

If the string corresponding to 'key' is not equal to 'oldvalue' then set key for 'value' and post given messages to the corresponding channels if key value was set successfully

## DELPUB key [key...] channel message

Time complexity: O(N) where N is the number of keys that will be removed + O(N+M) where N is the number of clients subscribed to the receiving channel and M is the total number of subscribed patterns (by any client)

Removes the specified keys and post a message to the given channel if delete key successfully(return >0)

## DELMPUB number_of_keys number_of_channel_message_pairs key [ key ... ] channel message [ channel message ... ]

Time complexity: O(N) where N is the number of keys that will be removed + O(N_1+M) [ + O(N_2+M) + ... ] where N_i are the number of clients subscribed to the receiving channel and M is the total number of subscribed patterns (by any client)

Remove the specified keys. If any of the keys was deleted succesfully (delete return value > 0) then post given messages to the corresponding channels.

## DELIEPUB key oldvalue channel message [channel message...]

Time complexity: O(1) + O(1) + O(1) + O(N_1+M) [ + O(N_2+M) + ...] where N_i are the number of clients subscribed to the corrensponding receiving channel and M is the total number of subscribed patterns (by any client)

If the string corresponding to 'key' is equal to 'oldvalue' then delete the key. If deletion was succesful (delete return value was 1) then post given messages to the corresponding channels.

## DELNEPUB key oldvalue channel message [channel message...]

Time complexity: O(1) + O(1) + O(1) + O(N_1+M) [ + O(N_2+M) + ...] where N_i are the number of clients subscribed to the corrensponding receiving channel and M is the total number of subscribed patterns (by any client)

If the string corresponding to 'key' is not equal to 'oldvalue' then delete the key. If deletion was succesful (delete return value was 1) then post given messages to the corresponding channels.

## NGET pattern

Time complexity: O(N) with N being the number of keys in the instance + O(N) where N is the number of keys to retrieve

Returns all key-value pairs matching pattern.

```
example:

redis> nget mykey*
(empty list or set)

redis> set mykey1 "myvalue1"
OK
redis> set mykey2 "myvalue2"
OK
redis> set mykey3 "myvalue3"
OK
redis> set mykey4 "myvalue4"
OK
redis> nget mykey*
1) "mykey2"
2) "myvalue2"
3) "mykey1"
4) "myvalue1"
5) "mykey4"
6) "myvalue4"
7) "mykey3"
8) "myvalue3"
```

## NDEL pattern

Time complexity: O(N) with N being the number of keys in the instance + O(N) where N is the number of keys that will be removed

Remove all key-value pairs matching pattern.

```
example:

redis> nget mykey*
1) "mykey2"
2) "myvalue2"
3) "mykey1"
4) "myvalue1"
5) "mykey4"
6) "myvalue4"
7) "mykey3"
8) "myvalue3"

redis> ndel mykey*
(integer) 4

redis> ndel mykey*
(integer) 0
```
