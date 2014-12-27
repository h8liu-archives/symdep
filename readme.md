`symdep` analyzes a dependencies of a go package at file level.
The 

Example:

```
$ symdep lonnie.io # prints the layers
0: inst symbol token
1: sym_table lex_operand lex_string lex_comment
2: sym_scope lex_all
3: builder stmt_lexer
4: parser
5: parse_arg parse_ops parse_sym parse_inst parse_label parse_reg
6: inst_jmp inst_reg inst_sys inst_br inst_imm
7: inst_all
8: parse_stmt
9: parse_func
10: build_func writer parse_all
11: bare_func

0: token
1: lex_comment lex_string lex_operand
2: lex_all
3: stmt_lexer
4: parser
5: parse_label parse_sym parse_reg inst parse_arg
6: inst_jmp inst_br inst_imm parse_inst inst_sys inst_reg symbol
7: inst_all parse_ops sym_table
8: parse_stmt sym_scope
9: parse_func builder
10: writer build_func
11: bare_func parse_all

$ symdep github.com/h8liu/d8/client # prints the dep graph
bug
client
     id_pool
     query
     bug
     exchange
     message
     job
exchange
     client
     query
     message
id_pool
     bug
job
     exchange
     message
     bug
message
     query
     client
query
     client
     message
error: has circular dependency
```