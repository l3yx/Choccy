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
		//把任务状态置为失败
		for _, runningAndQueuingTask := range runningAndQueuingTasks {
			runningAndQueuingTask.Status = -1
			runningAndQueuingTask.Logs += "任务执行中断\n"
			result = database.DB.Save(runningAndQueuingTask)
			if result.Error != nil {
				panic(result.Error.Error())
			}
		}
	} else {
		//恢复任务，相当于重新执行，要将参数重置，状态恢复到队列中
		for _, runningAndQueuingTask := range runningAndQueuingTasks {
			runningAndQueuingTask.Status = 0
			runningAndQueuingTask.Stage = 0
			//这个字段是中release模式中获取到的所有新版本（database模式只有一个版本），无论重置与否，都会被重新赋值
			runningAndQueuingTask.Versions = []string{}
			//这俩不要置空，这俩有值了说明扫描结果也已经出了，project表也会标记最后扫描版本，任务恢复只会扫没扫的部分
			//置空也会导致扫出来的结果关联不到对应的任务数据（现在结果通过id关联，不会通过AResults关联了）
			//runningAndQueuingTask.AnalyzedVersions = []string{}
			//runningAndQueuingTask.Results = []model.TaskResult{}
			runningAndQueuingTask.Logs += "任务已恢复\n"
			result = database.DB.Save(runningAndQueuingTask)
			if result.Error != nil {
				panic(result.Error.Error())
			}

			CH <- runningAndQueuingTask.ID
		}
	}
}

func InitTask() {
	CH = make(chan uint, 10000)
	recoveryTask()
}
