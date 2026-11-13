package config

// Features are the environment-dependent toggles. Branch on these, never on Env.
type Features struct {
	// EmailVerification requires a verified email before sign-in.
	EmailVerification bool
	// AccountPages mounts the sign-in, sign-up, dashboard and account pages.
	// Off: the site is the public portfolio only, and those paths are 404s.
	// The JSON API and packages/auth stay mounted either way.
	AccountPages bool
	// MobileDevOrigins trusts Expo Go / dev-client origins (exp://) for auth
	// requests. Never on in production.
	MobileDevOrigins bool
}

// FeaturesFor returns the toggles for env.
func FeaturesFor(env Env) Features {
	if env == Production {
		return Features{EmailVerification: true}
	}
	return Features{EmailVerification: false, MobileDevOrigins: true}
}
