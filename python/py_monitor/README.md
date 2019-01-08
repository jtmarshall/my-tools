# Python Monitor
A complete and automated monitoring daemon that constantly watches a list of domains(or specific urls) gathered from a database. Spin up on AWS for 100% uptime monitoring.

Use with docker/docker-compose.

- Continuously crawl homepages in domain list to check for domain health.
- Set interval for comprehensive crawl. Will course through domains starting on homepage, retrieving all links on the page and adding them to links list, and then re-iterate's for every unique page in list.
- Leverage Celery to spread tasks among sub-workers, celery-beat for scheduled tasks.
- Automated emails for 5xx alerts and weekly summary.
