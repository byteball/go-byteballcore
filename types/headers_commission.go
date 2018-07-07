package types

type(
	// [fyi] used by AddContributions in db/methods

	ContributionT struct{
		Payer_unit	UnitT
		Address		AddressT
		Amount		AmountT
	}
	ContributionsT	= []ContributionT
)
