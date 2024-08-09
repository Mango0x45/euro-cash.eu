if exists('b:current_syntax')
	finish
endif
let b:current_syntax = 'mintage'

syntax keyword Keyword BEGIN CIRC BU PROOF
syntax match   Label   /[^\s]\+\*\?:/
syntax match   Number  /[0-9\.]\+/
syntax match   Todo    /?/
