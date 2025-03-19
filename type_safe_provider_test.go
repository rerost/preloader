package preloader

import (
	"testing"
)

// TestTypedLoadableProviderRegistration tests that loadables can be registered and retrieved correctly
func TestTypedLoadableProviderRegistration(t *testing.T) {
	provider := NewTypedLoadableProvider()
	
	// Register a loadable
	bookLoadable := EmptyLoadable[UserID, *TypedUser, BookID, *TypedBook]()
	RegisterTypedLoadable(provider, LoadableKey("Books"), bookLoadable)
	
	// Retrieve the loadable
	retrievedLoadable := MustGetLoadable[UserID, *TypedUser, BookID, *TypedBook](provider, LoadableKey("Books"))
	
	// Verify that it's the same loadable
	if retrievedLoadable != bookLoadable {
		t.Errorf("Retrieved loadable is not the same as the registered loadable")
	}
	
	// Register a HasOneLoadable
	authorLoadable := EmptyHasOneLoadable[BookID, *TypedBook, UserID, *TypedUser]()
	RegisterTypedHasOneLoadable(provider, LoadableKey("Authors"), authorLoadable)
	
	// Retrieve the HasOneLoadable
	retrievedHasOneLoadable := MustGetHasOneLoadable[BookID, *TypedBook, UserID, *TypedUser](provider, LoadableKey("Authors"))
	
	// Verify that it's the same HasOneLoadable
	if retrievedHasOneLoadable != authorLoadable {
		t.Errorf("Retrieved HasOneLoadable is not the same as the registered HasOneLoadable")
	}
}

// TestTypedLoadableProviderPanic tests that the provider panics when there's a type mismatch
func TestTypedLoadableProviderPanic(t *testing.T) {
	provider := NewTypedLoadableProvider()
	
	// Register a loadable with the wrong type
	bookLoadable := EmptyLoadable[UserID, *TypedUser, BookID, *TypedBook]()
	provider.RegisterLoadable("Authors", bookLoadable) // Using the untyped RegisterLoadable
	
	// This should panic because we're trying to retrieve a HasOneLoadable but registered a Loadable
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	
	_ = MustGetHasOneLoadable[BookID, *TypedBook, UserID, *TypedUser](provider, LoadableKey("Authors"))
}
