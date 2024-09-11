;;; mintage-mode.el --- Major mode for mintage data files -*- lexical-binding: t; -*-

(defconst mintage--number  (rx (or digit ?.)))
(defconst mintage--string  (rx "\"" (*? (not "\"")) "\""))
(defconst mintage--unknown (rx ??))
(defconst mintage--comment (rx bol (* space) ?# (* any)))
(defconst mintage--year    (rx bol (= 4 digit) (? ?- (+ (not space)))))

(defun mintage--font-lock-keywords ()
  `((,mintage--year    . font-lock-constant-face)
    (,mintage--number  . 'font-lock-number-face)
    (,mintage--string  . font-lock-string-face)
    (,mintage--comment . font-lock-comment-face)
    (,mintage--unknown . font-lock-warning-face)))

;;;###autoload
(add-to-list 'auto-mode-alist (cons "mintages/[a-z]\\{2\\}\\'" 'mintage-mode))

;;;###autoload
(define-derived-mode mintage-mode conf-mode "Mintage"
  "Major mode for editing euro-cash.eu mintage files"
  (setq comment-start "#"
        comment-end "")
  (setq font-lock-defaults '(mintage--font-lock-keywords)))

(provide 'mintage-mode)
