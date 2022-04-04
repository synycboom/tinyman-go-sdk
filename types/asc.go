package types

// Logic represents logic
type Logic struct {
	Bytecode  string     `json:"bytecode"`
	Address   string     `json:"address"`
	Size      int        `json:"size"`
	Variables []Variable `json:"variables"`
	Source    string     `json:"source"`
}

// PoolLogicSig represents a pool logic signature
type PoolLogicSig struct {
	Type  string `json:"type"`
	Logic Logic  `json:"logic"`
	Name  string `json:"name"`
}

// Variable represents a variable
type Variable struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	Index  int    `json:"index"`
	Length int    `json:"length"`
}

// ApprovalProgram represents an approval program
type ApprovalProgram struct {
	Bytecode  string     `json:"bytecode"`
	Address   string     `json:"address"`
	Size      int        `json:"size"`
	Variables []Variable `json:"variables"`
	Source    string     `json:"source"`
}

// ClearProgram represents a clear program
type ClearProgram struct {
	Bytecode  string     `json:"bytecode"`
	Address   string     `json:"address"`
	Size      int        `json:"size"`
	Variables []Variable `json:"variables"`
	Source    string     `json:"source"`
}

// GlobalStateSchema represents a global state schema
type GlobalStateSchema struct {
	NumUints      int `json:"num_uints"`
	NumByteSlices int `json:"num_byte_slices"`
}

// LocalStateSchema represents a local state schema
type LocalStateSchema struct {
	NumUints      int `json:"num_uints"`
	NumByteSlices int `json:"num_byte_slices"`
}

// ValidatorApp represents a validator app
type ValidatorApp struct {
	Type              string            `json:"type"`
	ApprovalProgram   ApprovalProgram   `json:"approval_program"`
	ClearProgram      ClearProgram      `json:"clear_program"`
	GlobalStateSchema GlobalStateSchema `json:"global_state_schema"`
	LocalStateSchema  LocalStateSchema  `json:"local_state_schema"`
	Name              string            `json:"name"`
}

// Contracts represents contracts
type Contracts struct {
	PoolLogicSig PoolLogicSig `json:"pool_logicsig"`
	ValidatorApp ValidatorApp `json:"validator_app"`
}

// ASC represents Algorand Smart Contract
type ASC struct {
	Repo      string    `json:"repo"`
	Ref       string    `json:"ref"`
	Contracts Contracts `json:"contracts"`
}
