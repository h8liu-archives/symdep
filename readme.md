`symdep` analyzes a dependencies of a go package at file level.
The 

Example:

```
$ symdep lonnie.io/e8vm/asm8 # prints the layers
symdep lonnie.io/e8vm/asm8
0: token
1: lex_comment lex_string lex_operand
2: lex_all
3: stmt_lexer
4: symbol parser
5: sym_table parse_sym inst parse_label parse_reg parse_arg
6: var_stmt sym_scope lib func_stmt inst_sys inst_imm inst_reg inst_jmp parse_inst inst_br
7: var_decl builder func_decl data_str inst_all
8: build_var file build_func parse_var_stmt parse_func_stmt
9: package build_file parse_var parse_func
10: build_lib parse_file
11: bare_func single_file

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