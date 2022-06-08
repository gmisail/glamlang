open Glamlang.Lexer

let () =
  let source = {| hello 123 bazinga_babyyyy 123.45 h2o "test" "longggerrrrr stringgsssss" 
    "multi
    line
    string"
  let for while if else { } [ ] ( ) 
  + - * / = -> => >= <= ==
  |} in
  let context = Lexer.new_context source in 
  let tokens = Lexer.scan_tokens context in
  Lexer.print_tokens tokens