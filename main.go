package main

func main() {
	scheduler := NewJobScheduler()

	scheduler.AddJob("job1", []string{})
	scheduler.AddJob("job2", []string{"job1"})
	scheduler.AddJob("job3", []string{"job1", "job2"})

	scheduler.DisplayJobs()
	scheduler.DisplayJobQueue()

	scheduler.ProcessJob()      // job1 icra edilecek
	scheduler.DisplayJobQueue() // Novbede job2 olacaq

	scheduler.ProcessJob()      // job2 icra edilecek
	scheduler.DisplayJobQueue() // Novbede job3 olacaq

	scheduler.ProcessJob()      // job3 icra edilecek
	scheduler.DisplayJobQueue() // Novbe bos olacaq
}
