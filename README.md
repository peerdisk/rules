# rules

The PeerDisk backup service's ruleset parser.

## Syntax

A "Hello World" for rules:

```bash
# the shortest interval should be a common denominator for the others
interval 6h keep 1w
interval 2d keep 1m
interval 7d keep 1year
```

## Interval duration parsing


Duration parsing is very flexible.

A duration token consistens of a series of coefficient and unit tuples.

The coefficient must be positive, it may decimal.

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