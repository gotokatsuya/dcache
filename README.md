# dcache

Disk cache for Golang.

## Algorithm

Like LRU.

Prefer latest & reference count.

### System

One cache that can store 3 key.

Current state.
```
[1]-[2]-[3]
```

Get [2].

Then,
```
[1]-[3]-[2]
```

Set [4]

Then,
```
[3]-[2]-[4]
```

Set [5]

Then,
```
[2]-[4]-[5]
```

Get [5]

Then,
```
[2]-[4]-[5]
```
