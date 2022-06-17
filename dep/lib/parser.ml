open Token
open Ast

module Parser = struct
  type context = { current : int; tokens : Token.token list }

  let advance parser = { parser with current = parser.current + 1 }
  let current parser = List.nth_opt parser.tokens parser.current

  let get_string literal =
    match literal with
    | Some (Token.STRING_LITERAL string_val) -> string_val
    | _ -> "NULL"

  let get_number literal =
    match literal with
    | Some (Token.NUMBER_LITERAL num_val) -> num_val
    | _ -> Float.max_float

  let get_identifier literal =
    match literal with Some (Token.NAME_LITERAL name) -> name | _ -> "NULL"

  let expect_token parser token_type =
    let next_token = current parser in      
    match next_token.kind with
    | token_type -> (Some next_token, advance parser) 
    | _ -> raise (Failure "Expected token <insert token here>.")

  let parse_let_statement parser = 
    let (_, parser') = expect_token parser Token.LET in   
    let (id, parser'') = 
        try expect_token parser' Token.NAME
        with raise (Failure "Expected 'let' token.") 
    in
    
    


  (* returns a statement, if one exists. Returns None if the file is over. *)
  let parse_statement parser =
    let current_token = current parser in
    match current_token.kind with
    | Token.LET -> (parse_let_statement (advance parser))
    | _ -> raise (Failure "Could not parse statement.")


  let primary parser =
    let token = current parser in
    let token_val =
      if Option.is_some token then Option.get token
      else raise (Failure "Could not parse token.")
    in
    match token_val.kind with
    | Token.TRUE -> (Ast.Boolean true, advance parser)
    | Token.FALSE -> (Ast.Boolean false, advance parser)
    | Token.STRING -> (Ast.String , advance parser)
    | Token.NUMBER -> (Ast.Float , advance parser)
    | _ -> failwith "Could not parse primitive, invalid type."

  (* parses statements until there are none left, returns list of statements *)
  let rec parse_program parser =
    let statement, advanced_parser = parse_statement parser in
    match statement with
    | Some _ -> [ statement ] @ parse_program advanced_parser
    | None -> []

  (* returns list of statements based on a list of tokens *)
  let start tokens =
    let parser = { current = 0; tokens } in
    parse_program parser
end
