import re
import time
import urllib.request
import requests
import json
import redis
from datetime import datetime, timedelta
from app.alchemy import *
from sqlalchemy import exc
import app.taskdir.email.send_email as emailer

max_url_limit = 100
# New Regex that excludes anchors and query strings
default_link_regex = re.compile('<a\s(?:.*?\s)*?href=[\'"](?!\#)(\/?[^?]*?)[\'"].*?>')
red = redis.Redis(host='redis', port=6379)


# Only crawl homepages from list
def home_crawl():
    # Initialize domain list
    session = loadSession()
    res = session.query(Status.domain, Status.is_offline)
    domain_list = [[row.domain, row.is_offline] for row in res.all()]

    for domain, is_off in domain_list:
        try:
            home_url = "http://" + domain
            homepage = requests.get(home_url)

            # If we get error
            if homepage.status_code >= 500:
                if is_off == 1:
                    # already logged
                    continue
                else:
                    # Add to redis cache
                    red.set(domain, True)
                    err_time = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
                    out = Outage(domain=domain, page='home', datetime=err_time, status_code=homepage.status_code)
                    session.add(out)
                    # And update status table
                    stmt = update(Status).where(Status.domain == domain).values(is_offline=1)
                    session.execute(stmt)
                    session.commit()

                    emailer.alert_500(domain, homepage.status_code)
                    continue

            # No Error
            try:
                # If we get through without error and domain was previously off, update and send back online email
                if is_off == 1:
                    emailer.domain_reactive_email(domain)

                # And update status table
                stmt = update(Status).where(Status.domain == domain).values(is_offline=0)
                session.execute(stmt)
                session.commit()
            except Exception as e:
                print(e)
                pass
        except Exception as e:
            print(e)
            continue

    session.close()


# Metrics using SQLAlchemy
def alch_metrics(inDomain, crawl):
    session = loadSession()
    # Update Crawl end time
    db_end_log_crawl(crawl, session)
    try:
        avg_response = session.query(func.avg(Page.response_time)).filter(Page.domain == inDomain,
                                                                          Page.crawl_id == crawl).scalar()

        max_response = session.query(func.max(Page.response_time)).filter(Page.domain == inDomain,
                                                                          Page.crawl_id == crawl).scalar()

        total_urls = session.query(Page.url).filter(Page.domain == inDomain, Page.crawl_id == crawl).count()

        total_err = session.query(Page.url).filter(Page.domain == inDomain, Page.status_code != 200,
                                                   Page.crawl_id == crawl).count()

        avg_ttfb = session.query(func.avg(Page.ttfb)).filter(Page.domain == inDomain, Page.crawl_id == crawl).scalar()

        max_ttfb = session.query(func.max(Page.ttfb)).filter(Page.domain == inDomain, Page.crawl_id == crawl).scalar()

        end = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
        # Update Domain Status row for domain metrics
        stmt = update(Status).where(Status.domain == inDomain).values(datetime=end, avg_ttfb=avg_ttfb,
                                                                      max_ttfb=max_ttfb,
                                                                      avg_response=avg_response, total_errors=total_err,
                                                                      total_urls=total_urls, max_response=max_response)
        session.execute(stmt)
        session.commit()

    except:
        # Skip update if sql session error
        session.rollback()
        pass
    finally:
        session.close()


