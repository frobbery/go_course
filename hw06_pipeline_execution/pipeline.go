package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	currStageChannel := in
	for i := 0; i < len(stages); i++ {
		tmpStageChannel := make(Bi)
		go func(currStageChannel In) {
			defer close(tmpStageChannel)
			for i := range currStageChannel {
				select {
				case <-done:
					go func() {
						for i := range currStageChannel {
							_ = i
						}
					}()
					return
				default:
					tmpStageChannel <- i
				}
			}
		}(currStageChannel)
		currStageChannel = stages[i](tmpStageChannel)
	}
	return currStageChannel
}
