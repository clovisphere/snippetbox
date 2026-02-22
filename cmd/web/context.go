package main

// contextKey is a custom type used for keys in the request context.
// Using a private, custom type prevents "collisions" between our
// middleware and other external packages.
type contextKey string

// isAuthenticatedContextKey is the unique key used to store and
// retrieve the authentication status from a request's context.
const isAuthenticatedContextKey = contextKey("isAuthenticated")
