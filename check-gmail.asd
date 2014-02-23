(asdf:defsystem #:check-gmail
  :serial t
  :author "Jimmy Lu <gongchuo.lu@gmail.com>"
  :depends-on (#:alexandria
               #:cl-netrc
               #:cl-ppcre
               #:clonsigna
               #:cl-rfc2047)
  :components ((:file "package")
               (:file "check-gmail")))
