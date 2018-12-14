import app.app as app
import signal
import time
from celery import group


class DaemonKiller:
    kill_now = False

    def __init__(self):
        signal.signal(signal.SIGINT, self.exit_gracefully)
        signal.signal(signal.SIGTERM, self.exit_gracefully)

    def exit_gracefully(self, signum, frame):
        self.kill_now = True


def daemon_crawl():
    daemon_killer = DaemonKiller()
    task_overlap = 2  # tolerence for queue to allow next crawl to start
    inspector = app.celery_app.control.inspect()  # used to inspect reserved jobs in celery worker

    while True:
        worker_queue = inspector.reserved()  # worker_queue is a list of dicts

        # grab all queued ("reserved") tasks for this
        try:
            queued_tasks = worker_queue.get(next(iter(worker_queue)))
        except:
            # Will except on first iteration since queue will be empty
            queued_tasks = []

        num_tasks = len(queued_tasks)
        # print(num_tasks)

        job = group([app.homepage_crawl()])
        result = job.apply_async()
        result.join()

        # If there are less than 3 tasks in the queue, queue up next crawl
        if num_tasks < task_overlap:
            # app.crawl_everything()
            pass

        # Breakout if we get signal
        if daemon_killer.kill_now:
            break

        time.sleep(20)  # Sleep 120sec


if __name__ == '__main__':
    # Start Daemon
    daemon_crawl()
