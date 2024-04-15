package reqedit

import "net/http"

type EditRequestFunc func(*http.Request)
