# Code generation

At the moment, the Glam compiler will generate C code. This document will provide some examples
of Glam concepts in C.

## Variable Assignment

### Glam representation
```
let x : int = 100
let name : string = "graham"
let sum : int -> int = fn (x: int, y: int): int => x + y
```
### C representation
```
int x = 100;
GlamString name = "graham";
int sum(int x, int y) {
  return x + y;
}
```

## Records

### Glam representation
```
type Food {
    name: string,
    calories: int
}

type Scone(Food) {
  flavor: string
}
```

### C representation
```c
typedef struct Food;
typedef struct Scone;

/* all other forward declarations */

struct Food {
  GlamString name;
  GlamInt calories;
};

struct Scone {
  GlamString name;
  GlamInt calories;
  GlamString flavor; 
};
```

