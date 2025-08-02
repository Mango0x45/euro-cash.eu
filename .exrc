function s:SaveExcursion(cmd)
	let l:win = winsaveview()
	silent execute 'keepjumps %!' .. a:cmd
	call winrestview(l:win)
endfunction

autocmd BufNewFile,BufRead */cmd/extpo/*
	\ setlocal makeprg=go\ build\ ./cmd/extpo
autocmd BufNewFile,BufRead */cmd/extwiki/*
	\ setlocal makeprg=go\ build\ ./cmd/extwiki

autocmd FileType go autocmd BufWritePre <buffer>
	\ call s:SaveExcursion('gofmt -s')

nnoremap <silent> gM :wall \| make all-i18n<CR>
nnoremap <silent> <LocalLeader>t :vimgrep /\CTODO/ **/*<CR>

runtime plugin/netrwPlugin.vim
let &wildignore = netrw_gitignore#Hide() . ',.git/*,vendor/*'
let g:netrw_list_hide .= ',.*\\.gen\\..*'