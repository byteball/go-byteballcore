package types

type(
	AsyncFunctorT   = func () ErrorT
	AsyncFunctorsT  = []AsyncFunctorT
)
