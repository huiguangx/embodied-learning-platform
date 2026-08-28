package repository

import (
	"errors"
	"sync"
	"time"

	"eip-platform/backend/internal/domain"
)

type Reservation struct { ID, ProjectID, ResourceGroupID string; GPU int; StartsAt, EndsAt time.Time }
type ResourceRequest struct { ResourceGroupID string `json:"resourceGroupId"`; GPU, CPU, MemoryMB int }
type ResourceCheck struct { Allowed bool `json:"allowed"`; ReasonCode string `json:"reasonCode,omitempty"`; Remaining ResourceRequest `json:"remaining"`; Alternatives []string `json:"alternatives"` }

var resourceStore = struct{ sync.Mutex; groups []domain.ResourceGroup; queues []domain.Queue; reservations []Reservation }{
	groups: []domain.ResourceGroup{{ID:"rg-cloud",ProjectID:"00000000-0000-0000-0000-000000000001",Name:"Default Cloud",GPUQuota:16,CPUQuota:256,MemoryQuotaMB:1048576,Concurrency:8},{ID:"rg-idc",ProjectID:"00000000-0000-0000-0000-000000000001",Name:"Default IDC",GPUQuota:8,CPUQuota:128,MemoryQuotaMB:524288,Concurrency:4}},
	queues: []domain.Queue{{ID:"queue-default",ResourceGroupID:"rg-cloud",Name:"default",Priority:50,Weight:1,Concurrency:8},{ID:"queue-priority",ResourceGroupID:"rg-cloud",Name:"priority",Priority:100,Weight:2,Concurrency:4}},
}

func ResourceGroups(projectID string) []domain.ResourceGroup { resourceStore.Lock(); defer resourceStore.Unlock(); var out []domain.ResourceGroup; for _,g:=range resourceStore.groups {if g.ProjectID==projectID {out=append(out,g)}}; return out }
func AddResourceGroup(group domain.ResourceGroup) error { resourceStore.Lock(); defer resourceStore.Unlock(); if group.ID==""||group.ProjectID==""||group.GPUQuota<0 {return errors.New("invalid resource group")}; resourceStore.groups=append(resourceStore.groups,group); return nil }
func Queues(projectID string) []domain.Queue { groups:=ResourceGroups(projectID); allowed:=map[string]bool{}; for _,g:=range groups {allowed[g.ID]=true}; resourceStore.Lock(); defer resourceStore.Unlock(); var out []domain.Queue; for _,q:=range resourceStore.queues {if allowed[q.ResourceGroupID] {out=append(out,q)}}; return out }
func AddQueue(projectID string, queue domain.Queue) error { for _,g:=range ResourceGroups(projectID) {if g.ID==queue.ResourceGroupID {resourceStore.Lock(); resourceStore.queues=append(resourceStore.queues,queue); resourceStore.Unlock(); return nil}}; return errors.New("resource group not found") }
func Reservations(projectID string) []Reservation { resourceStore.Lock(); defer resourceStore.Unlock(); var out []Reservation; for _,v:=range resourceStore.reservations {if v.ProjectID==projectID {out=append(out,v)}}; return out }
func AddReservation(v Reservation) error { resourceStore.Lock(); defer resourceStore.Unlock(); if !v.StartsAt.Before(v.EndsAt) {return errors.New("invalid reservation window")}; for _,x:=range resourceStore.reservations {if x.ProjectID==v.ProjectID&&x.ResourceGroupID==v.ResourceGroupID&&v.StartsAt.Before(x.EndsAt)&&x.StartsAt.Before(v.EndsAt){return errors.New("reservation overlaps")}}; resourceStore.reservations=append(resourceStore.reservations,v); return nil }
func CheckResources(projectID string, request ResourceRequest) ResourceCheck { for _,g:=range ResourceGroups(projectID) {if g.ID==request.ResourceGroupID {remaining:=ResourceRequest{ResourceGroupID:g.ID,GPU:g.GPUQuota,CPU:g.CPUQuota,MemoryMB:g.MemoryQuotaMB}; ok:=request.GPU<=g.GPUQuota&&request.CPU<=g.CPUQuota&&request.MemoryMB<=g.MemoryQuotaMB; code:=""; if !ok {code="QUOTA_EXCEEDED"}; return ResourceCheck{Allowed:ok,ReasonCode:code,Remaining:remaining,Alternatives:[]string{"reduce resources","choose another resource group"}}}}; return ResourceCheck{Allowed:false,ReasonCode:"RESOURCE_GROUP_NOT_FOUND",Alternatives:[]string{}} }
