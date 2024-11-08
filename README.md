# Lsp in go

## Overview

This is an implementation of [educationlsp](https://github.com/tjdevries/educationalsp)  
i added the init.lue which has the code to setup the lsp and hook up the on attach  
this lsp works on markdown files

## installation steps

just run  
``` bash
  go build -o NAME main.go
```
and make sure to add the path to the neovim config  
``` lua
local test_client = vim.lsp.start_client({
	name = "lsp_in_go",
	cmd = { "PATH HERE" },
	on_attach = on_attach,
})
```
