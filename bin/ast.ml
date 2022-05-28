type expr = 
        | Atom of string                        (* hello, world   *)
        | Float of float                        (* 100.0          *)
        | Int of int                            (* 100            *)
        | Bool of bool                          (* true, false    *)
        | Assignment of string * expr           (* a = b          *)
        | FunctionCall of string * expr list    (* function(a, b) *)

type statement = 
        | Expr of expr
        | VariableDeclaration of string * string * expr         (* let x: int = 100 *)
        | VariableDeclarationExpression of string * string      (* let x: int       *)

type program = Program of atom list 
