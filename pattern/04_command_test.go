package pattern

import "testing"

func Test_Command(t *testing.T) {

	invoker := Invoker{}
	testObj := Receiver{
		fieldToChange: 10,
	}

	cmd1 := NewInc(&testObj)
	cmd2 := NewDouble(&testObj)
	cmd3 := NewInc(&testObj)
	cmd4 := NewDouble(&testObj)

	invoker.Execute(cmd1)
	if testObj.fieldToChange != 11 {
		t.FailNow()
	}

	invoker.Execute(cmd2)
	if testObj.fieldToChange != 22 {
		t.FailNow()
	}

	invoker.Execute(cmd3)
	if testObj.fieldToChange != 23 {
		t.FailNow()
	}

	invoker.Undo()
	if testObj.fieldToChange != 22 {
		t.FailNow()
	}

	invoker.Execute(cmd4)
	if testObj.fieldToChange != 44 {
		t.FailNow()
	}

	invoker.Undo()
	if testObj.fieldToChange != 22 {
		t.FailNow()
	}
}
