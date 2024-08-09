set runtimepath+=contrib

" We make use of this feature, so set this in the environment so that all
" Ex-calls to templ are aware of this
call setenv("TEMPL_EXPERIMENT", "rawgo")

function s:SaveExcursion(cmd)
	let l:win = winsaveview()
	silent execute "keepjumps %!" .. a:cmd
	call winrestview(l:win)
endfunction

autocmd BufNewFile,BufRead */data/mintages/* setfiletype mintage
autocmd BufNewFile,BufRead */cmd/mfmt/*      setlocal makeprg=go\ build\ ./cmd/mfmt

autocmd FileType go autocmd BufWritePre <buffer>
	\ call s:SaveExcursion("gofmt -s")
autocmd FileType mintage autocmd BufWritePre <buffer>
	\ call s:SaveExcursion("./mfmt")
autocmd FileType templ autocmd BufWritePre <buffer>
	\ call s:SaveExcursion("templ fmt")

nnoremap <silent> gM :wall \| make all-i18n<CR>
nnoremap <silent> <LocalLeader>t :vimgrep /TODO/ **/*<CR>

let &wildignore = netrw_gitignore#Hide() . ',*_templ.go,.git/*,vendor/*'
let g:netrw_list_hide .= ",.*_templ\\.go$"
