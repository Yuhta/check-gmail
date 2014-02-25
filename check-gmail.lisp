(in-package #:check-gmail)

(define-constant +host+ "imap.gmail.com" :test #'string=)

(defmacro with-imap-connection ((socket) &body body)
  (with-gensyms (netrc)
    `(let* ((,netrc (netrc:lookup (netrc:make-netrc) +host+))
            (,socket (make-imap :host +host+ :port (getf ,netrc :port)
                                :ssl-p t)))
       (unwind-protect (progn
                         (cmd-connect ,socket)
                         (cmd-login ,socket
                                    (getf ,netrc :login)
                                    (getf ,netrc :password))
                         ,@body)
         (cmd-logout ,socket)))))

(defun decode (string)
  (cl-rfc2047:decode* string :errorp nil))

(defun parse-envelope (reply)
  (ppcre:register-groups-bind (envelope-str)
      ("^\\* \\d+ FETCH \\(ENVELOPE (\\(.*\\))\\)$" reply :sharedp t)
    (when envelope-str
      (let ((envelope (read-from-string envelope-str)))
        (list :from (decode (format nil "~{~:[~;~:*~a ~]~*<~a@~a>~}"
                                    (first (third envelope))))
              :subject (decode (second envelope)))))))

(defun new-mail-summaries (imap-socket)
  (cmd-examine imap-socket "INBOX")
  (when-let ((ids (parse-search (cmd-search imap-socket
                                            :criteria "UNSEEN"))))
    (mapcar #'parse-envelope
            (cmd-fetch imap-socket (format nil "~{~a~^,~}" ids)
                       :criteria "ENVELOPE"))))

(defun main (&optional args)
  (declare (ignore args))
  (with-imap-connection (socket)
    (when-let ((summaries (new-mail-summaries socket)))
      (format t "~a new message~:*~[~;~:;s~] in your Gmail Inbox~%"
              (length summaries))
      (format t "~{~%From: ~a~%Subject: ~a~%~}"
              (mapcan (lambda (summary)
                        (list (getf summary :from)
                              (getf summary :subject)))
                      summaries)))))

;;; TODO: http://tools.ietf.org/html/rfc2047
