module Lexer = struct
  include Token

  type context = {
    source: string;
    start: int;
    line: int;
    current: int;
  }

  let null_term = '\x00'

  let is_at_end context = 
    context.current >= String.length context.source

  let advance context =
    { context with current = context.current + 1 }

  let current_char context =
    try
      Some (String.get context.source (context.current - 1))
    with
      Invalid_argument _ -> None

  (* adds a token that contains a literal, i.e. string, integer, etc... *)
  let add_token token_type literal_type context =
    let text = (String.sub context.source context.start (context.current - context.start)) in
    Token.create token_type text literal_type context.line context
    
  (* adds a token that does not contain a literal, i.e. "let" *)
  let add_keyword token_type context = 
    add_token token_type None context

  let scan_token context =
    let advanced_context = advance context in
    match (current_char, advanced_context) with
      | Some "(" -> add_keyword L_PAREN advanced_context
      | Some ")" -> add_keyword R_PAREN advanced_context
      | Some "{" -> add_keyword L_BRACKET advanced_context
      | Some "}" -> add_keyword R_BRACKET advanced_context
      | Some "[" -> add_keyword L_BRACE advanced_context
      | Some "]" -> add_keyword R_BRACE advanced_context
      | Some "," -> add_keyword COMMA advanced_context
      | Some "." -> add_keyword PERIOD advanced_context
      | Some "+" -> add_keyword ADD advanced_context
      | Some "-" -> add_keyword SUB advanced_context
      | Some "*" -> add_keyword MULT advanced_context
      | Some "/"  -> add_keyword DIV advanced_context
      | _ -> failwith "Disallowed character"

  let scan_tokens input = [ null_term ]
  
  (* create and execute the lexer on a given string *)
  let create_lexer (input : string) =
    let lexer =
      {
        source = input;
        start = 0;
        current = -1;
      }
    in advance lexer
end
