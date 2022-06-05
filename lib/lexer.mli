open Token

module Lexer : sig 
  type context

  val is_at_end : context -> bool
  val advance : context -> context

  val peek_char : context -> char option
  val current_char : context -> char option
  val expect_char : context -> char -> bool * context 
  
  val add_literal : Token.token_type -> Token.literal_type option -> context -> context 
  val add_keyword : Token.token_type -> context -> context 
  
  val scan_token : context -> context 
  val scan_tokens : 'a -> Token.token list
end