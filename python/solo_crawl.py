"""
This script is for command line use to crawl a single domain and write data to csv.
"""
import os, sys
import re
import requests
from bs4 import BeautifulSoup
import csv


html_link_regex = re.compile('<a\s(?:.*?\s)*?href=[\'"](.*?)[\'"].*?>')


def crawl_task(url):
    domain_name = url
    url = "https://" + url

    filename = domain_name+'2.csv'

    try:
        visited = set([])  # Use set to avoid duplicate links
        request_data = requests.get(url)
        links = html_link_regex.findall(request_data.text)  # Regex gives us a list of links/paths

        csvfile = open(filename, 'w+', newline='')
        fieldnames = ['URL', 'Page Title', 'Meta Title', 'Meta Description', 'H1']
        writer = csv.DictWriter(csvfile, fieldnames=fieldnames)
        writer.writeheader()

        while len(links) > 0:
            #  Pop off the front link from the list
            link = links.pop()

            if link not in visited:
                # Add to the visited set so we don't repeat requests
                visited.add(link)

                # Try to get a response to the link and add to data set
                try:
                    # Handle whether link is relative or not
                    if domain_name in link:
                        link_url = link
                    else:
                        link_url = str(url + link)  # Attach base url to beginning of link: "/about/contact"
                    link_data = requests.get(link_url, timeout=10)  # Now get data for page link of current iteration

                    final_url = link_data.url  # If 301 redirects, the final url could be different than initial url
                    landing_data = requests.get(final_url, timeout=10)  # Get sub-links from the final url we landed on
                    sub_links = html_link_regex.findall(landing_data.text)

                    al = link_data.text
                    page_title = al[al.find('<title>') + 7: al.find('</title>')]
                    soup = BeautifulSoup(al, 'lxml')
                    metas = soup.find_all('meta')
                    mdescription = [meta.attrs['content'] for meta in metas if
                                    'name' in meta.attrs and meta.attrs['name'] == 'description']
                    mtitle = [meta.attrs['content'] for meta in metas if
                              'name' in meta.attrs and meta.attrs['name'] == 'title']
                    try:
                        h1 = soup.find_all('h1')[0].text.strip()
                    except:
                        h1 = 'None'

                    # After we loop through, Write to file depending on the status code
                    writer.writerow({
                        'URL': link_url,
                        'Page Title': page_title,
                        'Meta Title': mtitle,
                        'Meta Description': mdescription[0],
                        'H1': h1
                    })

                    #  Loop through links inside the sub page
                    for path in sub_links:
                        # If the links on the sub page have not been visited append them to the global "links" list
                        if path not in links:
                            links.insert(0, path)  # Insert at front of list to simulate a queue

                except Exception as e:
                    # If we can't get response print the failure. (phone number links will be here etc...)
                    print('Skip: ' + str(link))
                    print(e)
                    continue

        print("%s Finished." % url)

        return

    # Throw if task is hanging
    except Exception as e:
        print(e)


if __name__ == "__main__":
    print(str(sys.argv[1]))
    crawl_task(sys.argv[1])
