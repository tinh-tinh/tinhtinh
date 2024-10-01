package helmet

import (
	"fmt"
	"net/http"
	"strings"
)

// CrossOriginEmbedderPolicy defines the policy for cross-origin resource isolation.
// Policy can be "require-corp", "credentialless", or other valid values.
type CrossOriginEmbedderPolicy struct {
	Policy  string // The policy string value, e.g., "require-corp".
	Enabled bool   // Whether the policy is enabled or not.
}

// CrossOriginOpenerPolicy defines the policy for controlling the browsing context group.
// Policy can be "same-origin", "same-origin-allow-popups", etc.
type CrossOriginOpenerPolicy struct {
	Policy  string // The policy string value.
	Enabled bool   // Whether the policy is enabled or not.
}

// CrossOriginResourcePolicy specifies how resources are shared across origins.
// Policy can be "same-origin", "same-site", or "cross-origin".
type CrossOriginResourcePolicy struct {
	Policy  string // The policy string value.
	Enabled bool   // Whether the policy is enabled or not.
}

// ReferrerPolicy specifies how much referrer information should be included with requests.
type ReferrerPolicy struct {
	Policy  interface{} // Referrer policy, e.g., "no-referrer", "origin-when-cross-origin".
	Enabled bool        // Whether the policy is enabled or not.
}

// StrictTransportSecurity sets the HSTS policy for the application.
type StrictTransportSecurity struct {
	MaxAge            int  // The time, in seconds, for how long the browser should remember that this site is only accessible using HTTPS.
	IncludeSubDomains bool // Whether to apply the policy to subdomains as well.
	Preload           bool // Whether the site is included in browsers' HSTS preload lists.
	Enabled           bool // Whether the policy is enabled or not.
}

// OptionSecurityPolicy defines the Content Security Policy (CSP) options.
type OptionSecurityPolicy struct {
	UseDefaults             bool     // Whether to use the default CSP settings.
	DefaultSrc              []string // The default source(s) from which resources can be loaded.
	StyleSrc                []string // The allowed source(s) for styles.
	ImgSrc                  []string // The allowed source(s) for images.
	FontSrc                 []string // The allowed source(s) for fonts.
	ScriptSrc               []string // The allowed source(s) for scripts.
	ManifestSrc             []string // The allowed source(s) for web manifests.
	FrameSrc                []string // The allowed source(s) for frames.
	UpgradeInsecureRequests bool     // Whether to automatically upgrade insecure requests (http to https).
	ObjectSrc               []string // The allowed source(s) for object tags.
	ReportOnly              bool     // If true, only report violations but do not block resources.
}

// ContentSecurityPolicy manages the security policy for content in a page.
type ContentSecurityPolicy struct {
	OptionSecurityPolicy *OptionSecurityPolicy // The detailed policy options.
}

// XContentTypeOptions controls MIME type sniffing protection.
type XContentTypeOptions struct {
	Enabled bool   // Whether the policy is enabled or not.
	Value   string // The value for X-Content-Type-Options, typically "nosniff".
}

// XDnsPrefetchControl controls browser DNS prefetching.
type XDnsPrefetchControl struct {
	Enabled bool   // Whether the policy is enabled or not.
	Value   string // The value for X-DNS-Prefetch-Control, typically "on" or "off".
}

// XDownloadOptions controls file download security.
type XDownloadOptions struct {
	Enabled bool   // Whether the policy is enabled or not.
	Value   string // The value for X-Download-Options, typically "noopen".
}

// XFrameOptions controls the ability to embed the page in a frame.
type XFrameOptions struct {
	Enabled bool   // Whether the policy is enabled or not.
	Action  string // The value for X-Frame-Options, typically "DENY" or "SAMEORIGIN".
}

// XPermittedCrossDomainPolicies specifies whether Flash and PDF files can be used in cross-domain requests.
type XPermittedCrossDomainPolicies struct {
	Enabled           bool   // Whether the policy is enabled or not.
	PermittedPolicies string // Specifies the allowed cross-domain policies, e.g., "none", "master-only".
}

// XPoweredBy controls the X-Powered-By header to hide or reveal server information.
type XPoweredBy struct {
	Enabled bool   // Whether the policy is enabled or not.
	Value   string // The value for X-Powered-By, typically the application or framework name.
}

// XXSSProtection controls the X-XSS-Protection header to mitigate cross-site scripting (XSS) attacks.
type XXSSProtection struct {
	Enabled bool   // Whether the policy is enabled or not.
	Value   string // The value for X-XSS-Protection, e.g., "1; mode=block" or "0" to disable.
}

