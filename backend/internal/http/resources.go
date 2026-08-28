package http

import (
	"net/http"
	"time"
	"github.com/gin-gonic/gin"
	"eip-platform/backend/internal/auth"
	"eip-platform/backend/internal/domain"
	"eip-platform/backend/internal/repository"
)

func listResourceGroups(c *gin.Context){items:=repository.ResourceGroups(c.Param("projectId")); c.JSON(200,gin.H{"items":items,"nextCursor":"","total":len(items)})}
func createResourceGroup(c *gin.Context){var v domain.ResourceGroup; if c.ShouldBindJSON(&v)!=nil {Error(c,400,"BAD_INPUT","invalid resource group",nil);return}; v.ProjectID=c.Param("projectId"); if repository.AddResourceGroup(v)!=nil {Error(c,400,"BAD_INPUT","invalid resource group",nil);return}; c.JSON(201,v)}
func listQueues(c *gin.Context){items:=repository.Queues(c.Param("projectId")); c.JSON(200,gin.H{"items":items,"nextCursor":"","total":len(items)})}
func createQueue(c *gin.Context){var v domain.Queue; if c.ShouldBindJSON(&v)!=nil||repository.AddQueue(c.Param("projectId"),v)!=nil {Error(c,400,"BAD_INPUT","invalid queue",nil);return}; c.JSON(201,v)}
func listReservations(c *gin.Context){items:=repository.Reservations(c.Param("projectId")); c.JSON(200,gin.H{"items":items,"nextCursor":"","total":len(items)})}
func createReservation(c *gin.Context){var body struct{ID,ResourceGroupID,StartsAt,EndsAt string;GPU int}; if c.ShouldBindJSON(&body)!=nil {Error(c,400,"BAD_INPUT","invalid reservation",nil);return}; start,e1:=time.Parse(time.RFC3339,body.StartsAt); end,e2:=time.Parse(time.RFC3339,body.EndsAt); if e1!=nil||e2!=nil||repository.AddReservation(repository.Reservation{ID:body.ID,ProjectID:c.Param("projectId"),ResourceGroupID:body.ResourceGroupID,GPU:body.GPU,StartsAt:start,EndsAt:end})!=nil {Error(c,409,"RESERVATION_CONFLICT","invalid or overlapping reservation",nil);return}; c.Status(201)}
func resourceCheck(c *gin.Context){var v repository.ResourceRequest;if c.ShouldBindJSON(&v)!=nil {Error(c,400,"BAD_INPUT","invalid resource request",nil);return}; result:=repository.CheckResources(c.Param("projectId"),v); status:=http.StatusOK;if !result.Allowed {status=http.StatusUnprocessableEntity};c.JSON(status,result)}

func registerResourceRoutes(group *gin.RouterGroup){group.GET("/resource-groups",listResourceGroups);group.POST("/resource-groups",auth.RequireWrite(),createResourceGroup);group.GET("/queues",listQueues);group.POST("/queues",auth.RequireWrite(),createQueue);group.GET("/reservations",listReservations);group.POST("/reservations",auth.RequireWrite(),createReservation);group.POST("/resource-checks",resourceCheck)}
