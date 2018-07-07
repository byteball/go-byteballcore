package types

type(
	// [fyi] used by AddPaidWitnessEvents in db/methods

	PaidWitnessEventT struct{
		Unit		UnitT
		Address		AddressT
		Delay		MCIndexT
	}
	PaidWitnessEventsT  = []PaidWitnessEventT
)
