package main

// In memory storage for skier data

type Skier struct {
	ResortID  int32  `json:"resortID"`
	SeasonID  string `json:"seasonID"`
	DayID     string `json:"dayID"`
	SkierID   int32  `json:"skierID"`
	TotalVert int32  `json:"totalVert"`
}

type TotalVertResponse struct {
	TotalVert int32 `json:"totalVert"`
}

// A very inefficient memory store.
// Also not thread safe but who cares.
// A list of Skiers.
type InMemoryStore struct {
	data []Skier
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		// Add 1 dummy Skier
		data: []Skier{
			{
				ResortID:  1,
				SeasonID:  "2023",
				DayID:     "5",
				SkierID:   123,
				TotalVert: 500,
			},
		},
	}
}

func (store *InMemoryStore) AddSkier(skier Skier) {
	store.data = append(store.data, skier)
}

// Get the total vertical for a skier for the specific day.
func (store *InMemoryStore) GetSkierDayVert(resortID int32, seasonID, dayID string, skierID int32) int32 {
	var totalVert int32 = 0
	for _, skier := range store.data {
		if skier.SkierID == skierID && skier.ResortID == resortID && skier.SeasonID == seasonID && skier.DayID == dayID {
			totalVert += skier.TotalVert
		}
	}
	return totalVert
}

// Get the total vertical for a skier. resortID and seasonID are optional.
func (store *InMemoryStore) GetSkierTotalVert(skierID int32, resortID *int32, seasonID *string) int32 {
	var totalVert int32 = 0
	for _, skier := range store.data {
		if skier.SkierID == skierID && (resortID == nil || skier.ResortID == *resortID) && (seasonID == nil || skier.SeasonID == *seasonID) {
			totalVert += skier.TotalVert
		}
	}
	return totalVert
}
