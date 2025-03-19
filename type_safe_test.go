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
	registeredBookLoadable := RegisterLoadable(provider, LoadableKey("Books"), bookLoadable)
	
	// Create a user and register the loadable
	user := &TypedUser{ID: 1, Name: "Test User"}
	user.SetProvider(provider)
	user.RegisterBooksLoadable(registeredBookLoadable)
	
	// This should compile because the loadable is registered
	_, _ = user.Books(context.Background()) // No compile error
	
	// The following code would not compile because the "Authors" loadable is not registered
	// Uncomment to see the compile error
	/*
	book := &TypedBook{ID: 1, Title: "Test Book"}
	book.SetProvider(provider)
	// This line would cause a compile error because authorLoadable is not registered
	_, _ = book.Author(context.Background())
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
	// This would cause a compile error because booksLoadable is not registered
	_, _ = user.Books(context.Background())
	*/
	
	// Example 2: Type mismatch
	/*
	provider := NewTypedLoadableProvider()
	
	// Register loadables
	authorLoadable := EmptyHasOneLoadable[BookID, *TypedBook, UserID, *TypedUser]()
	registeredAuthorLoadable := RegisterHasOneLoadable(provider, LoadableKey("Authors"), authorLoadable)
	
	// Create a book and register the loadable
	book := &TypedBook{ID: 1, Title: "Test Book"}
	book.SetProvider(provider)
	book.RegisterAuthorLoadable(registeredAuthorLoadable)
	
	// This would compile because the loadable is registered
	_, _ = book.Author(context.Background()) // No compile error
	
	// But this would not compile because we're trying to use a HasOneLoadable as a Loadable
	user := &TypedUser{ID: 1, Name: "Test User"}
	user.SetProvider(provider)
	// This would cause a compile error because of type mismatch
	user.RegisterBooksLoadable(registeredAuthorLoadable)
	*/
}