// Helmet is the structure that holds all security headers.
type Options struct {
	CrossOriginEmbedderPolicy     CrossOriginEmbedderPolicy     // COEP header configuration.
	CrossOriginOpenerPolicy       CrossOriginOpenerPolicy       // COOP header configuration.
	CrossOriginResourcePolicy     CrossOriginResourcePolicy     // CORP header configuration.
	ReferrerPolicy                ReferrerPolicy                // Referrer policy configuration.
	ContentSecurityPolicy         ContentSecurityPolicy         // CSP configuration.
	StrictTransportSecurity       StrictTransportSecurity       // HSTS configuration.
	XContentTypeOptions           XContentTypeOptions           // MIME type sniffing protection.
	XDnsPrefetchControl           XDnsPrefetchControl           // DNS prefetching control.
	XDownloadOptions              XDownloadOptions              // File download security.
	XFrameOptions                 XFrameOptions                 // Frame embedding control.
	XPermittedCrossDomainPolicies XPermittedCrossDomainPolicies // Cross-domain policy configuration.
	XPoweredBy                    XPoweredBy                    // X-Powered-By header configuration.
	XXSSProtection                XXSSProtection                // X-XSS-Protection header configuration.
}

func Handler(opt Options) func(h http.Handler) http.Handler {
	helmet := newDefault(opt)
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Cross-Origin-Embedder-Policy
			if helmet.CrossOriginEmbedderPolicy.Enabled {
				policy := "require-corp"
				if helmet.CrossOriginEmbedderPolicy.Policy == "credentialless" {
					policy = "credentialless"
				}
				w.Header().Set("Cross-Origin-Embedder-Policy", policy)
			}

			// Cross-Origin-Opener-Policy
			if helmet.CrossOriginOpenerPolicy.Enabled {
				openerPolicy := "same-origin"
				if helmet.CrossOriginOpenerPolicy.Policy == "same-origin-allow-popups" {
					openerPolicy = "same-origin-allow-popups"
				}
				w.Header().Set("Cross-Origin-Opener-Policy", openerPolicy)
			}

			// Cross-Origin-Resource-Policy
			if helmet.CrossOriginResourcePolicy.Enabled {
				resourcePolicy := "same-origin"
				if helmet.CrossOriginResourcePolicy.Policy == "same-site" {
					resourcePolicy = "same-site"
				}
				w.Header().Set("Cross-Origin-Resource-Policy", resourcePolicy)
			}

			// Strict-Transport-Security
			if helmet.StrictTransportSecurity.Enabled {
				if helmet.StrictTransportSecurity.MaxAge > 0 {
					hsts := fmt.Sprintf("max-age=%d", helmet.StrictTransportSecurity.MaxAge)
					if helmet.StrictTransportSecurity.IncludeSubDomains {
						hsts += "; includeSubDomains"
					}
					if helmet.StrictTransportSecurity.Preload {
						hsts += "; preload"
					}
					w.Header().Set("Strict-Transport-Security", hsts)
				}
			}

			// Referrer-Policy
			if helmet.ReferrerPolicy.Enabled {
				switch policy := helmet.ReferrerPolicy.Policy.(type) {
				case string:
					w.Header().Set("Referrer-Policy", policy)
				case []string:
					w.Header().Set("Referrer-Policy", strings.Join(policy, ","))
				}
			}

			// X-Content-Type-Options
			if helmet.XContentTypeOptions.Enabled {
				w.Header().Set("X-Content-Type-Options", helmet.XContentTypeOptions.Value)
			}

			// X-DNS-Prefetch-Control
			if helmet.XDnsPrefetchControl.Enabled {
				w.Header().Set("X-DNS-Prefetch-Control", helmet.XDnsPrefetchControl.Value)
			} else {
				w.Header().Set("X-DNS-Prefetch-Control", "off")
			}

			// X-Download-Options
			if helmet.XDownloadOptions.Enabled {
				w.Header().Set("X-Download-Options", helmet.XDownloadOptions.Value)
			}

			// X-Frame-Options
			if helmet.XFrameOptions.Enabled {
				w.Header().Set("X-Frame-Options", helmet.XFrameOptions.Action)
			}

			// X-Permitted-Cross-Domain-Policies
			if helmet.XPermittedCrossDomainPolicies.Enabled {
				w.Header().Set("X-Permitted-Cross-Domain-Policies", helmet.XPermittedCrossDomainPolicies.PermittedPolicies)
			}

			// X-Powered-By
			if helmet.XPoweredBy.Enabled {
				w.Header().Set("X-Powered-By", helmet.XPoweredBy.Value)
			} else {
				w.Header().Del("X-Powered-By") // Xóa header nếu không được phép
			}

			// X-XSS-Protection
			if helmet.XXSSProtection.Enabled {
				w.Header().Set("X-XSS-Protection", helmet.XXSSProtection.Value)
			}

			// Content-Security-Policy or Content-Security-Policy-Report-Only
			if helmet.ContentSecurityPolicy.OptionSecurityPolicy != nil {
				cspHeader := handleOptions(*helmet.ContentSecurityPolicy.OptionSecurityPolicy)
				if helmet.ContentSecurityPolicy.OptionSecurityPolicy.ReportOnly {
					w.Header().Set("Content-Security-Policy-Report-Only", cspHeader)
				} else {
					w.Header().Set("Content-Security-Policy", cspHeader)
				}
			}

			h.ServeHTTP(w, r)
		})
	}
}

