open Glamlang.Lexer

(*
hello 123 bazinga_babyyyy 123.45 h2o "test" "longggerrrrr stringgsssss" 
    "multi
    line
    string"
  let for while if else { } [ ] ( ) 
  + - * / = -> => >= <= ==   
*)

let () =
  let source = {| let x: int = 100 let square: ( int ) -> int = ( y : int ) => { return y * y } |} in
  let context = Lexer.new_context source in 
  let tokens = Lexer.scan_tokens context in
  Lexer.print_tokens tokens