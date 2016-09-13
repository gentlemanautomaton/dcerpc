package presentationcontext

// Result represents the result of a presentation context negotiation.
type Result uint16

const (
	// Acceptance indicates that the proposed presentation context was accepted.
	Acceptance Result = iota

	// UserRejection indicates that the proposed presentation context was rejected
	// by the user.
	UserRejection

	// ProviderRejection indicates that the proposed presentation context was rejected
	// by the provider.
	ProviderRejection
)
