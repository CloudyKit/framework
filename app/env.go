package app

import "os"

// LookupEnv same as os.LookupEnv, but with fallback support
func LookupEnv(Key string, fallbacks ...string) (value string, found bool) {
	value, found = os.LookupEnv(Key)
	if !found {
		for _, Key := range fallbacks {
			value, found = os.LookupEnv(Key)
			if found {
				return
			}
		}
	}
	return
}

// Getenv same as os.Getenv, but with fallback support
func Getenv(key string, fallbacks ...string) (value string) {
	value, _ = LookupEnv(key, fallbacks...)
	return
}

// Env works as Getenv but the first argument is a default value
func Env(def, key string, fallbacks ...string) string {
	if val, found := LookupEnv(key, fallbacks...); found {
		return val
	}
	return def
}
