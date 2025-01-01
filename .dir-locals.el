((nil . ((eval . (progn
                   (add-to-list
                    'load-path
                    (expand-file-name
                     "contrib/emacs"
                     (locate-dominating-file default-directory "contrib/emacs")))
                   (require 'mintage-mode)))))
 (go-ts-mode . ((comment-start . "/* ")
                (comment-end . " */")
                (require-final-newline . t)))
 (mhtml-mode . ((fill-column . 79))))