(ns yuhta.check-gmail
  (:require [clojure.string :refer [join]]
            [yuhta.googleauth :refer [*scopes* *http-transport* authorize]])
  (:import com.google.api.client.googleapis.util.Utils
           (com.google.api.services.gmail GmailScopes Gmail$Builder))
  (:gen-class))

(alter-var-root #'*scopes* conj GmailScopes/GMAIL_READONLY)

(def ^:private headers '("From" "Subject"))
(def ^:private line-sep (System/getProperty "line.separator"))

(defn -main [& args]
  (let [messages (-> (Gmail$Builder. *http-transport*
                                     (Utils/getDefaultJsonFactory)
                                     (authorize))
                     (.setApplicationName (str (:ns (meta #'-main))))
                     .build
                     .users
                     .messages)
        hdrs     (for [msg (-> messages
                               (.list "me")
                               (.setQ "is:unread")
                               (.setIncludeSpamTrash true)
                               .execute
                               (get "messages"))]
                   (into {} (map (juxt #(get % "name") #(get % "value"))
                                 (-> (.get messages "me" (get msg "id"))
                                     (.setFormat "metadata")
                                     (.setMetadataHeaders headers)
                                     .execute
                                     (get "payload")
                                     (get "headers")))))]
    (print (join line-sep
                 (for [m hdrs]
                   (join (for [k headers]
                           (str k ": " (get m k) line-sep))))))
    (flush)))
