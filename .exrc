" We make use of this feature, so set this in the environment so that all
" Ex-calls to templ are aware of this
call setenv("TEMPL_EXPERIMENT", "rawgo")

function s:SaveExcursion(cmd)
	let l:win = winsaveview()
	execute "%!" .. a:cmd
	call winrestview(l:win)
endfunction

autocmd BufWritePre *.go    call s:SaveExcursion("gofmt -s")
autocmd BufWritePre *.templ call s:SaveExcursion("templ fmt")

nnoremap <silent> gM :wall \| make all-i18n<CR>

let g:netrw_list_hide .= ",.*_templ\\.go$"
