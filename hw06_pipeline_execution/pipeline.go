package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	stagesChannels := make([]Bi, len(stages))
	for i := 0; i < len(stages); i++ {
		stagesChannels[i] = make(Bi)
	}
	subscribe(in, done, stages[0], stagesChannels[0])
	for i := 1; i < len(stages); i++ {
		subscribe(stagesChannels[i-1], done, stages[i], stagesChannels[i])
	}
	return stagesChannels[len(stages)-1]
}

func subscribe(in In, done In, stage Stage, stageChannel Bi) {
	outChannel := stage(in)
	go func() {
		defer closeChannel(stageChannel)
		for i := range outChannel {
			select {
			case <-done:
				return
			case stageChannel <- i:
			}
		}
	}()
}

func closeChannel(channel Bi) {
	go func() {
		for i := range channel {
			_ = i
		}
	}()
	close(channel)
}
