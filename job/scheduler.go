package job

import (
	"github.com/vmware/harbor/dao"
	"github.com/vmware/harbor/models"
	"github.com/vmware/harbor/utils/log"
)

var JobQueue chan int64 = make(chan int64)

func Schedule(jobID int64) {
	JobQueue <- jobID
}

func HandleRepJob(id int64) {
	sm := &JobSM{JobID: id}
	err := sm.Init()
	if err != nil {
		log.Errorf("Failed to initialize statemachine, error: %v")
		err2 := dao.UpdateRepJobStatus(id, models.JobError)
		if err2 != nil {
			log.Errorf("Failed to update job status to ERROR, error:%v", err2)
		}
		return
	}
	if sm.Parms.Enabled == 0 {
		log.Debugf("The policy of job:%d is disabled, will cancel the job")
		_ = dao.UpdateRepJobStatus(id, models.JobCanceled)
	} else {
		sm.Start(models.JobRunning)
	}
}
