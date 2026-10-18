package tests

import "os"

var Port = 7540
var DBFile = "../scheduler.db"
var FullNextDate = false
var Search = false

// Token заполняется через TODO_TOKEN или автоматически через /api/signin
// в TestMain. Не публикуйте действующие токены в репозитории.
var Token = os.Getenv("TODO_TOKEN")
