module Lexer = struct
  include Token

  type context = {
    source: string;
    start: int;
    line: int;
    current: int;
    tokens: Token.token list;
  }

  let new_context input = 
    {
      source = input;
      start = -1;
      line = 0;
      current = 0;
      tokens = [];
    }

  let is_at_end context = 
    context.current >= String.length context.source
  
  let is_digit character =
    match character with
      | Some c when c >= '0' && c <= '9' -> true 
      | _ -> false

  let is_alpha character = match character with
    | Some c when c >= 'a' && c <='z' -> true
    | Some c when c >= 'A' && c <='Z' -> true
    | Some c when c = '_' -> true
    | _ -> false

  let is_alpha_num character = is_alpha character || is_digit character

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

  let peek_next_char context =
    try
      Some (String.get context.source (context.current + 1))
    with
      Invalid_argument _ -> None

  let expect_char context expected_char =
    if is_at_end context then (false, context) else
    match peek_char context with
      | Some character when character = expected_char -> (true, advance context)
      | _ -> (false, context)
    
  let is_next_char context expected_char =
    if is_at_end context then false else
    match peek_char context with
      | Some character when character = expected_char -> true
      | _ -> false

  let add_token token context =
    { context with tokens = context.tokens @ [token] }

  (* adds a token that contains a literal, i.e. string, integer, etc... *)
  let add_literal token_type literal_type context =
    let text = (String.sub context.source context.start (context.current - context.start)) in
    add_token (Token.create token_type text literal_type context.line) context
    
  (* adds a token that does not contain a literal, i.e. "let" *)
  let add_keyword token_type context = 
    add_literal token_type None context

  (* check for any of the possible character sequences *)
  let add_conditional_token context possible_options token_fallback =
    try
      let (_, matching_token_type) = 
        List.find (fun (possible_char, _) -> is_next_char context possible_char) possible_options 
      in add_keyword matching_token_type (advance context)
    with
      Not_found -> add_keyword token_fallback context

  let rec add_string context =
    let (is_done, is_error) = match (peek_char context) with 
      | Some '"' -> (true, false)
      | Some _ when not (is_at_end context) -> (false, false)
      | None -> (true, true)
      | _ -> failwith "Fatal."
    in
    if not is_done && not is_error then
      context |> advance_line_if_newline (peek_char context) |> advance |> add_string
    else if is_error then 
      failwith "Unterminated string."
    else
      let advanced_context = advance context in 
      let string_value = 
        String.sub context.source (advanced_context.start + 1) (advanced_context.current - advanced_context.start) 
      in add_literal STRING (Some (STRING_LITERAL string_value)) advanced_context

  let rec scan_integer context =
    if is_digit (peek_char context) then 
      context |> advance |> scan_integer 
    else 
      context
  
  let scan_decimal context =
    let peeked_char = peek_char context in 
    match peeked_char with
    | Some '.' when (is_digit (peek_next_char context)) -> 
        context |> advance |> scan_integer
    | _ -> context

  let add_number context = 
    let advanced_context = context |> scan_integer |> scan_decimal in
    let number_literal = 
      (String.sub advanced_context.source advanced_context.start (advanced_context.current - advanced_context.start)) 
    in let number_value = float_of_string number_literal in 
    add_literal NUMBER (Some (NUMBER_LITERAL number_value)) advanced_context

  let rec add_identifier context = 
    if (is_alpha_num (peek_char context)) then 
      context |> advance |> add_identifier
    else 
    let literal = (String.sub context.source context.start (context.current - context.start)) in
    match literal with
      | "let" -> add_keyword LET context
      | "while" -> add_keyword WHILE context
      | "for" -> add_keyword FOR context
      | "if" -> add_keyword IF context
      | "else" -> add_keyword ELSE context
      | _ -> add_literal NAME (Some (NAME_LITERAL literal)) context

  let scan_token context =
    (* set the start to current before advancing it. *)
    let advanced_context = advance { context with start = context.current } in
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
      | Some '-' -> add_conditional_token advanced_context [('>', ARROW)] SUB
      | Some '*' -> add_keyword MULT advanced_context
      | Some '/'  -> add_keyword DIV advanced_context
      | Some '=' -> add_conditional_token advanced_context [('=', EQUALITY); ('>', THICK_ARROW)] EQUAL
      | Some '>' -> add_conditional_token advanced_context [('=', GT_EQ)] GT
      | Some '<' -> add_conditional_token advanced_context [('=', LT_EQ)] LT
      | Some '!' -> add_conditional_token advanced_context [('=', NOT_EQUAL)] EQUAL
      | Some '"' -> add_string advanced_context
      | c when is_digit c -> add_number advanced_context
      | c when is_alpha c -> add_identifier advanced_context
      | _ -> failwith "Unknown character."

  let rec scan_all_tokens context =
    if not (is_at_end context) then 
      context |> scan_token |> scan_all_tokens 
    else
      context

  (* generates list of tokens from a context *)
  let scan_tokens context =
    let scanned_context = scan_all_tokens context in
    scanned_context.tokens

  let print_tokens tokens = 
    let print_token (token: Token.token) = 
      Printf.printf "[ %s ]\n" (Token.type_to_string token.kind)
    in
    List.iter print_token tokens
end
