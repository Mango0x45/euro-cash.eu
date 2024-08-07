function s:SaveExcursion(cmd)
	let l:win = winsaveview()
	execute "%!" .. a:cmd
	call winrestview(l:win)
endfunction

autocmd BufWritePre *.go
	\ call s:SaveExcursion("gofmt -s")
autocmd BufWritePre *.templ
	\ call s:SaveExcursion("templ fmt | sed 's/{ {/{{/; s/} }/}}/'")

nnoremap gM :make all-i18n<CR>

let g:netrw_list_hide .= ",.*_templ.go$"
