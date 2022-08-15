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

let sum : (int, int) -> int = fun(x: int, y: int): int => x + y

sum(a, b)   # 4250
```
