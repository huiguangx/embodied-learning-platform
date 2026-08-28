package domain
import "testing"
func TestTransitionRejectsRegression(t *testing.T){if _,err:=Transition(JobSucceeded,"start");err==nil{t.Fatal("expected invalid transition")}}
