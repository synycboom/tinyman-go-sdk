package constants

const (
	// SwapFixedInput is a fixed-input swap type
	SwapFixedInput = "fixed-input"

	// SwapFixedOutput is a fixed-output swap type
	SwapFixedOutput = "fixed-output"
)

var (
	// SwapTypeMapping is a mapping for swap type
	SwapTypeMapping = map[string]string{
		SwapFixedInput:  "fi",
		SwapFixedOutput: "fo",
	}
)
