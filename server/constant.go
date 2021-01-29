package server

var allowHeaders = []string{"X-Requested-With", "Content-Type", "Authorization", "Username"}
var allowMethods = []string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}

const authorizationHeader = "Authorization"
const usernameHeader = "Username"

const bearerTokenPattern = "^Bearer\\s(.*)"
