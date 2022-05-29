(* use this as reference for this: https://craftinginterpreters.com/appendix-i.html, also the github for menhir lua *)

type name = Name of string

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

type unary_op = 
        | LogicalNot
        | Increment
        | Decrement

type expression = 
        | None
        | False
        | True
        | Integer of int
        | Float of float
        | String of string
        | BinOp of binary_op * expression * expression 
        | UnOp of unary_op * expression

type arguments = Arguments of expression list

type block = statement list

type statement = 
        | Assignment of name * expression
        | FunctionCall of name * arguments 
        | While of expression * block
        | If of (expression * block) list * block option
        | For of name * expression              (* for var in expr { ... } *) 

