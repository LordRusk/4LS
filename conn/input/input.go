package input

import (
	"fmt"

	errs "github.com/pkg/errors"

	"github.com/AlecAivazis/survey/v2"
)

type Inputter struct {
}

func NewInputter() *Inputter {
	return &Inputter{}
}

var question = survey.Select{
	Message: "Choose a color:",
	Options: []string{"red", "blue", "green"},
	// can pass a validator directly
	Validate: survey.Required}

func (i *Inputter) Test() error {
	var answer string
	if err := survey.AskOne(question, &answer); err != nil {
		return errs.Wrap(err, "Failed to cycle questions")
	}

	// fmt.Printf("%+v\n", answer)
	fmt.Printf("%s\n", answer)
	return nil
}
