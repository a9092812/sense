package pipeline

type JobChannel struct {
	ch chan []byte
}

func NewJobChannel(size int) *JobChannel {

	return &JobChannel{
		ch: make(chan []byte, size),
	}
}

func (j *JobChannel) Push(job []byte) {

	j.ch <- job
}

func (j *JobChannel) Channel() <-chan []byte {

	return j.ch
}

func (j *JobChannel) SendChannel() chan<- []byte {

	return j.ch
}

func (j *JobChannel) Close() {

	close(j.ch)
}
