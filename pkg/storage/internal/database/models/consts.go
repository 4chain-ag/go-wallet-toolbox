package models

// EstimatedInputSizeForP2PKH is the estimated size increase when adding and unlocking P2PKH input to transaction.
// 32 bytes txID
// + 4 bytes vout index
// + 1 byte script length
// + 107 bytes script pub key
// + 4 bytes nSequence
const EstimatedInputSizeForP2PKH = 148
