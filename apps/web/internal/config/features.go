package config

// Features are the environment-dependent toggles. Branch on these, never on Env.
type Features struct {
	// EmailVerification requires a verified email before sign-in.
	EmailVerification bool
	// AccountPages mounts the sign-in, sign-up, dashboard and account pages.
	// Off: the site is the public portfolio only, and those paths are 404s.
	// The JSON API and packages/auth stay mounted either way.
	AccountPages bool
	// Email sends auth email (verification, password reset) through Resend,
	// so production needs FROM_EMAIL and RESEND_API_KEY. It follows
	// AccountPages: without accounts nothing needs a mail provider, the
	// mailer drops every message (cmd/server), and with EmailVerification
	// on in production no sign-up through the JSON API can complete.
	Email bool
	// MobileDevOrigins trusts Expo Go / dev-client origins (exp://) for auth
	// requests. Never on in production.
	MobileDevOrigins bool
}

// accountPages turns the template's account UI and auth email back on for
// every environment.
const accountPages = false

// FeaturesFor returns the toggles for env.
func FeaturesFor(env Env) Features {
	features := Features{AccountPages: accountPages, Email: accountPages}
	if env == Production {
		features.EmailVerification = true
		return features
	}
	features.MobileDevOrigins = true
	return features
}
