package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	// Place your code here.

	doWork := func(in In, done In, s Stage) Out {
		out := make(Bi)
		res := s(in)
		go func() {
			defer close(out)
			for {
				select {
				case <-done:
					return
				case v, closeCh := <-res:
					if !closeCh {
						return
					}
					out <- v
				}
			}
		}()
		return out
	}
	out := make(Out)
	for _, s := range stages {
		out = doWork(in, done, s)
		in = out
	}
	return out
}
