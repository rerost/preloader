package preloader

import (
	"testing"
)

// TestTypedLoadableProviderRegistration tests that loadables can be registered and retrieved correctly
func TestTypedLoadableProviderRegistration(t *testing.T) {
	provider := NewTypedLoadableProvider()
	
	// Register a loadable
	bookLoadable := EmptyLoadable[UserID, *TypedUser, BookID, *TypedBook]()
	registeredBookLoadable := RegisterLoadable(provider, LoadableKey("Books"), bookLoadable)
	
	// Retrieve the loadable
	retrievedLoadable := GetRegisteredLoadable(provider, registeredBookLoadable)
	
	// Verify that it's the same loadable
	if retrievedLoadable != bookLoadable {
		t.Errorf("Retrieved loadable is not the same as the registered loadable")
	}
	
	// Register a HasOneLoadable
	authorLoadable := EmptyHasOneLoadable[BookID, *TypedBook, UserID, *TypedUser]()
	registeredAuthorLoadable := RegisterHasOneLoadable(provider, LoadableKey("Authors"), authorLoadable)
	
	// Retrieve the HasOneLoadable
	retrievedHasOneLoadable := GetRegisteredHasOneLoadable(provider, registeredAuthorLoadable)
	
	// Verify that it's the same HasOneLoadable
	if retrievedHasOneLoadable != authorLoadable {
		t.Errorf("Retrieved HasOneLoadable is not the same as the registered HasOneLoadable")
	}
}

// TestTypeSafetyAtCompileTime demonstrates that type safety is enforced at compile time
func TestTypeSafetyAtCompileTime(t *testing.T) {
	// This test doesn't actually run any code
	// It's just to demonstrate compile-time type safety
	
	// Example: Type mismatch between RegisteredLoadable and RegisteredHasOneLoadable
	/*
	provider := NewTypedLoadableProvider()
	
	// Register a loadable
	bookLoadable := EmptyLoadable[UserID, *TypedUser, BookID, *TypedBook]()
	registeredBookLoadable := RegisterLoadable(provider, LoadableKey("Books"), bookLoadable)
	
	// This would cause a compile error because we're trying to use GetRegisteredHasOneLoadable with a RegisteredLoadable
	_ = GetRegisteredHasOneLoadable(provider, registeredBookLoadable)
	*/
	
	// Example: Using an unregistered loadable
	/*
	provider := NewTypedLoadableProvider()
	
	// Create an unregistered loadable (using NotRegistered phantom type)
	unregisteredLoadable := RegisteredLoadable[NotRegistered, UserID, *TypedUser, BookID, *TypedBook]{
		Key: LoadableKey("Books"),
		Loadable: EmptyLoadable[UserID, *TypedUser, BookID, *TypedBook](),
	}
	
	// This would cause a compile error because GetRegisteredLoadable requires a RegisteredLoadable with Registered type
	_ = GetRegisteredLoadable(provider, unregisteredLoadable)
	*/
}
