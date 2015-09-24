package multithread_du

import "syscall"
import "sync"
import "container/list"
import "os"

// A queue that keeps track of the sum of (jobs enqueued + workers working)
// Also sums the values from each task to report total file size.
type JobQueue struct {
	q list.List
	l sync.Mutex

	total int64

	// needed to avoid double-counting hardlinks
	seen_inodes map[uint64]bool

	master_counter int
}

func NewJobQueue() JobQueue {
	var jq JobQueue
	jq.seen_inodes = make(map[uint64]bool)
	return jq
}

// Push into queue
func (q *JobQueue) EnqueueTask(data interface{}) {
	q.l.Lock()
	q.q.PushBack(data)
	q.master_counter += 1
	q.l.Unlock()
}

// Pop from the queue (and assume that one starts working on task)
func (q *JobQueue) AssignTask() interface{} {
	q.l.Lock()
	data := q.q.Remove(q.q.Front())
	q.l.Unlock()
	return data
}

func (q *JobQueue) SignalFinishedTask(inode_num uint64, size int64) {
	q.l.Lock()
	if !q.seen_inodes[inode_num] {
		q.seen_inodes[inode_num] = true
		q.total += size
	}
	q.master_counter -= 1
	q.l.Unlock()
}

func (q *JobQueue) IsFinished() bool {
	q.l.Lock()
	is_done := (q.master_counter == 0)
	q.l.Unlock()
	return is_done
}

// number of tasks enqueued
func (q *JobQueue) Len() int {
	return q.q.Len()
}

// for reference, see http://lxr.free-electrons.com/source/include/uapi/linux/stat.h#L21
func isDirectory(mode uint32) bool {
	var s_ifmt uint32 = 00170000
	var s_ifdir uint32 = 0040000

	return (((mode) & s_ifmt) == s_ifdir)
}

func ProcessFile(filename string, q *JobQueue) (uint64, int64) {
	var x syscall.Stat_t
	err := syscall.Lstat(filename, &x)

	if err == nil {
		nblocks := x.Blocks
		var mode uint32 = uint32(x.Mode)

		if isDirectory(mode) {
			file, _ := os.Open(filename)
			glob_results, _ := file.Readdirnames(-1)
			file.Close()

			for _, elt := range glob_results {
				q.EnqueueTask(filename + "/" + elt)
			}
		}

		inode_id := x.Ino

		return inode_id, nblocks
	}
	return 0, 0
}

func TotalFileSize(root string) int64 {
	c := make(chan string)

	job_queue := NewJobQueue()
	job_queue.EnqueueTask(root)

	for i := 0; i < 16; i++ {
		go func() {
			for nf := range c {
				inode_id, counts := ProcessFile(nf, &job_queue)
				job_queue.SignalFinishedTask(inode_id, counts)
			}
		}()
	}

	for !(job_queue.IsFinished()) {
		if job_queue.Len() > 0 {
			// the queue is used in addition to the channels because we want an infinite buffer
			nextFile, _ := job_queue.AssignTask().(string)
			c <- nextFile
		}
	}
	close(c)

	return job_queue.total
}
