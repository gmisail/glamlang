# ✨ glamlang

Glam is a experimental programming language inspired by Javascript, Haskell, and OCaml. 

## Build
Glam requires `task` for running project commands. To build, simply run:
```
task build
```

## Tests

```
go test -v
```

## Example

```
let a : int = 150
let b : int = 4100

let sum : (int, int) -> int = fn(x: int, y: int): int => x + y

sum(a, b)   # 4250

type NumberPair {
    first: int
    second: int
}

let pair: NumberPair = NumberPair { first: 10, second: 20 }
sum(pair.a, pair.b)   # 30

let mul: (NumberPair) -> int = fn(pair: NumberPair): int => {
    return pair.first * pair.second
}

mul(pair)   # 300
```
