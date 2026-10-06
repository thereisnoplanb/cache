package RemovalSource

// Specifies the source that caused a cache item to be removed.
type Enum int

const (
	// The item was removed explicitly by a caller.
	Manual = 1
	// The item was removed automatically after its expiration duration elapsed.
	Expiration = 2
)

const undefined string = "Unknown removal source"

var dictionary = map[Enum]string{
	Manual:     "Manual",
	Expiration: "Expiration",
}

// Returns the display name of the removal source.
//
// # Returns
//
//	name string
//
// The name associated with the value.
func (enum Enum) String() (name string) {
	if name, isDefined := dictionary[enum]; isDefined {
		return name
	}
	return undefined
}

// Determines whether the value is a defined removal source.
//
// # Returns
//
//	isDefined bool
//
// True if the value is defined; otherwise, false.
func (enum Enum) IsDefined() (isDefined bool) {
	_, isDefined = dictionary[enum]
	return isDefined
}

// Determines whether the specified value is a defined removal source.
//
// # Parameters
//
//	enum Enum
//
// The value to validate.
//
// # Returns
//
//	isDefined bool
//
// True if the value is defined; otherwise, false.
func IsDefined(enum Enum) (isDefined bool) {
	return enum.IsDefined()
}

// Returns all defined removal source values.
//
// # Returns
//
//	enums []Enum
//
// A slice containing all defined removal source values.
func GetValues() (enums []Enum) {
	return []Enum{
		Manual,
		Expiration,
	}
}