# Main comprehensive crawl task for each domain
def crawl_task(url, d_filter, is_off, crawl_id):
    domain_name = url
    html_link_regex = default_link_regex
    url = "http://" + url

    # Update regex to watch for filters set for domain in DB
    if d_filter is not None:
        json_filter = json.loads(d_filter)
        filter_string = ''
        for item in json_filter:
            filter_string += '(?!' + item + ')'
        # Insert added filters into regex string to exclude from crawl
        html_link_regex = re.compile('<a\s(?:.*?\s)*?href=[\'"](?!\#)' + filter_string + '(\/?[^?]*?)[\'"].*?>')

    try:
        session = loadSession()  # DB Session
        past_offline = is_off
        start_time = datetime.utcnow()
        count500 = 0  # Use to track 500 statuses
        visited = {url}  # Use set to avoid duplicate links
        homepage_data = requests.get(url)  # Homepage request
        links = html_link_regex.findall(homepage_data.text)  # Links on homepage

        # Switch to https if found in homepage
        if 'https' in homepage_data.url:
            url = url.replace("http", "https")

        while len(links) > 0 and len(visited) < max_url_limit:
            #  Pop off the front link from the list
            link = links.pop()

            if link not in visited:
                # Add to the visited set so we don't repeat requests
                visited.add(link)
                time_of_link = datetime.now().strftime("%Y-%m-%d %H:%M:%S")

                # Try to get a response to the link and add to data set
                try:
                    # Handle whether link is relative or not
                    if domain_name in link:
                        link_url = link
                    elif link[0] == '/':
                        # Attach base url to relative link
                        link_url = str(url + link)
                    else:
                        # Skip to next link if not in our domain
                        continue

                    # get data for page link of current iteration; will EXCEPT here if failed request
                    link_data = requests.get(link_url, timeout=10)
                    num_redirects = 0
                    error_txt = ''
                    # Whole response time
                    resp_time = link_data.elapsed.total_seconds()

                    #  If status code 504, flag in DB, requeue domain for crawl, and immediately return
                    if link_data.status_code == 504:
                        if past_offline == 1:
                            # If also 504 on previous crawl, send alert and exit.
                            emailer.alert_500(str(link_url), "504 on multiple crawls.")
                            out = Outage(domain=domain_name, datetime=time_of_link, status_code=link_data.status_code,
                                         page=link_url)
                            session.add(out)
                            session.commit()
                            break
                        else:
                            # Tag for check on next crawl in Status Table
                            offline_check = update(Status).where(Status.domain == domain_name) \
                                .values(is_offline=1, datetime=time_of_link)
                            session.execute(offline_check)
                            # Log 500 status in Page Table
                            page_out = Page(domain=domain_name, url=link, datetime=time_of_link,
                                            status_code=link_data.status_code, crawl_id=crawl_id)
                            session.add(page_out)
                            session.commit()
                            break

                    # Pass checks for each link
                    if link_data.status_code >= 500:
                        # Skip if url already in redis
                        if red.get(str(link_url)) is not None:
                            continue
                        else:
                            # Add to redis cache
                            red.set(str(link_url), True)
                        # If 5xx error send alert immediately, and increment count500
                        count500 += 1
                        emailer.alert_500(str(link_url), link_data.status_code)
                        if count500 > 2:
                            # If we get to 3 5xx status log domain in outages and kill crawl
                            out = Outage(domain=domain_name, datetime=time_of_link, status_code=link_data.status_code)
                            session.add(out)
                            session.commit()
                            break
                        # If we don't break out, log page normally in Page Table
                        page_out = Page(domain=domain_name, url=link, datetime=time_of_link,
                                        status_code=link_data.status_code, crawl_id=crawl_id)
                        session.add(page_out)
                        session.commit()
                        continue

                    if link_data.status_code != 200 or link_data.elapsed.total_seconds() > 4:
                        # If not 5xx or 200 then set error_txt for entry
                        error_txt = link_data.reason

                    # If content is html content, not an image/pdf/etc
                    if 'html' in link_data.headers.get('content-type'):
                        # Send alert if page is missing closing </html>(or </HTML>) tag
                        if '</html>' not in link_data.text and '</HTML>' not in link_data.text:
                            error_txt = '/html'
                    if link_data.history:
                        # If there are redirects then history will be more than 0
                        num_redirects = len(link_data.history)

                    # Get all links on current page
                    sub_links = html_link_regex.findall(link_data.text)

                    #  Loop through links inside the sub page
                    for path in sub_links:
                        # Only add if link has domain name or is relative
                        if (domain_name in path) or (path[0] == '/'):
                            # If the link has not already been added, add it to the "links" list
                            if path not in links:
                                links.insert(0, path)  # Insert at front of list to simulate a queue

                    # START: Time to First Byte, use urllib
                    opener = urllib.request.build_opener()
                    requester = urllib.request.Request(link_url)
                    requester.headers['Range'] = 'bytes=%s-%s' % (0, 1)
                    resp = opener.open(requester)
                    start = time.time()  # Start TTFB
                    resp.read(1)  # read one byte
                    ttfb = time.time() - start  # Sub difference after we get the first byte
                    # END: Time to First Byte

                    # Insert page (link) into page Table
                    sql = Page(domain=domain_name, url=link, datetime=time_of_link, status_code=link_data.status_code,
                               response_time=resp_time, ttfb=ttfb, redirects=num_redirects, crawl_id=crawl_id,
                               error=error_txt)
                    session.add(sql)
                    session.commit()

                    # If we get all info without error and the url was previously in cache, REMOVE from cache
                    if red.get(link_url) is not None:
                        try:
                            red.delete(link_url)
                        except:
                            pass

                except requests.exceptions.Timeout:
                    sql = Page(domain=domain_name, url=link, datetime=time_of_link, error="timeout", crawl_id=crawl_id)
                    session.add(sql)
                    session.commit()
                    continue

                except requests.exceptions.ConnectionError:
                    if domain_name in link:
                        out = Page(domain=domain_name, url=link, datetime=time_of_link, status_code=0, error="refused",
                                   crawl_id=crawl_id)
                        session.add(out)
                        session.commit()
                    continue

                except Exception as e:
                    # If we can't get response print the failure. (phone number links will be here etc...)
                    # print(e)
                    continue

        # if past_offline == 1 and count500 == 0:
        #     # Swap back the offline flag to 0, if no 500 errors
        #     update_outage = update(Status).where(Status.domain == domain_name).values(is_offline=0, datetime=start_time)
        #     session.execute(update_outage)
        #     session.commit()
        #     # Send email alert to signal domain is back online
        #     # emailer.domain_reactive_email(str(domain_name))

        # Close DB Session
        session.close()
        # Send for metrics Status Table update with crawl id
        alch_metrics(domain_name, crawl_id)

    except requests.exceptions.ConnectionError:
        out_time = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
        out = Outage(domain=domain_name, datetime=out_time, status_code=0)
        session.add(out)
        session.commit()
        session.close()
        return

    except exc.SQLAlchemyError as e:
        print("SQLAlchemy Error: ")
        print(e)
        return


