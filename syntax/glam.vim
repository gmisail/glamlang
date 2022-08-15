" Glam Programming Language syntax
" Language: Glam 
" Maintainer: Graham Misail

:if exists("b:current_syntax")
:  finish
:endif

" keywords
:syntax keyword glamKeyword let
:syntax keyword glamKeyword and or 
:syntax keyword glamKeyword struct while if return
:syntax keyword glamKeyword fun module

" booleans
:syntax keyword glamBoolean true false

" operators
:syntax match glamOperator "\v\*"
:syntax match glamOperator "\v\+"
:syntax match glamOperator "\v\-"
:syntax match glamOperator "\v/"
:syntax match glamOperator "\v\="
:syntax match glamOperator "\v!"

" conditionals
:syntax keyword glamConditional if else and or else

" numbers
:syntax match glamNumber "\v\-?\d*(\.\d+)?"

" strings
:syntax region glamString start="\v\"" end="\v\""

" comments
" :syntax match glamComment "\v//.*$"

:highlight link glamKeyword Keyword
:highlight link glamBoolean Boolean
:highlight link glamFunction Function
:highlight link glamOperator Operator
:highlight link glamConditional Conditional
:highlight link glamNumber Number
:highlight link glamString String
":highlight link glamComment Comment

:let b:current_syntax = "gl"