// Hàm khởi tạo Helmet với tùy chọn cấu hình
func newDefault(options Options) *Options {
	h := &Options{
		StrictTransportSecurity:       options.StrictTransportSecurity,
		ContentSecurityPolicy:         options.ContentSecurityPolicy,
		XContentTypeOptions:           options.XContentTypeOptions,
		XDnsPrefetchControl:           options.XDnsPrefetchControl,
		XDownloadOptions:              options.XDownloadOptions,
		XFrameOptions:                 options.XFrameOptions,
		CrossOriginEmbedderPolicy:     options.CrossOriginEmbedderPolicy,
		CrossOriginOpenerPolicy:       options.CrossOriginOpenerPolicy,
		CrossOriginResourcePolicy:     options.CrossOriginResourcePolicy,
		XPoweredBy:                    options.XPoweredBy,
		XXSSProtection:                options.XXSSProtection,
		XPermittedCrossDomainPolicies: options.XPermittedCrossDomainPolicies,
		ReferrerPolicy:                options.ReferrerPolicy,
	}

	if options.ContentSecurityPolicy.OptionSecurityPolicy != nil {
		h.ContentSecurityPolicy = ContentSecurityPolicy{
			OptionSecurityPolicy: options.ContentSecurityPolicy.OptionSecurityPolicy,
		}
	} else {
		h.ContentSecurityPolicy = ContentSecurityPolicy{
			OptionSecurityPolicy: &OptionSecurityPolicy{
				DefaultSrc: []string{"'self'"},
			},
		}
	}

	if !h.XContentTypeOptions.Enabled {
		h.XContentTypeOptions = XContentTypeOptions{
			Enabled: true,
			Value:   "nosniff",
		}
	}

	if !h.XDnsPrefetchControl.Enabled {
		h.XDnsPrefetchControl = XDnsPrefetchControl{
			Enabled: true,
			Value:   "off",
		}
	}

	if !h.XFrameOptions.Enabled {
		h.XFrameOptions = XFrameOptions{
			Enabled: true,
			Action:  "SAMEORIGIN",
		}
	}

	if !h.XDownloadOptions.Enabled {
		h.XDownloadOptions = XDownloadOptions{
			Enabled: true,
			Value:   "noopen",
		}
	}

	if !h.XXSSProtection.Enabled {
		h.XXSSProtection = XXSSProtection{
			Enabled: true,
			Value:   "0",
		}
	}

	if !h.StrictTransportSecurity.Enabled {
		h.StrictTransportSecurity = StrictTransportSecurity{
			Enabled: true,
			MaxAge:  31536000,
			Preload: true,
		}
	}

	return h
}
func handleOptions(cspOptions OptionSecurityPolicy) string {
	policies := []string{}

	if !cspOptions.UseDefaults {
		if len(cspOptions.DefaultSrc) > 0 {
			policies = append(policies, fmt.Sprintf("default-src %s", strings.Join(cspOptions.DefaultSrc, " ")))
		}
		if len(cspOptions.ScriptSrc) > 0 {
			policies = append(policies, fmt.Sprintf("script-src %s", strings.Join(cspOptions.ScriptSrc, " ")))
		}
		if len(cspOptions.ImgSrc) > 0 {
			policies = append(policies, fmt.Sprintf("img-src %s", strings.Join(cspOptions.ImgSrc, " ")))
		}
		if len(cspOptions.StyleSrc) > 0 {
			policies = append(policies, fmt.Sprintf("style-src %s", strings.Join(cspOptions.StyleSrc, " ")))
		}
		if len(cspOptions.FontSrc) > 0 {
			policies = append(policies, fmt.Sprintf("font-src %s", strings.Join(cspOptions.FontSrc, " ")))
		}
		if len(cspOptions.ManifestSrc) > 0 {
			policies = append(policies, fmt.Sprintf("manifest-src %s", strings.Join(cspOptions.ManifestSrc, " ")))
		}
		if len(cspOptions.FrameSrc) > 0 {
			policies = append(policies, fmt.Sprintf("frame-src %s", strings.Join(cspOptions.FrameSrc, " ")))
		}
		if len(cspOptions.ObjectSrc) > 0 {
			policies = append(policies, fmt.Sprintf("object-src %s", strings.Join(cspOptions.ObjectSrc, " ")))
		}
	}

	if cspOptions.UpgradeInsecureRequests {
		policies = append(policies, "upgrade-insecure-requests")
	}
	return strings.Join(policies, "; ")
}
