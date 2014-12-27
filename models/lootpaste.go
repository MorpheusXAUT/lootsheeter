// lootpaste
package models

type LootPaste struct {
	ID        int64
	FleetID   int64
	PastedBy  int64
	RawPaste  string
	Value     float64
	PasteType LootPasteType
}

func NewLootPaste(id int64, fleet int64, pasted int64, raw string, value float64, pasteType LootPasteType) *LootPaste {
	paste := &LootPaste{
		ID:        id,
		FleetID:   fleet,
		PastedBy:  pasted,
		RawPaste:  raw,
		Value:     value,
		PasteType: pasteType,
	}

	return paste
}

type LootPasteType int

const (
	LootPasteTypeUnknown LootPasteType = 1 << iota
	LootPasteTypeProfit
	LootPasteTypeLoss
)

func (t LootPasteType) String() string {
	switch t {
	case LootPasteTypeUnknown:
		return "U"
	case LootPasteTypeProfit:
		return "P"
	case LootPasteTypeLoss:
		return "L"
	default:
		return "U"
	}
}
