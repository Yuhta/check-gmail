(defpackage #:check-gmail
  (:use #:cl
        #:clonsigna)
  (:import-from #:alexandria
                #:with-gensyms
                #:define-constant
                #:when-let))
