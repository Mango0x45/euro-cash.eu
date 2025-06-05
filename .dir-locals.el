((csv-mode . ((number-format-predicate
               . (lambda (beg end)
                   (save-excursion
                     (goto-char beg)
                     (not (bolp)))))))
 (go-ts-mode . ((comment-start . "/* ")
                (comment-end . " */")
                (comment-continue . "   ")
                (require-final-newline . t)))
 (mhtml-mode . ((fill-column . 79))))