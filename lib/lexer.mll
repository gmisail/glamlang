{
open Lexing
open Ast
open Parser
}

let white = [' ' '\t']+
let digit = ['0'-'9']
let int = '-'? digit+
let letter = ['a'-'z' 'A'-'Z']
let id = letter+

(*

%token L_PAREN
%token R_PAREN
%token L_BRACK
%token R_BRACK
%token L_BRACE
%token R_BRACE

%token EOF

%token <float> FLOAT
%token <string> LITERAL_STRING
*)

rule read =
  parse
  | white { read lexbuf }
  | "true"  { TRUE }
  | "false" { FALSE }
  | "let"   { LET }
  | "none"  { NONE }
  | "if"    { IF }
  | "else"  { ELSE }
  | "break" { BREAK }
  | "for"   { FOR }
  | "return" { RETURN }
  | "while" { WHILE }
  | "="     { EQUAL }
  | ":"     { COLON }
  | ","     { COMMA }
  | "."     { DOT }
  | "("     { L_PAREN }
  | ")"     { R_PAREN }
  | "&&"    { AND }
  | "||"    { OR }
  | "=="    { DOUBLE_EQ }
  | ">="    { GEQ }
  | "<="    { LEQ }
  | ">"     { GT }
  | "<"     { LT }
  | "+"     { PLUS }
  | "-"     { MINUS }
  | "*"     { ASTERISK }
  | "/"     { SLASH }
  | "!"     { EXCLAMATION }
  | "->"    { THIN_ARROW }
  | "=>"    { THICK_ARROW }
  | id      { NAME (Lexing.lexeme lexbuf) }
  | int     { INTEGER (int_of_string (Lexing.lexeme lexbuf)) }
  | eof     { EOF }
