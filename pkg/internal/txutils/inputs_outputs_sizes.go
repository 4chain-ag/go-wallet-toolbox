package txutils

const P2PKHUnlockingScriptLength = 107

var P2PKHOutputSize = TransactionOutputSize(25)
var P2PKHEstimatedInputSize = TransactionInputSize(P2PKHUnlockingScriptLength)
