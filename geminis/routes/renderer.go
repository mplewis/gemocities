package routes

import "git.sr.ht/~adnano/go-gemini"

type Renderer func(w gemini.ResponseWriter, tplName string, data any)
