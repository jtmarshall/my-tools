# Send emails with gmail
import smtplib
import email.message


# Setup variables
fromaddr = 'test@test.com'
toaddr = ['tester@test.com', 'tester2@test2.com']


def send_crawl_summary(domain, error_set):
    ### AWS Config ###
    EMAIL_HOST = 'email-smtp.us-east-1.amazonaws.com'
    EMAIL_HOST_USER = 'your_user'
    EMAIL_HOST_PASSWORD = 'your_pass'
    EMAIL_PORT = 587

    # Setup email message
    message = email.message.Message()
    message['Subject'] = "%s Crawl Summary Report" % domain
    message['From'] = fromaddr
    message['To'] = ", ".join(toaddr)
    message.add_header('Content-Type', 'text/html')

    # HTML
    template = """\
    <html>
      <head></head>
      <h1>%s Report:</h1>
      <body>
    """ % domain

    for key in error_set:
        # Concatenates the 'key', (domain page), with the status code, [0], error text, [1], and response time, [2]
        err_line = """\
            <p><h3>%s</h3> <br>
               <b>Status:</b> %s <br>
               <b>Response Time:</b> %s <br>
               <b>Error Message:</b> %s <br>
            </p>
        """ % (key, error_set[key][0], error_set[key][2], error_set[key][1])

        # Append single error string to list to join after loop
        template = template + err_line

    # Join the string list and send it as the error body message
    template_end = """\
        </body>
    </html>
    """

    template = template + template_end

    # Add html template as payload, and "stringify" content
    message.set_payload(template)
    msg_full = message.as_string()

    # Send message
    s = smtplib.SMTP(EMAIL_HOST, EMAIL_PORT)
    s.starttls()
    s.login(EMAIL_HOST_USER, EMAIL_HOST_PASSWORD)
    s.sendmail(fromaddr, toaddr, msg_full)
    s.quit()


def alert_500(webpage, notes):
    ### AWS Config ###
    EMAIL_HOST = 'email-smtp.us-east-1.amazonaws.com'
    EMAIL_HOST_USER = 'your_user'
    EMAIL_HOST_PASSWORD = 'your_pass'
    EMAIL_PORT = 587

    # Setup email message
    message = email.message.Message()
    message['Subject'] = "500 Alert!"
    message['From'] = fromaddr
    message['To'] = ", ".join(toaddr)
    message.add_header('Content-Type', 'text/html')

    # HTML
    template = """\
        <html>
          <head></head>
          <h2>500 Status Return For:</h2>
          <body> <h3>%s</h3> <p><b>Notes:</b> <br> %s</p> </body>
        </html>
        """ % (webpage, notes)

    # Add html template as payload, and "stringify" content
    message.set_payload(template)
    msg_full = message.as_string()

    # Send message
    s = smtplib.SMTP(EMAIL_HOST, EMAIL_PORT)
    s.starttls()
    s.login(EMAIL_HOST_USER, EMAIL_HOST_PASSWORD)
    s.sendmail(fromaddr, toaddr, msg_full)
    s.quit()


def monitor_alert(webpage, notes):
    ### AWS Config ###
    EMAIL_HOST = 'email-smtp.us-east-1.amazonaws.com'
    EMAIL_HOST_USER = 'your_user'
    EMAIL_HOST_PASSWORD = 'your_pass'
    EMAIL_PORT = 587

    # Setup email message
    message = email.message.Message()
    message['Subject'] = "Monitor Alert!"
    message['From'] = fromaddr
    message['To'] = ", ".join(toaddr)
    message.add_header('Content-Type', 'text/html')

    # HTML
    template = """\
        <html>
          <head></head>
          <h2>Error Returned On:</h2>
          <body> <h3>%s</h3> <p><b>Notes:</b> <br> %s</p> </body>
        </html>
        """ % (webpage, notes)

    # Add html template as payload, and "stringify" content
    message.set_payload(template)
    msg_full = message.as_string()

    # Send message
    s = smtplib.SMTP(EMAIL_HOST, EMAIL_PORT)
    s.starttls()
    s.login(EMAIL_HOST_USER, EMAIL_HOST_PASSWORD)
    s.sendmail(fromaddr, toaddr, msg_full)
    s.quit()


def domain_reactive_email(domain_name):
    ### AWS Config ###
    EMAIL_HOST = 'email-smtp.us-east-1.amazonaws.com'
    EMAIL_HOST_USER = 'your_user'
    EMAIL_HOST_PASSWORD = 'your_pass'
    EMAIL_PORT = 587

    # Setup email message
    message = email.message.Message()
    message['Subject'] = "Domain Back Online"
    message['From'] = fromaddr
    message['To'] = ", ".join(toaddr)
    # message.add_header('Content-Type', 'text/html')

    error_msg = '%s is back online.' % domain_name

    # Add error as payload, and "stringify" content
    message.set_payload(error_msg)
    msg_full = message.as_string()

    # Send message
    s = smtplib.SMTP(EMAIL_HOST, EMAIL_PORT)
    s.starttls()
    s.login(EMAIL_HOST_USER, EMAIL_HOST_PASSWORD)
    s.sendmail(fromaddr, toaddr, msg_full)
    s.quit()


def weekly_summary_email(summ_set, to_time, from_time):
    ### AWS Config ###
    EMAIL_HOST = 'email-smtp.us-east-1.amazonaws.com'
    EMAIL_HOST_USER = 'your_user'
    EMAIL_HOST_PASSWORD = 'your_pass'
    EMAIL_PORT = 587

    # Setup email message
    message = email.message.Message()
    message['Subject'] = 'Weekly Report: %s to %s \n\n' % (from_time, to_time)
    message['From'] = fromaddr
    message['To'] = ", ".join(toaddr)
    message.add_header('Content-Type', 'text/html')

    # HTML
    template = """\
    <html>
        <head></head>
        <h1>Weekly Report:</h1>
        <body>
    """

    for key in summ_set:
        summary_line = """\
            <h3> %s </h3>
            <p>
                <b>Avg Response Time:</b> %s <br>
                <b>Max Response Time:</b> %s <br>
                <b>Total URLs:</b> %s <br>
                <b>Total Errors:</b> %s <br>
                <b>Avg TTFB:</b> %s <br>
                <b>Max TTFB:</b> %s <br>
            </p>
        """ % (key, summ_set[key][0], summ_set[key][1], summ_set[key][2], summ_set[key][3],
               summ_set[key][4], summ_set[key][5])
        # Append single error string to list to join after loop
        template += summary_line

    template_end = """\
        </body>
    </html>
    """

    template = template + template_end

    # Add html template as payload, and "stringify" content
    message.set_payload(template)
    msg_full = message.as_string()

    # Send message
    s = smtplib.SMTP(EMAIL_HOST, EMAIL_PORT)
    s.starttls()
    s.login(EMAIL_HOST_USER, EMAIL_HOST_PASSWORD)
    s.sendmail(fromaddr, toaddr, msg_full)
    s.quit()
