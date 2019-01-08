from flask import Flask, render_template
from celery import Celery
from celery.exceptions import SoftTimeLimitExceeded
from celery.schedules import crontab
from app.taskdir.tasks import crawl_task, weekly_report, alch_metrics, home_crawl
from app.alchemy import *
import datetime


def make_celery(app):
    celery = Celery(app.import_name, backend=app.config['CELERY_BACKEND'],
                    broker=app.config['CELERY_BROKER_URL'])
    celery.conf.update(app.config)

    task_base = celery.Task

    class ContextTask(task_base):
        abstract = True

        def __call__(self, *args, **kwargs):
            with app.app_context():
                return task_base.__call__(self, *args, **kwargs)
    celery.Task = ContextTask
    return celery


app = Flask(__name__)
app.config['CELERY_BACKEND'] = "redis://redis:6379/0"
app.config['CELERY_BROKER_URL'] = "redis://redis:6379/0"
app.config['CELERY_TIMEZONE'] = 'UTC'
app.config['CELERYBEAT_SCHEDULE'] = {
    # 'crawl_everything': {
    #     'task': 'crawl_everything',
    #     'schedule': crontab(hour='*/1')
    # },
    'weekly_email': {
        'task': 'weekly_email',
        'schedule': crontab(minute=0, hour=8, day_of_week="mon")
    },
    'monthly_purge': {
        'task': 'clear_old_data',
        'schedule': crontab(day_of_month=1)
    },
}

celery_app = make_celery(app)


@celery_app.task(name='crawl_domain', time_limit=600, soft_time_limit=500)
def crawl_domain(domain, d_filter, is_off, crawl_id):
    try:
        res = crawl_task(domain, d_filter, is_off, crawl_id)
        res_out = res.wait(interval=0.5)
    except SoftTimeLimitExceeded:
        alch_metrics(domain, crawl_id)
        return


@celery_app.task(name='homepage_crawl')
def homepage_crawl():
    print('home crawl start', '')
    home_crawl()


@celery_app.task(name='weekly_email')
def weekly_email():
    weekly_report()


@celery_app.task(name='clear_old_data')
def clear_old_data():
    too_old = datetime.datetime.today() - datetime.timedelta(weeks=4)
    session = loadSession()
    # Delete old from Page
    session.query(Page).filter(Page.datetime < too_old).delete()
    # Delete old from Crawl
    session.query(Crawl).filter(Crawl.start_time < too_old).delete()
    # Commit/Close
    session.commit()
    session.close()


@celery_app.task(name='crawl_everything')
def crawl_everything():
    task_overlap = 2  # tolerence for queue to allow next crawl to start
    inspector = celery_app.control.inspect()  # used to inspect reserved jobs in celery worker
    worker_queue = inspector.reserved()  # worker_queue is a list of dicts

    # grab all queued ("reserved") tasks for this
    try:
        queued_tasks = worker_queue.get(next(iter(worker_queue)))
    except:
        # Will except on first iteration since queue will be empty
        queued_tasks = []

    num_tasks = len(queued_tasks)
    # print(num_tasks)

    # If there are less than threshold tasks in the queue, queue up next crawl
    if num_tasks < task_overlap:
        # New DB session for each iteration
        session = loadSession()
        res = session.query(Status.domain, Status.filter, Status.is_offline)
        domain_list = [[row.domain, row.filter, row.is_offline] for row in res.all()]

        # Reset log start time, and log start time of this iteration of crawl
        log_start = datetime.datetime.now().strftime("%Y-%m-%d %H:%M:%S")
        # Insert new crawl into crawl table
        try:
            # SQLAlchemy: insert new crawl to crawl table
            new_crawl = Crawl(start_time=log_start)
            session.add(new_crawl)
            session.commit()

        except Exception as e:
            print('Unable to insert new crawl ERROR: ', e)
            return

        # Get the id for the new crawl, *(it is the first index of the first returned value)
        current_crawl_id = session.query(Crawl.id).filter(Crawl.start_time == log_start).first()[0]

        # Then we queue up all domains in the list
        for domain, d_filter, is_off in domain_list:
            # if domain not in queued_tasks:
            # task_list.append(crawl_task(domain, https_status, current_crawl_id))  # Append domain crawl to list
            celery_app.send_task('crawl_domain', args=[domain, d_filter, is_off, current_crawl_id], kwargs={})

        # Commit/Close session for iteration
        session.commit()
        session.close()


@app.route('/')
def home():
    # Get list of domains for cards
    session = loadSession()
    res = session.query(Status.domain, Status.is_offline)
    domain_list = [[row.domain, row.is_offline] for row in res.all()]
    session.close()

    return render_template('index.html', domain_list=domain_list,
                           time=datetime.datetime.now().strftime("%Y-%m-%d %H:%M:%S"))
