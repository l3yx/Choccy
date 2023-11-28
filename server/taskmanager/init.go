package taskmanager

import (
	"choccy/server/database"
	"choccy/server/database/model"
)

var CH chan uint

func recoveryTask() {
	var setting model.Setting
	result := database.DB.Take(&setting)
	if result.Error != nil {
		panic(result.Error.Error())
	}

	var runningAndQueuingTasks []model.Task
	result = database.DB.Where("status in ?", []int{1, 0}).Find(&runningAndQueuingTasks)
	if result.Error != nil {
		panic(result.Error.Error())
	}

	if !setting.AutoRecoveryTask {
		//Set Task Status to Failed
		for _, runningAndQueuingTask := range runningAndQueuingTasks {
			runningAndQueuingTask.Status = -1
			runningAndQueuingTask.Logs += "Mission Execution Interruption\n"
			result = database.DB.Save(runningAndQueuingTask)
			if result.Error != nil {
				panic(result.Error.Error())
			}
		}
	} else {
		//Recovery task, equivalent to re-executing, to reset the parameters, restore the status to the queue
		for _, runningAndQueuingTask := range runningAndQueuingTasks {
			runningAndQueuingTask.Status = 0
			runningAndQueuingTask.Stage = 0
			//This field is all new versions obtained in release mode (only one version in database mode) and will be reassigned regardless of whether it is reset or not
			runningAndQueuingTask.Versions = []string{}
			//These two should not be empty，These two values show that the scan results have also been out，Project table will also mark the final scan version，Task recovery will only sweep the parts that were not scanned
			//置Null also causes the scanned results to be not associated with the corresponding task data (now the results are associated by id, not by AR Results)
			//runningAndQueuingTask.AnalyzedVersions = []string{}
			//runningAndQueuingTask.Results = []model.TaskResult{}
			runningAndQueuingTask.Logs += "Mission has resumed.\n"
			result = database.DB.Save(runningAndQueuingTask)
			if result.Error != nil {
				panic(result.Error.Error())
			}

			CH <- runningAndQueuingTask.ID
		}
	}
}

func InitTask() {
	CH = make(chan uint, 1000)
	recoveryTask()
}
