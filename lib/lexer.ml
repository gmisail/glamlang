module Lexer = struct
include Token

  type context = {
    source: string;
    start: int;
    line: int;
    current: int;
    tokens: Token.token list;
  }

  let is_at_end context = 
    context.current >= String.length context.source

  let advance context =
    { context with current = context.current + 1 }

  let advance_line context = 
    { context with line = context.line + 1 }

  let advance_line_if_newline next_char context = 
    if next_char = Some '\n' then advance_line context else context
    
  let current_char context =
    try
      Some (String.get context.source (context.current - 1))
    with
      Invalid_argument _ -> None

  let peek_char context =
    try 
      Some (String.get context.source context.current)
  with
    Invalid_argument _ -> None

  let expect_char context expected_char =
    if is_at_end context then (false, context) else
    match peek_char context with
      | Some character when character = expected_char -> (true, advance context)
      | _ -> (false, context)

  let add_token token context =
    { context with tokens = context.tokens @ [token] }

  (* adds a token that contains a literal, i.e. string, integer, etc... *)
  let add_literal token_type literal_type context =
    let text = (String.sub context.source context.start (context.current - context.start)) in
    add_token (Token.create token_type text literal_type context.line) context
    
  (* adds a token that does not contain a literal, i.e. "let" *)
  let add_keyword token_type context = 
    add_literal token_type None context

  let add_conditional_token context expected_char token token_fallback =
    match expect_char context expected_char with
      | (true, updated_context) -> add_keyword token updated_context
      | (false, updated_context) -> add_keyword token_fallback updated_context

  let rec add_string context =
    let next_char = peek_char context in 
    if not (next_char = Some '"') && not (is_at_end context) then
      context |> advance_line_if_newline next_char |> advance |> add_string
    else
      let string_value = String.sub context.source (context.start + 1) (context.current - 1) in
      add_literal STRING (Some (STRING_LITERAL string_value)) context

  let scan_token context =
    let advanced_context = advance context in
    let current_char = current_char advanced_context in
    match current_char with
      | Some '\r' | Some '\t' | Some ' ' -> advanced_context  (* ignore whitespace *)
      | Some '\n' -> advance_line advanced_context 
      | Some '(' -> add_keyword L_PAREN advanced_context
      | Some ')' -> add_keyword R_PAREN advanced_context
      | Some '{' -> add_keyword L_BRACKET advanced_context
      | Some '}' -> add_keyword R_BRACKET advanced_context
      | Some '[' -> add_keyword L_BRACE advanced_context
      | Some ']' -> add_keyword R_BRACE advanced_context
      | Some ',' -> add_keyword COMMA advanced_context
      | Some '.' -> add_keyword PERIOD advanced_context
      | Some '+' -> add_keyword ADD advanced_context
      | Some '-' -> add_keyword SUB advanced_context
      | Some '*' -> add_keyword MULT advanced_context
      | Some '/'  -> add_keyword DIV advanced_context
      | Some '=' -> add_conditional_token context '=' EQUALITY EQUAL
      | Some '>' -> add_conditional_token context '=' GT_EQ GT
      | Some '<' -> add_conditional_token context '=' LT_EQ LT
      | Some '!' -> add_conditional_token context '=' NOT_EQUAL EQUAL
      | Some '"' -> add_string { advanced_context with start = advanced_context.current }
      | _ -> failwith "Disallowed character"

  let scan_tokens input = []
end
