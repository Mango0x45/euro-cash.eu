if exists('b:current_syntax')
	finish
endif
let b:current_syntax = 'mintage'

syntax match Comment    /^\s*#.*/
syntax match Number     /[0-9.]\+/
syntax match Identifier /\v\d{4}(-\S+)?/
syntax match String     /"[^"]\{-}"/

" ‘Todo’ is semantically a better syntax group, but it looks bad
syntax match Error      /?/
