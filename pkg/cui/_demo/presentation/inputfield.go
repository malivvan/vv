package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/malivvan/vv/pkg/cui"
)

const inputField = `[green]package[white] main

[green]import[white] (
    [red]"strconv"[white]

    [red]"github.com/gdamore/tcell/v2"[white]
    [red]"github.com/malivvan/vv/pkg/cui"[white]
)

[green]func[white] [yellow]main[white]() {
    input := cui.[yellow]NewInputField[white]().
        [yellow]SetLabel[white]([red]"Enter a number: "[white]).
        [yellow]SetAcceptanceFunc[white](
            cui.InputFieldInteger,
        ).[yellow]SetDoneFunc[white]([yellow]func[white](key tcell.Key) {
            text := input.[yellow]GetText[white]()
            n, _ := strconv.[yellow]Atoi[white](text)
            [blue]// We have a number.[white]
        })
    cui.[yellow]NewApplication[white]().
        [yellow]SetRoot[white](input, true).
        [yellow]Run[white]()
}`

// InputField demonstrates the InputField.
func InputField(nextSlide func()) (title string, info string, content cui.Primitive) {
	input := cui.NewInputField()
	input.SetLabel("Enter a number: ")
	input.SetAcceptanceFunc(cui.InputFieldInteger)
	input.SetDoneFunc(func(key tcell.Key) {
		nextSlide()
	})
	return "InputField", "", Code(input, 30, 1, inputField)
}
