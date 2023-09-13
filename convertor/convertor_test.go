package convertor

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

type Event struct {
	Employee      People
	LaborActivity LaborActivity
	CustomFields  Addition
	staffField    string
}

type People struct {
	Name   string `ref:"Name"`
	Age    int    `ref:"Age"`
	IsMale bool   `ref:"IsMale"`
}

type LaborActivity struct {
	Experience     int      `ref:"Experience"`
	PastJobs       []string `ref:"PastJobs"`
	CurentPost     Post
	PastExperience []int `ref:"PastExperience"`
}

type Post struct {
	JobTitle *string `ref:"JobTitle"`
	Salary   float32 `ref:"Salary"`
}

type Addition struct {
	Coments []string `ref:"Coments"`
}

func TestUnmarshalMap(t *testing.T) {
	jobTitle := "разработчик"

	tests := []struct {
		name     string
		wantData any
		data     string
		wantErr  error
	}{
		{
			name: "succes",
			wantData: &Event{
				Employee: People{
					Name:   "Максим",
					Age:    15,
					IsMale: true,
				},
				LaborActivity: LaborActivity{
					Experience: 10,
					PastJobs:   []string{"mts", "вмф"},
					CurentPost: Post{
						JobTitle: &jobTitle,
						Salary:   500.23,
					},
					PastExperience: []int{11, 2},
				},
				CustomFields: Addition{
					Coments: []string{"comment_1", "comment_2"},
				},
			},
			data:    `{"Employee":{"Name":"Максим","Age":15,"IsMale":true},"LaborActivity":{"Experience":10,"PastJobs":["mts","вмф"],"PastExperience":[11,2],"CurentPost":{"JobTitle":"разработчик","Salary":500.23}},"CustomFields":{"Coments":["comment_1","comment_2"]}}`,
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var dataRaw = make(map[string]interface{})
			var data = make(map[string]interface{})
			got := &Event{}
			err := json.Unmarshal([]byte(tt.data), &dataRaw)
			require.NoError(t, err, "Unmarshal in data error")

			DimensionalMap(dataRaw, data)

			err = UnmarshalMap(got, data)

			require.Equal(t, tt.wantErr, err, "UnmarshalMap error")
			require.Equal(t, tt.wantData, got, "not equal data")
		})
	}
}
