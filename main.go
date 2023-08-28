package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/creditdb/go-creditdb"
	"github.com/google/uuid"
)

type Mood string

const (
	Happy Mood = "happy"
	Sad   Mood = "sad"
)

type Person struct {
	ID    string    `json:"id"`
	Name  string    `json:"name"`
	Age   int       `json:"age"`
	DOB   time.Time `json:"dob"`
	Moods []Mood    `json:"moods"`
}

func (p *Person) MarshalBinary() ([]byte, error) {
	return json.Marshal(p)
}

func (p *Person) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, p)
}

func main() {
	c := creditdb.NewClient().WithPage(2)
	ctx := context.Background()
	defer c.Close(ctx)

	person := Person{
		ID:    uuid.New().String(),
		Name:  "John Doe",
		Age:   30,
		DOB:   time.Now(),
		Moods: []Mood{Happy, Sad, Happy, Happy},
	}

	personJSON, err := person.MarshalBinary()
	if err != nil {
		panic(err)
	}

	personLine := creditdb.Line{
		Key:   person.Name,
		Value: string(personJSON),
	}

	// SetLine
	if err := c.SetLine(ctx, personLine.Key, personLine.Value); err != nil {
		panic(err)
	}

	// GetLine
	getPersonLine, err := c.GetLine(ctx, personLine.Key)
	if err != nil {
		panic(err)
	}

	getPerson := Person{}
	if err := getPerson.UnmarshalBinary([]byte(getPersonLine.Value)); err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", getPerson)
	fmt.Printf("%+v\n", getPerson.Name)
	fmt.Printf("%+v\n", getPerson.Age)
	fmt.Printf("%+v\n", getPerson.DOB)

	/*

		$ go run main.go
		{ID:f7534b30-94ab-4659-80c4-d412c27c36fb Name:John Doe Age:30 DOB:2023-08-28 22:40:21.5176114 +0200 CEST Moods:[happy sad happy happy]}
		John Doe
		30
		2023-08-28 22:40:21.5176114 +0200 CEST

	*/

	fmt.Println("=====================================")

	// Set and Array of Objects

	people := []Person{
		{
			ID:    uuid.New().String(),
			Name:  "John Doe",
			Age:   30,
			DOB:   time.Now(),
			Moods: []Mood{Happy, Sad, Happy, Happy},
		},
		{
			ID:    uuid.New().String(),
			Name:  "Jane Doe",
			Age:   25,
			DOB:   time.Now(),
			Moods: []Mood{Happy, Happy, Happy, Happy},
		},
	}

	peopleJSON, err := json.Marshal(people)
	if err != nil {
		panic(err)
	}

	peopleLine := creditdb.Line{
		Key:   "people",
		Value: string(peopleJSON),
	}

	// SetLine
	if err := c.SetLine(ctx, peopleLine.Key, peopleLine.Value); err != nil {
		panic(err)
	}

	// GetLine
	getPeopleLine, err := c.GetLine(ctx, peopleLine.Key)
	if err != nil {
		panic(err)
	}

	getPeople := []Person{}
	if err := json.Unmarshal([]byte(getPeopleLine.Value), &getPeople); err != nil {
		panic(err)
	}

	for _, person := range getPeople {
		fmt.Printf("%+v\n", person)
		fmt.Printf("%+v\n", person.ID)
		fmt.Printf("%+v\n", person.Name)
		fmt.Printf("%+v\n", person.Age)
		fmt.Printf("%+v\n", person.DOB)
		fmt.Println()
	}

	/*

		$ go run main.go
		{ID:c51bfa89-0db8-4be1-a40d-0fe059fa1992 Name:John Doe Age:30 DOB:2023-08-28 22:40:21.5250386 +0200 CEST Moods:[happy sad happy happy]}
		c51bfa89-0db8-4be1-a40d-0fe059fa1992
		John Doe
		30
		2023-08-28 22:40:21.5250386 +0200 CEST

		{ID:764b797d-9e19-4afd-877e-080a72754af3 Name:Jane Doe Age:25 DOB:2023-08-28 22:40:21.5250386 +0200 CEST Moods:[happy happy happy happy]}
		764b797d-9e19-4afd-877e-080a72754af3
		Jane Doe
		25
		2023-08-28 22:40:21.5250386 +0200 CEST

	*/
}
