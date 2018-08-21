# Keyset

A Sorted Set for storing 64bit keys (such as record ids) in a compact `[]byte` slice. 100% test coverage.

This was created to handle boltdb/goleveldb indexes where you want to add/remove record ids quickly and with minimum storage requirements.

This is a more limited, but faster version of the [bolthold keyList](https://github.com/timshannon/bolthold/blob/master/index.go#L108) which supports storing any type using gob encoding.
