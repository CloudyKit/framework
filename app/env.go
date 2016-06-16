package app

import "os"

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

func Getenv(key string, fallbacks ...string) (value string) {
	value, _ = LookupEnv(key, fallbacks...)
	return
}

func Env(def, key string, fallbacks ...string) string {
	if val, found := LookupEnv(key, fallbacks...); found {
		return val
	}
	return def
}
