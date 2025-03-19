package preloader

import (
	"context"
	"testing"
)

// TestTypedLoadableProvider tests the TypedLoadableProvider
func TestTypedLoadableProvider(t *testing.T) {
	provider := NewTypedLoadableProvider()
	
	// Register a loadable
	bookLoadable := EmptyLoadable[UserID, *TypedUser, BookID, *TypedBook]()
	RegisterTypedLoadable(provider, LoadableKey("Books"), bookLoadable)
	
	// This should compile because the loadable is registered
	user := &TypedUser{ID: 1, Name: "Test User"}
	user.SetProvider(provider)
	_, _ = user.Books(context.Background()) // No compile error
	
	// The following code would not compile because the "Authors" loadable is not registered
	// Uncomment to see the compile error
	/*
	book := &TypedBook{ID: 1, Title: "Test Book"}
	book.SetProvider(provider)
	_, _ = book.Author(context.Background()) // Compile error: "Authors" loadable not registered
	*/
}

// TestCompileTimeErrors demonstrates compile-time errors for unregistered loadables
func TestCompileTimeErrors(t *testing.T) {
	// This test doesn't actually run any code
	// It's just to demonstrate compile-time errors
	
	// Example 1: Missing loadable registration
	/*
	provider := NewTypedLoadableProvider()
	user := &TypedUser{ID: 1, Name: "Test User"}
	user.SetProvider(provider)
	_, _ = user.Books(context.Background()) // Compile error: "Books" loadable not registered
	*/
	
	// Example 2: Type mismatch
	/*
	provider := NewTypedLoadableProvider()
	authorLoadable := EmptyHasOneLoadable[BookID, *TypedBook, UserID, *TypedUser]()
	provider.RegisterTypedHasOneLoadable(LoadableKey("Authors"), authorLoadable)
	
	// This would compile because the loadable is registered
	book := &TypedBook{ID: 1, Title: "Test Book"}
	book.SetProvider(provider)
	_, _ = book.Author(context.Background()) // No compile error
	
	// But this would not compile because we're trying to use a HasOneLoadable as a Loadable
	provider.RegisterTypedHasOneLoadable(LoadableKey("Books"), authorLoadable)
	user := &TypedUser{ID: 1, Name: "Test User"}
	user.SetProvider(provider)
	_, _ = user.Books(context.Background()) // Compile error: type mismatch
	*/
}
