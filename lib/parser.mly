%{ 
    open Ast
%}

(* symbols *)
%token EQUAL
%token COLON
%token COMMA
%token DOT

%token L_PAREN
%token R_PAREN
%token L_BRACK
%token R_BRACK
%token L_BRACE
%token R_BRACE

%token EOF

%token PLUS
%token MINUS
%token ASTERISK
%token SLASH
%token EXCLAMATION

%token DOUBLE_EQ
%token GEQ
%token LEQ
%token GT
%token LT
%token AND
%token OR

%token THIN_ARROW   // ->
%token THICK_ARROW  // =>

(* keywords *)
%token LET

%token NONE

%token BREAK

%token IF
%token ELSE

%token FALSE
%token TRUE

%token FOR
%token RETURN
%token WHILE

%token <int> INTEGER
%token <float> FLOAT
%token <Ast.name> NAME
%token <string> LITERAL_STRING