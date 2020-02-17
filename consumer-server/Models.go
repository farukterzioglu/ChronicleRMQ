package consumerserver

type Transaction struct {
	TransactionHash string `json:"transaction_hash"`
	BlockHeight     uint64 `json:"block_height"`
	Trace           Trace
}

type Trace struct {
	Status       string        `json:"status"`
	ActionTraces []ActionTrace `json:"action_traces"`
}

type ActionTrace struct {
}
