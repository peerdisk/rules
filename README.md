# rules

The PeerDisk backup service's ruleset parser.

## Syntax

```bash
# this is a comment
# the shortest interval should be a common denominator
interval 6h keep 1w
interval 2d keep 1m
interval 7d keep 1y
```

## Interval duration parsing


Duration parsing is very flexible.

A duration token consists of a series of coefficient/unit tuples.

The coefficient must be positive. It may be decimal.

The following units are available: 

| Abbreviation | Description |
|--------------|-------------|
| s            | 1 second    |
| m            | 60 seconds  |
| h            | 60 minutes  |
| d            | 24 hours    |
| w            | 7 days      |
| mo           | 30 days     |
| y            | 12 months   |


The following durations are valid:

- 4h
- 4.5h4m
- 1y5.5mo