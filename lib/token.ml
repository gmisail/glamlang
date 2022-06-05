module Token = struct
  type token_type =
    (* misc *)
    | ILLEGAL
    | END_OF_FILE

    (* types *)
    | NAME
    | INTEGER
    | FLOAT
    | BOOL

    (* math *)
    | EQUAL
    | ADD
    | SUB
    | MULT
    | DIV
    | GT
    | GT_EQ
    | LT
    | LT_EQ
    | BANG
    | EQUALITY
    | NOT_EQUAL
    | COMMA
    | PERIOD
    | COLON
    | L_PAREN
    | R_PAREN
    | L_BRACE
    | R_BRACE
    | L_BRACKET
    | R_BRACKET
    | QUOTE
    | STRING

    (* registered keywords *)
    | LET 
    | WHILE
    | FOR
    | IF
    | ELSE

  type literal_type =
    | STRING_LITERAL of string
    | INTEGER_LITERAL of int

  type token = { 
    literal : literal_type option;
    lexeme : string; 
    kind : token_type;
    line : int;
  }

  let type_to_string kind =
    match kind with
    | ILLEGAL -> "ILLEGAL"
    | END_OF_FILE -> "EOF"

    (* types *)
    | NAME -> "NAME"
    | INTEGER -> "INTEGER"
    | FLOAT -> "FLOAT"
    | BOOL -> "BOOL"

    (* math *)
    | EQUAL -> "EQUALS"
    | ADD -> "ADD"
    | SUB -> "SUB"
    | MULT -> "MULT"
    | DIV -> "DIV"
    | GT -> "GREATER_THAN"
    | GT_EQ -> "GREATER_THAN_EQ"
    | LT -> "LESS_THAN"
    | LT_EQ -> "LESS_THAN_EQ"
    | BANG -> "BANG"
    | EQUALITY -> "DEQUAL"
    | NOT_EQUAL -> "NOT_EQUAL"
    | PERIOD -> "PERIOD"
    | COMMA -> "COMMA"
    | COLON -> "COLON"
    | L_PAREN -> "L_PAREN"
    | R_PAREN -> "R_PAREN"
    | L_BRACE -> "L_BRACE"
    | R_BRACE -> "R_BRACE"
    | L_BRACKET -> "L_BRACKET"
    | R_BRACKET -> "R_BRACKET"
    | QUOTE -> "QUOTE"
    | STRING -> "STRING"

    (* registered keywords *)
    | LET -> "LET"
    | WHILE -> "WHILE"
    | FOR -> "FOR"
    | IF -> "IF"
    | ELSE -> "ELSE"

  let create kind lexeme literal line = { 
      literal; kind; 
      lexeme; line;
    }
end
