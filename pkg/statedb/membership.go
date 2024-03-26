package statedb

type MEMBERSHIP_TYPE string

const (
	INCLUSION         = MEMBERSHIP_TYPE("inclusion")
	EXCLUSION         = MEMBERSHIP_TYPE("exclusion")
	PARTIAL_INCLUSION = MEMBERSHIP_TYPE("partial_inclusion")
)
