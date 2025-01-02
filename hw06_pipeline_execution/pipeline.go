package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	currStageChannel := subscribe(in, done, stages[0])
	for i := 1; i < len(stages); i++ {
		currStageChannel = subscribe(currStageChannel, done, stages[i])
	}
	return currStageChannel
}

func subscribe(in In, done In, stage Stage) Out {
	stageChannel := make(Bi)
	go func() {
		for i := range in {
			select {
			case <-done:
				//nolint
				defer readAllUnread(stageChannel)
			default:
				stageChannel <- i
			}
		}
		close(stageChannel)
	}()
	return stage(stageChannel)
}

func readAllUnread(channel Bi) {
	go func() {
		for i := range channel {
			_ = i
		}
	}()
}
