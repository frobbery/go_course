package hw06pipelineexecution

import "sync"

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
	readAllUnread(&stagesChannels)
	return stagesChannels[len(stages)-1]
}

func subscribe(in In, done In, stage Stage, stageChannel Bi) {
	go func() {
		wg := sync.WaitGroup{}
		outChannel := stage(in)
		for i := range outChannel {
			wg.Add(1)
			go func() {
				defer wg.Done()
				select {
				case <-done:
					return
				case stageChannel <- i:
				}
			}()
		}
		go func() {
			wg.Wait()
			close(stageChannel)
		}()
	}()
}

func readAllUnread(stages *[]Bi) {
	wg := sync.WaitGroup{}
	for stage := range *stages {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := range stage {
				_ = i
			}
		}()
	}
	wg.Wait()
}
