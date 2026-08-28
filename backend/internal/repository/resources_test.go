package repository

import (
	"testing"
	"time"
)

func TestResourceCheckRejectsQuotaExceeded(t *testing.T){r:=CheckResources("00000000-0000-0000-0000-000000000001",ResourceRequest{ResourceGroupID:"rg-cloud",GPU:17});if r.Allowed||r.ReasonCode!="QUOTA_EXCEEDED"{t.Fatalf("unexpected result: %+v",r)}}
func TestReservationRejectsOverlap(t *testing.T){start:=time.Date(2026,8,28,1,0,0,0,time.UTC); first:=Reservation{ID:"r1",ProjectID:"test-overlap",ResourceGroupID:"rg",StartsAt:start,EndsAt:start.Add(time.Hour)};if err:=AddReservation(first);err!=nil{t.Fatal(err)};if err:=AddReservation(Reservation{ID:"r2",ProjectID:"test-overlap",ResourceGroupID:"rg",StartsAt:start.Add(time.Minute),EndsAt:start.Add(2*time.Hour)});err==nil{t.Fatal("expected overlap error")}}
