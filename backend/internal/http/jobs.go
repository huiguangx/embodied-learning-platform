package http

import("net/http";"fmt";"time";"github.com/gin-gonic/gin";"eip-platform/backend/internal/auth";"eip-platform/backend/internal/domain";"eip-platform/backend/internal/repository";"eip-platform/backend/internal/scheduler")
var localScheduler scheduler.Local
func listJobs(c *gin.Context){items:=repository.ListJobs(c.Param("projectId"));c.JSON(200,gin.H{"items":items,"nextCursor":"","total":len(items)})}
func validateJob(c *gin.Context){var j domain.TrainingJob;if c.ShouldBindJSON(&j)!=nil||j.ImageDigest==""||j.CodeVersion==""||j.DatasetVersion==""{Error(c,400,"BAD_INPUT","image, code and dataset versions are required",nil);return};c.JSON(200,gin.H{"valid":true,"job":j})}
func createJob(c *gin.Context){var j domain.TrainingJob;if c.ShouldBindJSON(&j)!=nil{Error(c,400,"BAD_INPUT","invalid training job",nil);return};j.ID=stringID();j.ProjectID=c.Param("projectId");created,err:=repository.CreateJob(j,c.GetHeader("Idempotency-Key"));if err!=nil{Error(c,400,"BAD_INPUT",err.Error(),nil);return};if err=localScheduler.Submit(c.Request.Context(),created);err!=nil{Error(c,500,"SCHEDULER_ERROR",err.Error(),nil);return};created,_=repository.GetJob(created.ID);c.JSON(http.StatusCreated,created)}
func cancelJob(c *gin.Context){j,ok:=repository.GetJob(c.Param("id"));if !ok||j.ProjectID!=c.Param("projectId"){Error(c,404,"JOB_NOT_FOUND","job not found",nil);return};if err:=localScheduler.Cancel(c.Request.Context(),j.ID);err!=nil{Error(c,409,"INVALID_STATE",err.Error(),nil);return};j,_=repository.GetJob(j.ID);c.JSON(200,j)}
func stringID()string{return fmt.Sprintf("job-%d",time.Now().UnixNano())}
func registerJobRoutes(g *gin.RouterGroup){g.GET("/training-jobs",listJobs);g.POST("/training-jobs:validate",validateJob);g.POST("/training-jobs",auth.RequireWrite(),createJob);g.POST("/training-jobs/:id:cancel",auth.RequireWrite(),cancelJob)}
