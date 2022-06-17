type name = Name of string

type primitive =
  | Null
  | String of string
  | Int of int
  | Float of float
  | Boolean of bool

type variable_type =
  | VariableType of name
  | Signature of name list * variable_type (* (int, int) -> int *)

type binary_op =
  | Add
  | Sub
  | Mul
  | Div
  | Mod
  | And
  | Or
  | Equality
  | Inequality
  | Less
  | Greater
  | LessEqual
  | GreaterEqual

type unary_op = LogicalNot | Increment | Decrement | Negative

type expression =
  | Constant of primitive
  | Identifier of name
  | BinOp of binary_op * expression * expression
  | UnOp of unary_op * expression
  | FunctionCall of name * arguments
  | FunctionDeclaration of (name * name) list * block
  | Index of expression * expression

and arguments = expression list
and block = statement list

and statement =
  | Assignment of name * expression
  | SignatureDefinition of
      name * variable_type (* define square: (int) -> int *)
  | VariableDeclaration of name * variable_type * expression option
  | While of expression * statement
  | If of (expression * block) list * block option
  | For of name * expression (* for var in expr { ... } *)
  | Return of expression
  | Block of block
  | Eval of expression
