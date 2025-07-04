package sqx

import (
	"github.com/blink-io/sq"
)

type Constructor func(sq.DB) sq.DB

// Chain acts as a list of sq.DB constructors.
// Chain is effectively immutable:
// once created, it will always hold
// the same set of constructors in the same order.
type Chain struct {
	constructors []Constructor
}

// NewChain creates a new chain,
// memorizing the given list of middleware constructors.
// New serves no other function,
// constructors are only called upon a call to Then().
func NewChain(constructors ...Constructor) Chain {
	return Chain{append(([]Constructor)(nil), constructors...)}
}

// Then chains the middleware and returns the final sq.DB.
//
//	New(m1, m2, m3).Then(db)
//
// is equivalent to:
//
//	m1(m2(m3(db)))
//
// When the request comes in, it will be passed to m1, then m2, then m3
// and finally, the given handler
// (assuming every middleware calls the following one).
//
// A chain can be safely reused by calling Then() several times.
//
//	stdStack := alice.New(ratelimitHandler, csrfHandler)
//	indexPipe = stdStack.Then(indexHandler)
//	authPipe = stdStack.Then(authHandler)
//
// Note that constructors are called on every call to Then()
// and thus several instances of the same middleware will be created
// when a chain is reused in this way.
// For proper middleware, this should cause no problems.
func (c Chain) Then(db sq.DB) sq.DB {
	if db != nil {
		for i := range c.constructors {
			db = c.constructors[len(c.constructors)-1-i](db)
		}
		return db
	}
	return nil
}

// Append extends a chain, adding the specified constructors
// as the last ones in the request flow.
//
// Append returns a new chain, leaving the original one untouched.
//
//	stdChain := alice.New(m1, m2)
//	extChain := stdChain.Append(m3, m4)
//	// requests in stdChain go m1 -> m2
//	// requests in extChain go m1 -> m2 -> m3 -> m4
func (c Chain) Append(constructors ...Constructor) Chain {
	newCons := make([]Constructor, 0, len(c.constructors)+len(constructors))
	newCons = append(newCons, c.constructors...)
	newCons = append(newCons, constructors...)

	return Chain{newCons}
}

// Extend extends a chain by adding the specified chain
// as the last one in the request flow.
//
// Extend returns a new chain, leaving the original one untouched.
//
//	stdChain := alice.New(m1, m2)
//	ext1Chain := alice.New(m3, m4)
//	ext2Chain := stdChain.Extend(ext1Chain)
//	// requests in stdChain go  m1 -> m2
//	// requests in ext1Chain go m3 -> m4
//	// requests in ext2Chain go m1 -> m2 -> m3 -> m4
//
// Another example:
//
//	aHtmlAfterNosurf := alice.New(m2)
//	aHtml := alice.New(m1, func(db sq.DB) sq.DB {
//		csrf := nosurf.New(h)
//		csrf.SetFailureHandler(aHtmlAfterNosurf.ThenFunc(csrfFail))
//		return csrf
//	}).Extend(aHtmlAfterNosurf)
//	// requests to aHtml hitting nosurfs success handler go m1 -> nosurf -> m2 -> target-handler
//	// requests to aHtml hitting nosurfs failure handler go m1 -> nosurf -> m2 -> csrfFail
func (c Chain) Extend(chain Chain) Chain {
	return c.Append(chain.constructors...)
}

func ChainFunc(db sq.DB, constructors ...Constructor) sq.DB {
	if len(constructors) == 0 || db == nil {
		return db
	}
	chain := NewChain(constructors...)
	return chain.Then(db)
}