def db_end_log_crawl(crawl_id, session):
    try:
        end_time = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
        total_urls = session.query(Page.id).filter(Page.crawl_id == crawl_id).count()
        stmt = update(Crawl).where(Crawl.id == crawl_id).values(end_time=end_time, total_crawled=total_urls)
        session.execute(stmt)
        session.commit()

    except Exception as e:
        print(e)


# Send a comprehensive status email weekly
def weekly_report():
    # Get datetime for now and last week
    end = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
    start = (datetime.now() - timedelta(days=7)).strftime("%Y-%m-%d %H:%M:%S")
    # DB Session
    session = loadSession()
    domain_list = [row.domain for row in session.query(Status.domain).all()]  # Switch list of dicts to list of domains
    summ_report = {}  # Dictionary to save all summary data for each domain

    for inDomain in domain_list:
        try:
            # Get metrics for each domain in the last 7 day period
            avg_response = session.query(func.avg(Page.response_time)).filter(Page.domain == inDomain,
                                                                              Page.datetime.between(start,
                                                                                                    end)).scalar()

            max_response = session.query(func.max(Page.response_time)).filter(Page.domain == inDomain,
                                                                              Page.datetime.between(start,
                                                                                                    end)).scalar()

            total_urls = session.query(Page.url).filter(Page.domain == inDomain,
                                                        Page.datetime.between(start, end)).count()

            total_err = session.query(Page.url).filter(Page.domain == inDomain, Page.status_code != 200,
                                                       Page.datetime.between(start, end)).count()

            avg_ttfb = session.query(func.avg(Page.ttfb)).filter(Page.domain == inDomain,
                                                                 Page.datetime.between(start, end)).scalar()

            max_ttfb = session.query(func.max(Page.ttfb)).filter(Page.domain == inDomain,
                                                                 Page.datetime.between(start, end)).scalar()

            # Insert metrics into dictionary
            summ_report[inDomain] = [avg_response, max_response, total_urls, total_err, avg_ttfb, max_ttfb]
        except:
            continue

    # Send dictionary, along with time frame, to weekly email func for parsing/sending
    emailer.weekly_summary_email(summ_report, end, start)
    # Close
    session.close()
