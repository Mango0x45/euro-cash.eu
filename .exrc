set runtimepath+=contrib/vim

function s:SaveExcursion(cmd)
	let l:win = winsaveview()
	silent execute 'keepjumps %!' .. a:cmd
	call winrestview(l:win)
endfunction

autocmd BufNewFile,BufRead */data/mintages/*
	\   setfiletype mintage
	\ | setlocal nowrap
autocmd BufNewFile,BufRead */cmd/exttmpl/*
	\ setlocal makeprg=go\ build\ ./cmd/exttmpl
autocmd BufNewFile,BufRead */cmd/mfmt/*
	\ setlocal makeprg=go\ build\ ./cmd/mfmt

autocmd FileType go autocmd BufWritePre <buffer>
	\ call s:SaveExcursion('gofmt -s')
autocmd FileType mintage autocmd BufWritePre <buffer>
	\ call s:SaveExcursion('./mfmt')

nnoremap <silent> gM :wall \| make all-i18n<CR>
nnoremap <silent> <LocalLeader>t :vimgrep /\CTODO/ **/*<CR>

let &wildignore = netrw_gitignore#Hide() . ',.git/*,vendor/*'
let g:netrw_list_hide .= ',.*\\.gen\\..*'